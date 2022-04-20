package base

import (
	"sync/atomic"
)

type Worker interface {
	StartWork()
	FinishWork()
	IsBusy() bool
}

type SimpleWorker struct {
	num int32
}

func (w *SimpleWorker) StartWork() {
	atomic.AddInt32(&w.num, 1)
}

func (w *SimpleWorker) FinishWork() {
	atomic.AddInt32(&w.num, -1)
}

func (w *SimpleWorker) IsIdle() bool {
	return w.num == 0
}
