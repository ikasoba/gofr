package main

import (
	"syscall/js"

	"github.com/ikasoba/gofr"
)

func SiteHeader() js.Value {
	return gofr.H("header", map[string]any{
		"style": gofr.StyleDeclaration{
			"display":         "flex",
			"align-content":   "center",
			"justify-content": "space-between",
			"padding":         "0.25rem 0.5rem",
			"background":      "rgb(1 173 216)",
			"color":           "white",
		},
	},
		gofr.H("h1", map[string]any{
			"style": gofr.StyleDeclaration{
				"margin": "0",
			},
		}, "Gofr"),
		gofr.H("div", map[string]any{
			"style": gofr.StyleDeclaration{
				"display":       "flex",
				"align-content": "center",
				"flex-wrap":     "wrap",
				"gap":           "1rem",
			},
		},
			gofr.H("a", map[string]any{
				"href": "/",
			}, "トップ"),
			gofr.H("a", map[string]any{
				"href": "https://github.com/ikasoba/gofr",
			}, "GitHub"),
		),
	)
}

func Page(children ...any) js.Value {
	return gofr.H("div", nil,
		SiteHeader(),
		gofr.H("main", map[string]any{
			"style": gofr.StyleDeclaration{
				"padding": "0.5rem 2rem",
			},
		},
			children...,
		),
	)
}
