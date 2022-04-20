package base

import (
	"os"
	"os/signal"
	"syscall"
)

func OnSignal(f func()) {
	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)
		<-signals
		f()
	}()
}
