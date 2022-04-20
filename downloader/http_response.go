package downloader

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type HttpResponse struct {
	Request  *HttpRequest
	Response *http.Response
	body     []byte
}

func NewHttpResponse(request *HttpRequest,
	response *http.Response) (req *HttpResponse, err error) {
	if request == nil {
		err = fmt.Errorf("HttpRequest is nil")
	} else if response == nil {
		err = fmt.Errorf("http.Response is nil")
	} else {
		req = &HttpResponse{
			Request:  request,
			Response: response,
		}
	}
	return
}

func (r *HttpResponse) GetURL() string {
	return r.Request.GetURL()
}

func (r *HttpResponse) GetMeta() Meta {
	return r.Request.GetMeta()
}

func (r *HttpResponse) GetMetaValue(key string) interface{} {
	return r.Request.GetMetaValue(key)
}

func (r *HttpResponse) GetRequest() Request {
	return r.Request
}

func (r *HttpResponse) GetBody() []byte {
	if r.body == nil {
		r.body, _ = ioutil.ReadAll(r.Response.Body)
		r.Response.Body.Close()
		if r.body == nil {
			r.body = make([]byte, 0)
		}
	}
	return r.body
}

func (r *HttpResponse) GetText() string {
	return string(r.GetBody())
}

func (r *HttpResponse) Json() map[string]interface{} {
	j := make(map[string]interface{})
	if err := json.Unmarshal(r.GetBody(), &j); err != nil {
		return nil
	}
	return j
}
