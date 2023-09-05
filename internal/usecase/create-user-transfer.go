package usecase

import (
	"errors"

	user "github.com/LuisKpBeta/api-transfer/internal/domain"
	repository "github.com/LuisKpBeta/api-transfer/internal/infra/database/interface"
)

type CreateUserTransfer struct {
	UserRepository repository.UserRepostory
}

func (c *CreateUserTransfer) CreateUserTransfer(senderId, receiverId string, total string) error {
	value, err := user.GetBalanceValueFromString(total)
	if err != nil {
		return errors.New("invalid value for total:" + total)
	}
	newTransfer := user.Transfer{
		SenderId:   senderId,
		ReceiverId: receiverId,
		Total:      value,
	}
	senderBalance, err := c.UserRepository.LoadUserBalanceById(newTransfer.SenderId)
	if err != nil {
		return err
	}
	if senderBalance < value {
		return user.ErrSenderBalanceInvalid
	}

	return nil
}
