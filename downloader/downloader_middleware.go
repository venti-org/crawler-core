package downloader

type DownloaderMiddleware interface {
	ProcessRequest(request Request) interface{}
	ProcessResponse(response Response) interface{}
}
