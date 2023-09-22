package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	user "github.com/LuisKpBeta/api-transfer/internal/domain"
	"github.com/LuisKpBeta/api-transfer/internal/usecase"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

type ApiRouter struct {
	GetUserBalance     usecase.GetUserBalance
	CreateUserTransfer usecase.CreateUserTransfer
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

func (a *ApiRouter) CreateTransfer(ctx *fasthttp.RequestCtx) {
	ctxBody := ctx.Request.Body()
	var bodyRequest CreateTransferDto
	err := json.Unmarshal(ctxBody, &bodyRequest)
	if err != nil {
		response := ErrorMessage{Message: err.Error()}
		MakeJsonResponse(ctx, response, http.StatusBadRequest)
		return
	}
	err = a.CreateUserTransfer.CreateUserTransfer(bodyRequest.Sender, bodyRequest.Receiver, bodyRequest.Total)
	if err != nil {
		if errors.Is(err, user.ErrSenderBalanceInvalid) || errors.Is(err, user.ErrReceiverNotFound) || errors.Is(err, user.ErrSenderNotFound) {
			response := ErrorMessage{Message: err.Error()}
			MakeJsonResponse(ctx, response, http.StatusBadRequest)
			return
		}
		response := ErrorMessage{Message: err.Error()}
		MakeJsonResponse(ctx, response, http.StatusInternalServerError)
		return
	}
	ctx.Response.SetStatusCode(http.StatusOK)
}
