package database

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

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
func (u *UserRepository) StartTransaction() *sql.Tx {
	tx, _ := u.Db.Begin()
	return tx
}
func (u *UserRepository) AbortTransaction(tx *sql.Tx) {
	log.Println("aborting transaction")
	err := tx.Rollback()
	log.Println(err.Error())
}
func (u *UserRepository) EndTransaction(tx *sql.Tx) error {
	log.Println("transaction successfull")
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
func (u *UserRepository) TrxLoadUserById(tx *sql.Tx, id string) (*user.User, error) {
	var user user.User
	err := tx.QueryRow("SELECT id, name, balance FROM client WHERE id = $1", id).
		Scan(&user.Id, &user.Name, &user.Balance)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
func (u *UserRepository) UpdateUserBalanceById(tx *sql.Tx, userId string, value int) error {
	_, err := tx.Exec("UPDATE client SET balance = balance + $1 WHERE id = $2", value, userId)
	if err != nil {
		return err
	}
	return nil
}
func (u *UserRepository) CreateUserTransfer(tx *sql.Tx, transfer user.Transfer) error {
	transfer.Id = uuid.NewString()
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	_, err := tx.ExecContext(ctx, "INSERT INTO transfers (id, sender_id, receiver_id, total, operation_date) VALUES ($1, $2, $3, $4, $5)",
		transfer.Id, transfer.SenderId, transfer.ReceiverId, transfer.Total, transfer.OperationDate)
	if err != nil {
		return err
	}
	return nil
}
