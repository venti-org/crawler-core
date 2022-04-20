package seeder

import (
	"github.com/venti-org/crawler-core/downloader"
)

type Request = downloader.Request

type Seeder interface {
	NextRequest() Request
	Close()
}
