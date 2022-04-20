package extensions

import (
	"github.com/venti-org/crawler-core/downloader"
	"github.com/venti-org/crawler-core/pipeline"
)

type Request = downloader.Request
type Response = downloader.Response
type Item = pipeline.Item

type Registry struct {
	seedRequestCallbacks     []SignalSeedRequest
	scheduledCallbacks       []SignalScheduled
	downloadSuccessCallbacks []SignalDownloadSuccess
	downloadFailureCallbacks []SignalDownloadFailure
	parsedRequestCallbacks   []SignalParsedRequest
	parsedItemCallbacks      []SignalParsedItem
	parsedUnknownCallbacks   []SignalParsedUnknown
	parsedResultCallbacks    []SignalParsedResult
	beforePipelineCallbacks  []SignalBeforePipeline
	afterPipelineCallbacks   []SignalAfterPipeline
}

func (sr *Registry) RegisterSeedRequest(m SignalSeedRequest) {
	sr.seedRequestCallbacks = append(sr.seedRequestCallbacks, m)
}

func (sr *Registry) RegisterScheduled(m SignalScheduled) {
	sr.scheduledCallbacks = append(sr.scheduledCallbacks, m)
}

func (sr *Registry) RegisterDownloadSuccess(m SignalDownloadSuccess) {
	sr.downloadSuccessCallbacks = append(sr.downloadSuccessCallbacks, m)
}

func (sr *Registry) RegisterDownloadFailure(m SignalDownloadFailure) {
	sr.downloadFailureCallbacks = append(sr.downloadFailureCallbacks, m)
}

func (sr *Registry) RegisterParsedRequest(m SignalParsedRequest) {
	sr.parsedRequestCallbacks = append(sr.parsedRequestCallbacks, m)
}

func (sr *Registry) RegisterParsedItem(m SignalParsedItem) {
	sr.parsedItemCallbacks = append(sr.parsedItemCallbacks, m)
}

func (sr *Registry) RegisterParsedUnknown(m SignalParsedUnknown) {
	sr.parsedUnknownCallbacks = append(sr.parsedUnknownCallbacks, m)
}

func (sr *Registry) RegisterParsedResult(m SignalParsedResult) {
	sr.parsedResultCallbacks = append(sr.parsedResultCallbacks, m)
}

func (sr *Registry) RegisterBeforePipeline(m SignalBeforePipeline) {
	sr.beforePipelineCallbacks = append(sr.beforePipelineCallbacks, m)
}

func (sr *Registry) RegisterAfterPipeline(m SignalAfterPipeline) {
	sr.afterPipelineCallbacks = append(sr.afterPipelineCallbacks, m)
}

func (sr *Registry) ExecSeedRequest(r Request) {
	for _, m := range sr.seedRequestCallbacks {
		m.OnSeedRequest(r)
	}
}

func (sr *Registry) ExecScheduled(r Request) {
	for _, m := range sr.scheduledCallbacks {
		m.OnScheduled(r)
	}
}

func (sr *Registry) ExecDownloadSucess(r Response) {
	for _, m := range sr.downloadSuccessCallbacks {
		m.OnDownloadSuccess(r)
	}
}

func (sr *Registry) ExecDownloadFailure(r Request, err error) {
	for _, m := range sr.downloadFailureCallbacks {
		m.OnDownloadFailure(r, err)
	}
}

func (sr *Registry) ExecParsedRequest(r Request) {
	for _, m := range sr.parsedRequestCallbacks {
		m.OnParsedRequest(r)
	}
}

func (sr *Registry) ExecParsedItem(item Item) {
	for _, m := range sr.parsedItemCallbacks {
		m.OnParsedItem(item)
	}
}

func (sr *Registry) ExecParsedUnknown(item interface{}) {
	for _, m := range sr.parsedUnknownCallbacks {
		m.OnParsedUnknown(item)
	}
}

func (sr *Registry) ExecParsedResult(result []interface{}) {
	for _, m := range sr.parsedResultCallbacks {
		m.OnParsedResult(result)
	}
}

func (sr *Registry) ExecBeforePipline(item Item) {
	for _, m := range sr.beforePipelineCallbacks {
		m.OnBeforePipeline(item)
	}
}

func (sr *Registry) ExecAfterPipline(item Item) {
	for _, m := range sr.afterPipelineCallbacks {
		m.OnAfterPipeline(item)
	}
}

func (sr *Registry) AddExtension(e interface{}) {
	if s, ok := e.(SignalSeedRequest); ok {
		sr.RegisterSeedRequest(s)
	}
	if s, ok := e.(SignalScheduled); ok {
		sr.RegisterScheduled(s)
	}
	if s, ok := e.(SignalDownloadSuccess); ok {
		sr.RegisterDownloadSuccess(s)
	}
	if s, ok := e.(SignalDownloadFailure); ok {
		sr.RegisterDownloadFailure(s)
	}
	if s, ok := e.(SignalParsedRequest); ok {
		sr.RegisterParsedRequest(s)
	}
	if s, ok := e.(SignalParsedItem); ok {
		sr.RegisterParsedItem(s)
	}
	if s, ok := e.(SignalParsedUnknown); ok {
		sr.RegisterParsedUnknown(s)
	}
	if s, ok := e.(SignalParsedResult); ok {
		sr.RegisterParsedResult(s)
	}
	if s, ok := e.(SignalBeforePipeline); ok {
		sr.RegisterBeforePipeline(s)
	}
	if s, ok := e.(SignalAfterPipeline); ok {
		sr.RegisterAfterPipeline(s)
	}
}
