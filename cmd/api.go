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

	apiRouter := api.ApiRouter{
		GetUserBalance: getbalance,
	}

	router.GET("/user-balance/:id", apiRouter.GetUserBalanceById)

	log.Fatal(fasthttp.ListenAndServe(":8080", router.Handler))
}
