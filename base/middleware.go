package base

type OpenSpider interface {
	OpenSpider()
}

type CloseSpider interface {
	CloseSpider()
}

type BaseMiddleware struct {
}

type MiddlewareManager struct {
}
