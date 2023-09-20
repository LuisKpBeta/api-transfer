package database

import (
	"database/sql"
	"errors"

	user "github.com/LuisKpBeta/api-transfer/internal/domain"
	"github.com/google/uuid"
)

type UserRepository struct {
	Db *sql.DB
}

func CreateUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

var (
	ErrUserNotFound = errors.New("user not found")
)

func (u *UserRepository) LoadUserById(id string) (*user.User, error) {
	var user user.User
	err := u.Db.QueryRow("SELECT id, name, balance FROM client WHERE id = $1", id).
		Scan(&user.Id, &user.Name, &user.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (u *UserRepository) UpdateUserBalanceById(user *user.User) error {
	_, err := u.Db.Exec("UPDATE client SET balance=$1 WHERE id = $2", user.Balance, user.Id)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepository) CreateUserTransfer(transfer user.Transfer) error {
	transfer.Id = uuid.NewString()
	_, err := u.Db.Exec("INSERT INTO transfers (id, sender_id, receiver_id, total, operation_date) VALUES ($1, $2, $3, $4, $5)",
		transfer.Id, transfer.SenderId, transfer.ReceiverId, transfer.Total, transfer.OperationDate)
	if err != nil {
		return err
	}
	return nil
}
