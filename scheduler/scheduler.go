package scheduler

import (
	"github.com/venti-org/crawler-core/base"
	"github.com/venti-org/crawler-core/downloader"
)

type Request = downloader.Request

type Scheduler interface {
	Init(base.QueueFlag) error
	Put(Request)
	Pop() Request
	IsIdle() bool
	Close()
}

type DefaultScheduler struct {
	requestC base.Queue
	base.SimpleWorker
}

func NewDefaultScheduler(q base.Queue) *DefaultScheduler {
	return &DefaultScheduler{
		requestC: q,
	}
}

func (s *DefaultScheduler) Init(flag base.QueueFlag) error {
	return s.requestC.Init(flag)
}

func (s *DefaultScheduler) Put(request Request) {
	s.StartWork()
	s.requestC.Push(request)
}

func (s *DefaultScheduler) Pop() Request {
	request, _ := s.requestC.Pop().(Request)
	if request != nil {
		s.FinishWork()
	}
	return request
}

func (s *DefaultScheduler) Close() {
	s.requestC.Close()
}
