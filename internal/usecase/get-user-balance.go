package usecase

import (
	repository "github.com/LuisKpBeta/api-transfer/internal/infra/database/interface"
)

type GetUserBalance struct {
	UserRepository repository.UserRepostory
}

type ReadUser struct {
	Id      string
	Balance string
}

func (g *GetUserBalance) GetUserBalance(id string) (*ReadUser, error) {
	user, err := g.UserRepository.LoadUserById(id)
	if err != nil || user == nil {
		return nil, err
	}

	readUser := ReadUser{
		Id:      user.Id,
		Balance: user.GetBalanceFormated(),
	}
	return &readUser, nil
}
