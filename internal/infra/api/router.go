package api

import (
	"fmt"
	"net/http"

	"github.com/LuisKpBeta/api-transfer/internal/usecase"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type ApiRouter struct {
	GetUserBalance usecase.GetUserBalance
}

func (a *ApiRouter) GetUserBalanceById(ctx *fasthttp.RequestCtx) {
	param := ctx.Value("id")
	ctx.Response.Header.SetContentType("application/json")

	if param == nil {
		ctx.Response.SetStatusCode(http.StatusBadRequest)
		return
	}
	id := fmt.Sprint(param)
	_, err := uuid.Parse(id)
	if err != nil {
		response := ErrorMessage{Message: "invalid user id"}
		MakeJsonResponse(ctx, response, http.StatusInternalServerError)
		return
	}
	user, err := a.GetUserBalance.GetUserBalance(id)
	if err != nil {
		response := ErrorMessage{Message: err.Error()}
		MakeJsonResponse(ctx, response, http.StatusInternalServerError)
		return
	}
	if user == nil {
		ctx.Response.SetStatusCode(http.StatusNotFound)
		return
	}
	MakeJsonResponse(ctx, user, http.StatusOK)
}
