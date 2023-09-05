package database

import (
	"database/sql"

	user "github.com/LuisKpBeta/api-transfer/internal/domain"
)

type UserRepository struct {
	Db *sql.DB
}

func CreateUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{
		Db: db,
	}
}

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
