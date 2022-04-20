package downloader

type Meta = map[string]interface{}

type Request interface {
	GetURL() string
	GetMeta() Meta
	GetMetaValue(key string) interface{}
	SetMetaValue(key string, value interface{})
}

type DefaultRequest struct {
	meta Meta
}

func NewDefaultRequest() DefaultRequest {
	return DefaultRequest{
		meta: make(Meta),
	}
}

func (r *DefaultRequest) GetMeta() Meta {
	return r.meta
}

func (r *DefaultRequest) GetMetaValue(key string) interface{} {
	return r.meta[key]
}

func (r *DefaultRequest) SetMetaValue(key string, value interface{}) {
	r.meta[key] = value
}
