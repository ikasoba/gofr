package main

import (
	"syscall/js"

	"github.com/ikasoba/gofr"
)

func MainPage() js.Value {
	return Page(
		gofr.H("p", nil,
			"Gofr は好奇心から生まれたgolang用のフロントエンドフレームワークです。",
		),
		gofr.H("p", nil,
			"既存の golang 製のライブラリを使いながらWebページを構築できます。", gofr.H("br", nil),
			"また、Signal を用いて動的にDOMを書き換えることができます。",
		),
		gofr.H("dl", nil,
			gofr.H("dt", nil, "タイマー"),
			gofr.H("dd", map[string]any{
				"style": gofr.StyleDeclaration{
					"padding": "0.5rem",
				},
			}, Timer()),
			gofr.H("dt", nil, "カウンター"),
			gofr.H("dd", map[string]any{
				"style": gofr.StyleDeclaration{
					"padding": "0.5rem",
				},
			}, Counter()),
		),
	)
}

func main() {
	gofr.Render(js.Global().Get("document").Get("body"), MainPage)

	select {}
}
