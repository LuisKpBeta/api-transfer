package api_test

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

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
func (suite *TaskApiTestSuite) UserTransferExists(id string) bool {
	var transferId string
	err := suite.Db.QueryRow("SELECT id FROM transfers WHERE sender_id = $1", id).
		Scan(&transferId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false
		}
	}
	return true
}
func (suite *TaskApiTestSuite) createTransferAndMarshall(from user.User, to user.User, total string) []byte {
	newTransfer := api.CreateTransferDto{
		Sender:   from.Id,
		Receiver: to.Id,
		Total:    total,
	}
	jsonData, _ := json.Marshal(newTransfer)
	return jsonData
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
	jsonData := suite.createTransferAndMarshall(fromUser, toUser, "50.00")

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

	jsonData := suite.createTransferAndMarshall(fromUser, toUser, "100.00")
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

func (suite *TaskApiTestSuite) Test_CreateTransfer_ParalleTest() {
	fromUser1 := user.User{
		Name:    "from user 1",
		Balance: 5000,
	}
	fromUser2 := user.User{
		Name:    "from user 2",
		Balance: 2000,
	}
	toUser := user.User{
		Name:    "to user",
		Balance: 5000,
	}
	suite.InsertUserDb(&fromUser1)
	suite.InsertUserDb(&fromUser2)
	suite.InsertUserDb(&toUser)
	firstTransfer := suite.createTransferAndMarshall(fromUser1, toUser, "50.00")
	secondTransfer := suite.createTransferAndMarshall(fromUser2, toUser, "10.00")

	testCases := []struct {
		name    string
		total   string
		result  int
		payload []byte
	}{
		{"success parallel transfer", "50.00", http.StatusOK, firstTransfer},
		{"fail parallel transfer", "20.00", http.StatusOK, secondTransfer},
	}
	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			tc := tc
			t.Parallel()
			response := suite.MakeHttpRequest(http.MethodPost, "/transfer", tc.payload)
			if !suite.Equal(tc.result, response.StatusCode()) {
				t.Errorf("got %d but waited %d", response.StatusCode(), tc.result)
			}

		})
	}
	suite.T().Cleanup(func() {
		fromUser1Balance := suite.LoadUserBalance(fromUser1.Id)
		fromUser2Balance := suite.LoadUserBalance(fromUser2.Id)
		toUserBalance := suite.LoadUserBalance(toUser.Id)

		suite.Equal(0, fromUser1Balance)
		suite.Equal(1000, fromUser2Balance)
		suite.Equal(11000, toUserBalance)
	})
}

func (suite *TaskApiTestSuite) Test_CreateTransferFailOperation() {
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
	jsonData := suite.createTransferAndMarshall(fromUser, toUser, "50.00")

	tx, _ := suite.Db.Begin()
	defer tx.Rollback()
	tx.Exec("lock transfers in exclusive mode")
	response := suite.MakeHttpRequest(http.MethodPost, "/transfer", jsonData)
	suite.Equal(http.StatusInternalServerError, response.StatusCode())
	tx.Commit()
	fromUserBalance := suite.LoadUserBalance(fromUser.Id)
	toUserBalance := suite.LoadUserBalance(toUser.Id)
	hasTransfer := suite.UserTransferExists(fromUser.Id)
	suite.Equal(5000, fromUserBalance)
	suite.Equal(5000, toUserBalance)
	suite.False(hasTransfer)

}
