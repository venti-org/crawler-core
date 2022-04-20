package extensions

type SignalSeedRequest interface {
	OnSeedRequest(Request)
}

type SignalScheduled interface {
	OnScheduled(Request)
}

type SignalDownloadSuccess interface {
	OnDownloadSuccess(Response)
}

type SignalDownloadFailure interface {
	OnDownloadFailure(Request, error)
}

type SignalParsedRequest interface {
	OnParsedRequest(Request)
}

type SignalParsedUnknown interface {
	OnParsedUnknown(interface{})
}

type SignalParsedItem interface {
	OnParsedItem(Item)
}

type SignalParsedResult interface {
	OnParsedResult([]interface{})
}

type SignalBeforePipeline interface {
	OnBeforePipeline(Item)
}

type SignalAfterPipeline interface {
	OnAfterPipeline(Item)
}
