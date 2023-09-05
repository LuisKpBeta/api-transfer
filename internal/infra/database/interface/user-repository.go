package repository

import user "github.com/LuisKpBeta/api-transfer/internal/domain"

type UserRepostory interface {
	LoadUserById(id string) (*user.User, error)
}
