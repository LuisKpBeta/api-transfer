package main

import (
	"log"

	"github.com/LuisKpBeta/api-transfer/internal/infra/api"
	"github.com/LuisKpBeta/api-transfer/internal/infra/database"
	"github.com/LuisKpBeta/api-transfer/internal/usecase"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func main() {
	router := fasthttprouter.New()
	dbConn := database.ConnectToDatabase()
	repository := database.CreateUserRepository(dbConn)

	getbalance := usecase.GetUserBalance{
		UserRepository: repository,
	}
	createTransfer := usecase.CreateUserTransfer{
		UserRepository: repository,
	}

	apiRouter := api.ApiRouter{
		GetUserBalance:     getbalance,
		CreateUserTransfer: createTransfer,
	}

	router.GET("/user-balance/:id", apiRouter.GetUserBalanceById)
	router.POST("/transfer", apiRouter.CreateTransfer)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}

//curl -X POST -H "Content-Type: application/json" -d '{"sender_id": "c7fbc542-4c2f-11ee-be56-0242ac120002", "receiver_id": "88ee33d6-4c33-11ee-be56-0242ac120002", "total":"100"}' http://localhost:8080/transfer
//curl http://localhost:8080/user-balance/c7fbc542-4c2f-11ee-be56-0242ac120002
