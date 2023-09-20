package api_test

import (
	"database/sql"
	"testing"

	user "github.com/LuisKpBeta/api-transfer/internal/domain"
	"github.com/LuisKpBeta/api-transfer/internal/infra/api"
	"github.com/LuisKpBeta/api-transfer/internal/infra/database"
	"github.com/LuisKpBeta/api-transfer/internal/usecase"
	"github.com/buaazp/fasthttprouter"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"github.com/valyala/fasthttp"
)

type TaskApiTestSuite struct {
	suite.Suite
	Db *sql.DB
	r  *fasthttprouter.Router
}

func (suite *TaskApiTestSuite) SetupSuite() {
	suite.r = fasthttprouter.New()
	suite.Db = database.ConnectToDatabase()
	suite.SetupHttpServer()
}
func (suite *TaskApiTestSuite) TearDownSuite() {
	suite.Db.Close()
}
func (suite *TaskApiTestSuite) SetupTest() {
	suite.Db.Exec("DELETE FROM client")
}
func (suite *TaskApiTestSuite) SetupHttpServer() {
	repository := database.CreateUserRepository(suite.Db)

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

	suite.r.GET("/user-balance/:id", apiRouter.GetUserBalanceById)
	suite.r.POST("/transfer", apiRouter.CreateTransfer)
}
func (suite *TaskApiTestSuite) MakeHttpRequest(method string, url string, body []byte) *fasthttp.Response {
	ctx := &fasthttp.RequestCtx{}
	ctx.Request.Header.SetContentEncoding("application/json")
	ctx.Request.Header.SetMethod(method)
	if body != nil {
		ctx.Request.SetBody(body)
	}
	ctx.Request.SetRequestURI(url)
	suite.r.Handler(ctx)
	return &ctx.Response
}

func (suite *TaskApiTestSuite) InsertUserDb(u *user.User) {
	u.Id = uuid.NewString()
	suite.Db.Exec("INSERT INTO client (id, name, balance)  VALUES ($1, $2, $3)", u.Id, u.Name, u.Balance)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(TaskApiTestSuite))
}
