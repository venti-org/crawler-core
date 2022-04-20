package downloader

import (
	"fmt"
	"net/http"

	"github.com/venti-org/crawler-core/base"
)

type AsyncDownloadResult struct {
	Request
	Response
	Err error
}

type Downloader interface {
	Download(Request) (Response, error)
	AsyncDownload(Request)
	Init(base.QueueFlag) error
	GetAsyncDownloadResult() *AsyncDownloadResult
	Start()
	Close()
	IsIdle() bool
}

type HttpDownloader struct {
	base.SimpleWorker
	client      *http.Client
	requestC    chan Request
	resultsC    base.Queue
	concurrency int
}

func NewHttpDownloader(resQ base.Queue) Downloader {
	concurrency := 10
	return &HttpDownloader{
		client:      &http.Client{},
		requestC:    make(chan Request, concurrency),
		resultsC:    resQ,
		concurrency: concurrency,
	}
}

func (d *HttpDownloader) Init(flag base.QueueFlag) error {
	return d.resultsC.Init(flag)
}

func (d *HttpDownloader) Download(request Request) (Response, error) {
	switch request := request.(type) {
	case *HttpRequest:
		response, err := d.client.Do(request.Request)
		if err != nil {
			return nil, err
		}
		return NewHttpResponse(request, response)
	default:
		return nil, fmt.Errorf("HttpDownloader not support non HttpRequest")
	}
}

func (d *HttpDownloader) AsyncDownload(request Request) {
	d.requestC <- request
}

func (d *HttpDownloader) GetAsyncDownloadResult() *AsyncDownloadResult {
	result, _ := d.resultsC.Pop().(*AsyncDownloadResult)
	return result
}

func (d *HttpDownloader) Start() {
	for i := 0; i < d.concurrency; i++ {
		go d.asyncDownloadTask()
	}
}

func (d *HttpDownloader) asyncDownloadTask() {
	for request := range d.requestC {
		d.StartWork()
		result := &AsyncDownloadResult{
			Request: request,
		}
		result.Response, result.Err = d.Download(request)
		d.resultsC.Push(result)
		d.FinishWork()
	}
}

func (d *HttpDownloader) Close() {
	close(d.requestC)
	d.resultsC.Close()
}
