package downloader

import (
	"fmt"
	http "net/http"
)

type HttpRequest struct {
	DefaultRequest
	Request *http.Request
}

func (req *HttpRequest) GetURL() string {
	return req.Request.URL.String()
}

func NewHttpRequest(request *http.Request) (req *HttpRequest, err error) {
	if request != nil {
		req = &HttpRequest{
			DefaultRequest: NewDefaultRequest(),
			Request:        request,
		}
	} else {
		err = fmt.Errorf("http.Request is nil")
	}
	return
}

func NewHttpRequestWithError(request *http.Request, e error) (*HttpRequest, error) {
	if e != nil {
		return nil, e
	}
	return NewHttpRequest(request)
}
