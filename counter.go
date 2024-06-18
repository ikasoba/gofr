package main

import (
	"syscall/js"

	"github.com/ikasoba/gofr"
)

func Counter() js.Value {
	count := gofr.NewSignal(0)

	return gofr.H("button", map[string]any{
		"onclick": js.FuncOf(func(this js.Value, args []js.Value) any {
			count.Set(count.Get() + 1)

			return nil
		}),
	}, "count: ", count)
}
