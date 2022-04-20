package spider

import (
	"github.com/venti-org/crawler-core/parser"
	"github.com/venti-org/crawler-core/seeder"
)

type Spider interface {
	seeder.Seeder
	parser.Parser
}
