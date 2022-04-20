package parser

import (
	"github.com/venti-org/crawler-core/downloader"
)

type Response = downloader.Response

type Parser interface {
	Parse(response Response) []interface{}
}
