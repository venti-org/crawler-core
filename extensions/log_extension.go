package extensions

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type LogExtension struct {
}

func NewLogExtension() *LogExtension {
	return &LogExtension{}
}

func (e *LogExtension) OnSeedRequest(request Request) {
	logrus.WithFields(logrus.Fields{
		"url": request.GetURL(),
	}).Infoln("generate seed")
}

func (e *LogExtension) OnScheduled(request Request) {
	logrus.WithFields(logrus.Fields{
		"url": request.GetURL(),
	}).Infoln("request to download")
}

func (e *LogExtension) OnDownloadSuccess(response Response) {
	logrus.WithFields(logrus.Fields{
		"url": response.GetURL(),
	}).Infoln("response to parse")
}

func (e *LogExtension) OnDownloadFailure(request Request, err error) {
	logrus.WithFields(logrus.Fields{
		"url": request.GetURL(),
		"err": err,
	}).Infoln("download error")
}

func (e *LogExtension) OnParsedUnknown(item interface{}) {
	logrus.Errorf("parser error return type: %T\n", item)
}

func (e *LogExtension) OnAfterPipeline(item Item) {
	if body, err := json.Marshal(item); err != nil {
		logrus.Errorln(err)
	} else {
		json_str := string(body)
		logrus.Println(json_str)
	}
}
