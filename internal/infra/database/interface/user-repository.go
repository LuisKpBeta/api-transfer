package repository

import (
	"database/sql"

	user "github.com/LuisKpBeta/api-transfer/internal/domain"
)

type UserRepostory interface {
	LoadUserById(id string) (*user.User, error)
	UpdateUserBalanceById(tx *sql.Tx, userId string, value int) error
	CreateUserTransfer(tx *sql.Tx, transfer user.Transfer) error
	TrxLoadUserById(tx *sql.Tx, id string) (*user.User, error)
	StartTransaction() *sql.Tx
	EndTransaction(tx *sql.Tx) error
	AbortTransaction(tx *sql.Tx)
}
