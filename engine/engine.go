package engine

import (
	"fmt"
	"sync"
	"time"

	"github.com/venti-org/crawler-core/base"
	"github.com/venti-org/crawler-core/downloader"
	"github.com/venti-org/crawler-core/extensions"
	"github.com/venti-org/crawler-core/parser"
	"github.com/venti-org/crawler-core/pipeline"
	"github.com/venti-org/crawler-core/scheduler"
	"github.com/venti-org/crawler-core/seeder"
	"github.com/venti-org/crawler-core/spider"
	"github.com/sirupsen/logrus"
)

type Request = downloader.Request
type Response = downloader.Response
type Item = pipeline.Item

type Engine struct {
	extensions.Registry

	Downloader   downloader.Downloader
	Seeder       seeder.Seeder
	Parser       parser.Parser
	Scheduler    scheduler.Scheduler
	Pipeline     pipeline.Pipeline
	seederWg     sync.WaitGroup
	downloaderWg sync.WaitGroup
	parserWg     sync.WaitGroup
	pipelineWg   sync.WaitGroup

	itemC base.Queue

	parserWorker   base.SimpleWorker
	pipelineWorker base.SimpleWorker

	parserConcurrency   int
	pipelineConcurrency int

	runing bool
}

func NewEngine(s spider.Spider) (*Engine, error) {
	return NewEngineBuilder().WithSpider(s).Build()
}

func (engine *Engine) mustInitFlag(name string, push bool, pop bool) {
	flag := base.QueueFlag{
		Push: push,
		Pop:  pop,
	}
	var err error
	if name == "scheduler" {
		err = engine.Scheduler.Init(flag)
	} else if name == "downloader" {
		err = engine.Downloader.Init(flag)
	} else if name == "itemC" {
		err = engine.itemC.Init(flag)
	} else {
		err = fmt.Errorf("not support %s init", name)
	}

	if err != nil {
		panic(fmt.Errorf("init %s error: %s", name, err.Error()))
	}
}

func (engine *Engine) Run() {
	if engine.runing {
		panic("not Run when running")
	}

	engine.mustInitFlag("scheduler", true, true)
	engine.mustInitFlag("downloader", true, true)
	engine.mustInitFlag("itemC", true, true)

	engine.seederWg.Add(1)
	go engine.generate_seed_task()

	engine.Downloader.Start()
	engine.downloaderWg.Add(1)
	go engine.download_task()

	for i := 0; i < engine.parserConcurrency; i++ {
		engine.parserWg.Add(1)
		go engine.parse_task()
	}

	for i := 0; i < engine.pipelineConcurrency; i++ {
		engine.pipelineWg.Add(1)
		go engine.pipeline_task()
	}

	engine.seederWg.Wait()

	go engine.monitor_sleep_task()

	engine.pipelineWg.Wait()
}

func (engine *Engine) RunModle(module string) {
	if engine.runing {
		panic("not Run when running")
	}

	engine.runing = true
	if module == "seeder" {
		engine.mustInitFlag("scheduler", true, false)
		engine.seederWg.Add(1)
		go engine.generate_seed_task()
		base.OnSignal(func() {
			engine.Seeder.Close()
		})
		engine.seederWg.Wait()
	} else if module == "downloader" {
		engine.mustInitFlag("scheduler", false, true)
		engine.mustInitFlag("downloader", true, false)
		engine.Downloader.Start()
		engine.downloaderWg.Add(1)
		go engine.download_task()
		base.OnSignal(func() {
			engine.Scheduler.Close()
		})
		engine.downloaderWg.Wait()
	} else if module == "parser" {
		engine.mustInitFlag("downloader", false, true)
		engine.mustInitFlag("scheduler", true, false)
		engine.mustInitFlag("itemC", true, false)
		for i := 0; i < engine.parserConcurrency; i++ {
			engine.parserWg.Add(1)
			go engine.parse_task()
		}
		base.OnSignal(func() {
			engine.Downloader.Close()
		})
		engine.parserWg.Wait()
	} else if module == "pipeline" {
		engine.mustInitFlag("itemC", true, true)
		for i := 0; i < engine.pipelineConcurrency; i++ {
			engine.pipelineWg.Add(1)
			go engine.pipeline_task()
		}
		base.OnSignal(func() {
			engine.itemC.Close()
		})
		engine.pipelineWg.Wait()
	} else {
		panic("not support module: " + module)
	}
}

func (engine *Engine) monitor_sleep_task() {
	ticker := time.NewTicker(time.Second)
	times := 0
	for range ticker.C {
		if engine.Scheduler.IsIdle() &&
			engine.Downloader.IsIdle() &&
			engine.parserWorker.IsIdle() &&
			engine.pipelineWorker.IsIdle() {
			times++
			if times > 1 {
				logrus.Infoln("stop ticker")
				ticker.Stop()
				break
			}
		} else {
			times = 0
		}
	}

	engine.Scheduler.Close()
	logrus.Infoln("wait downloaderWg")
	engine.downloaderWg.Wait()
	engine.Downloader.Close()
	logrus.Infoln("wait parserWg")
	engine.parserWg.Wait()
	engine.itemC.Close()
	logrus.Infoln("close Monitor")
}

func (engine *Engine) generate_seed_task() {
	defer engine.seederWg.Done()

	for {
		request := engine.Seeder.NextRequest()
		if request == nil {
			break
		}
		engine.ExecSeedRequest(request)
		engine.Scheduler.Put(request)
	}
	logrus.Infoln("close Seeder")
}

func (engine *Engine) download_task() {
	defer engine.downloaderWg.Done()

	for {
		request := engine.Scheduler.Pop()
		if request == nil {
			break
		}
		engine.ExecScheduled(request)
		engine.Downloader.AsyncDownload(request)
	}
	logrus.Infoln("close downloader")
}

func (engine *Engine) parse_task() {
	defer engine.parserWg.Done()

	for {
		result := engine.Downloader.GetAsyncDownloadResult()
		if result == nil {
			break
		}
		engine.parserWorker.StartWork()
		if result.Err != nil {
			engine.ExecDownloadFailure(result.Request, result.Err)
			engine.parserWorker.FinishWork()
			continue
		}
		engine.ExecDownloadSucess(result.Response)
		response := result.Response
		items := engine.Parser.Parse(response)
		engine.ExecParsedResult(items)
		for _, item := range items {
			switch item := item.(type) {
			case Item:
				engine.ExecParsedItem(item)
				engine.itemC.Push(item)
			case Request:
				engine.ExecParsedRequest(item)
				engine.Scheduler.Put(item)
			default:
				engine.ExecParsedUnknown(item)
			}
		}
		engine.parserWorker.FinishWork()
	}
	logrus.Infoln("close parser")
}

func (engine *Engine) pipeline_task() {
	defer engine.pipelineWg.Done()

	for {
		item, _ := engine.itemC.Pop().(Item)
		if item == nil {
			break
		}
		engine.ExecBeforePipline(item)
		engine.pipelineWorker.StartWork()
		engine.Pipeline.ProcessItem(item)
		engine.pipelineWorker.FinishWork()
		engine.ExecAfterPipline(item)
	}
	logrus.Infoln("close pipeline")
}
