package engine

import (
	"fmt"

	"github.com/venti-org/crawler-core/base"
	"github.com/venti-org/crawler-core/downloader"
	"github.com/venti-org/crawler-core/parser"
	"github.com/venti-org/crawler-core/pipeline"
	"github.com/venti-org/crawler-core/scheduler"
	"github.com/venti-org/crawler-core/seeder"
	"github.com/venti-org/crawler-core/spider"
)

type EngineBuilder struct {
	engine *Engine
}

func NewEngineBuilder() *EngineBuilder {
	return &EngineBuilder{
		engine: &Engine{
			Pipeline:            pipeline.NewDefaultPipeline(),
			parserConcurrency:   1,
			pipelineConcurrency: 1,
		},
	}
}

func (builder *EngineBuilder) WithItemQueue(q base.Queue) *EngineBuilder {
	builder.engine.itemC = q
	return builder
}

func (builder *EngineBuilder) WithScheduler(s scheduler.Scheduler) *EngineBuilder {
	builder.engine.Scheduler = s
	return builder
}

func (builder *EngineBuilder) WithSchedulerQueue(q base.Queue) *EngineBuilder {
	return builder.WithScheduler(scheduler.NewDefaultScheduler(q))
}

func (builder *EngineBuilder) WithDownloader(d downloader.Downloader) *EngineBuilder {
	builder.engine.Downloader = d
	return builder
}

func (builder *EngineBuilder) WithDownloaderQueue(resQ base.Queue) *EngineBuilder {
	return builder.WithDownloader(downloader.NewHttpDownloader(resQ))
}

func (builder *EngineBuilder) WithSeeder(s seeder.Seeder) *EngineBuilder {
	builder.engine.Seeder = s
	return builder
}

func (builder *EngineBuilder) WithParser(p parser.Parser) *EngineBuilder {
	builder.engine.Parser = p
	return builder
}

func (builder *EngineBuilder) WithSpider(s spider.Spider) *EngineBuilder {
	return builder.WithSeeder(s).WithParser(s)
}

func (builder *EngineBuilder) SetParserConcurrency(c int) *EngineBuilder {
	builder.engine.parserConcurrency = c
	return builder
}

func (builder *EngineBuilder) SetPipelineConcurrency(c int) *EngineBuilder {
	builder.engine.pipelineConcurrency = c
	return builder
}

func (builder *EngineBuilder) AppendComponents(components ...pipeline.Component) *EngineBuilder {
	builder.engine.Pipeline.AppendComponents(components...)
	return builder
}

func (builder *EngineBuilder) Build() (engine *Engine, err error) {
	engine = builder.engine
	if engine.Seeder == nil {
		err = fmt.Errorf("Engine.Seeder is Nil")
	} else if engine.Parser == nil {
		err = fmt.Errorf("Engine.Parser is Nil")
	}
	if err != nil {
		return nil, err
	}
	if engine.Downloader == nil {
		builder.WithDownloaderQueue(base.NewChannel(make(chan interface{}, 100)))
	}
	if engine.Scheduler == nil {
		builder.WithSchedulerQueue(base.NewChannel(make(chan interface{}, 100)))
	}
	if engine.itemC == nil {
		builder.WithItemQueue(base.NewChannel(make(chan interface{}, 100)))
	}
	return
}
