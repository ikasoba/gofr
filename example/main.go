//go:build js

package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/ikasoba/gofr"
)

func Counter() js.Value {
	count := gofr.NewSignal(0)

	return gofr.H("button", map[string]any{
		"onclick": js.FuncOf(func(this js.Value, args []js.Value) any {
			count.Set(count.Get() + 1)

			return nil
		}),
		"style": gofr.Style(map[string]any{
			"border": gofr.Computed(func() string {
				return fmt.Sprint("0.25rem solid hsl(", count.Get()*10%360, "deg 100% 50%)")
			}),
			"border-radius": "0.5rem",
		}),
	},
		"count: ", count,
	)
}

func Spoiler(Component func() js.Value) js.Value {
	isHidden := gofr.NewSignal(true)

	return gofr.H("div", map[string]any{
		"style": gofr.Style(map[string]any{
			"display":        "inline-flex",
			"flex-direction": "column",
			"gap":            "0.125rem",
		}),
	},
		gofr.H("button", map[string]any{
			"onclick": js.FuncOf(func(this js.Value, args []js.Value) any {
				isHidden.Set(!isHidden.Get())

				return nil
			}),
		},
			gofr.Computed(func() string {
				if isHidden.Get() {
					return "show"
				} else {
					return "hide"
				}
			}),
		),
		gofr.Computed(func() any {
			if isHidden.Get() {
				return nil
			} else {
				return Component()
			}
		}),
	)
}

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

	return gofr.H("span", map[string]any{}, count, "s")
}

func App() js.Value {
	return gofr.H("div", map[string]any{
		"style": gofr.Style(map[string]any{
			"display":        "inline-flex",
			"flex-direction": "column",
			"gap":            "0.5rem",
		}),
	},
		Counter(),
		Spoiler(Timer),
	)
}

func main() {
	gofr.Render(js.Global().Get("document").Get("body"), App)

	select {}
}
