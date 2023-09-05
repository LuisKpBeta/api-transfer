package api

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type ErrorMessage struct {
	Message string `json:"message"`
}

func MakeJsonResponse(ctx *fasthttp.RequestCtx, value interface{}, statusCode int) {
	ctx.Response.SetStatusCode(statusCode)
	responseJSON, _ := json.Marshal(value)
	ctx.Write(responseJSON)
}
