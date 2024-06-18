package main

import (
	"syscall/js"
	"time"

	"github.com/ikasoba/gofr"
)

func Timer() js.Value {
	count := gofr.NewSignal(0)

	ticker := time.NewTicker(time.Second)
	cleanupped := make(chan bool)

	go func() {
	l:
		for {
			select {
			case <-ticker.C:
				count.Set(count.Get() + 1)

			case <-cleanupped:
				break l
			}
		}
	}()

	gofr.OnCleanup(func() {
		ticker.Stop()
		cleanupped <- true
	})

	return gofr.H("span", nil, count, "秒")
}
