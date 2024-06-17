package gofr

import (
	"fmt"
	"strings"
	"syscall/js"
)

var (
	document = js.Global().Get("document")
)

type AnyTransmitter interface {
	Transmitter

	GetAsAny() any
}

type NodeConverter interface {
	ToNode() js.Value
}

type NodeAttacher interface {
	AttachNode(key string, el js.Value)
}

func newNode(value any) js.Value {
	switch v := value.(type) {
	case js.Value:
		return v

	default:
		if v == nil {
			return document.Call("createComment", "")
		}

		return document.Call("createTextNode", v)
	}
}

func replaceNode(value any, prev js.Value) js.Value {
	switch v := value.(type) {
	case js.Value:
		prev.Get("parentNode").Call("replaceChild", v, prev)

		return v

	default:
		if v == nil {
			c := document.Call("createComment", "")
			prev.Get("parentNode").Call("replaceChild", c, prev)

			return c
		}

		prev.Set("textContent", fmt.Sprint(v))
		return prev
	}
}

func H(name string, attrs map[string]any, children ...any) js.Value {
	el := document.Call("createElement", name)

	for k, value := range attrs {
		if strings.HasPrefix(k, "on") {
			switch v := value.(type) {
			case AnyTransmitter:
				v.Subscribe(func() {
					el.Call("addEventListener", k[2:], v.GetAsAny())
				})

				el.Call("addEventListener", k[2:], v.GetAsAny())

			default:
				el.Call("addEventListener", k[2:], v)
			}
		} else {
			switch v := value.(type) {
			case AnyTransmitter:
				v.Subscribe(func() {
					el.Call("setAttribute", k, v.GetAsAny())
				})

				el.Call("setAttribute", k, v.GetAsAny())

			case NodeAttacher:
				v.AttachNode(k, el)

			default:
				el.Call("setAttribute", k, v)
			}
		}
	}

	var nodes []any

	for _, item := range children {
		switch val := item.(type) {
		case js.Value:
			nodes = append(nodes, val)

		case NodeConverter:
			nodes = append(nodes, val.ToNode())

		default:
			nodes = append(nodes, fmt.Sprint(val))
		}
	}

	el.Call("append", nodes...)

	return el
}

func Render(root js.Value, fn func() js.Value) *LifeTime {
	escapeLifeTime := BeginLifeTime()
	el := fn()
	lifeTime := escapeLifeTime()

	root.Call("append", el)

	lifeTime.DispatchMount()

	return lifeTime
}
