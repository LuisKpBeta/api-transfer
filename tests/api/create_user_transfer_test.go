package api_test

import (
	"encoding/json"
	"net/http"

	user "github.com/LuisKpBeta/api-transfer/internal/domain"
	"github.com/LuisKpBeta/api-transfer/internal/infra/api"
)

func (suite *TaskApiTestSuite) LoadUserBalance(id string) int {
	var balance int
	err := suite.Db.QueryRow("SELECT balance FROM client WHERE id = $1", id).
		Scan(&balance)
	if err != nil {
		panic(err)
	}
	return balance
}
func (suite *TaskApiTestSuite) Test_CreateTransferSuccess() {
	fromUser := user.User{
		Name:    "from user",
		Balance: 10000,
	}
	toUser := user.User{
		Name:    "to user",
		Balance: 5000,
	}
	suite.InsertUserDb(&fromUser)
	suite.InsertUserDb(&toUser)
	newTransfer := api.CreateTransferDto{
		Sender:   fromUser.Id,
		Receiver: toUser.Id,
		Total:    "50.00",
	}
	jsonData, _ := json.Marshal(newTransfer)

	response := suite.MakeHttpRequest(http.MethodPost, "/transfer", jsonData)

	testErr := suite.Equal(http.StatusOK, response.StatusCode())
	if !testErr {
		return
	}

	fromUserBalance := suite.LoadUserBalance(fromUser.Id)
	toUserBalance := suite.LoadUserBalance(toUser.Id)

	suite.Equal(5000, fromUserBalance)
	suite.Equal(10000, toUserBalance)
}
func (suite *TaskApiTestSuite) Test_CreateTransferInsufficientBalance() {
	fromUser := user.User{
		Name:    "from user",
		Balance: 5000,
	}
	toUser := user.User{
		Name:    "to user",
		Balance: 5000,
	}
	suite.InsertUserDb(&fromUser)
	suite.InsertUserDb(&toUser)
	newTransfer := api.CreateTransferDto{
		Sender:   fromUser.Id,
		Receiver: toUser.Id,
		Total:    "100.00",
	}
	jsonData, _ := json.Marshal(newTransfer)

	response := suite.MakeHttpRequest(http.MethodPost, "/transfer", jsonData)

	testErr := suite.Equal(http.StatusBadRequest, response.StatusCode())
	if !testErr {
		return
	}
	errMessage := &api.ErrorMessage{}
	if err := json.Unmarshal(response.Body(), errMessage); err != nil {
		suite.Error(err)
	}
	suite.Equal(user.ErrSenderBalanceInvalid.Error(), errMessage.Message)

	fromUserBalance := suite.LoadUserBalance(fromUser.Id)
	toUserBalance := suite.LoadUserBalance(toUser.Id)

	suite.Equal(5000, fromUserBalance)
	suite.Equal(5000, toUserBalance)
}
