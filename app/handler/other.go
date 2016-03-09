package handler

import (
	"github.com/dinever/golf"
)

func NotFoundHandler(ctx *Golf.Context) {
	ctx.StatusCode = 404
	data := map[string]interface{}{}
	ctx.Loader("theme").Render("404.html", data)
}
