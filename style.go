package gofr

import "syscall/js"

type StyleDeclaration map[string]any

func Style(style map[string]any) StyleDeclaration {
	return StyleDeclaration(style)
}

func (s StyleDeclaration) AttachNode(_k string, el js.Value) {
	for key, value := range s {
		k := key

		switch v := value.(type) {
		case AnyTransmitter:
			v.Subscribe(func() {
				el.Get("style").Call("setProperty", k, v.GetAsAny())
			})

			el.Get("style").Call("setProperty", k, v.GetAsAny())

		default:
			el.Get("style").Call("setProperty", k, v)
		}
	}
}
