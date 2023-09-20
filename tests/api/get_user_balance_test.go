package api_test

import (
	"encoding/json"
	"net/http"

	user "github.com/LuisKpBeta/api-transfer/internal/domain"
	"github.com/LuisKpBeta/api-transfer/internal/usecase"
)

func (suite *TaskApiTestSuite) Test_GetClientBalance() {
	newUser := user.User{
		Name:    "test user",
		Balance: 100,
	}
	suite.InsertUserDb(&newUser)
	url := "/user-balance/" + newUser.Id

	response := suite.MakeHttpRequest(http.MethodGet, url, nil)

	suite.Equal(http.StatusOK, response.StatusCode())
	readUser := &usecase.ReadUser{}
	if err := json.Unmarshal(response.Body(), readUser); err != nil {
		suite.Error(err)
	}

	suite.Equal(readUser.Id, newUser.Id)
	suite.Equal(readUser.Balance, "1.00")

}
