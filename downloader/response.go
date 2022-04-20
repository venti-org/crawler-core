package downloader

type Response interface {
	GetURL() string
	GetMeta() Meta
	GetRequest() Request
	GetBody() []byte
}
