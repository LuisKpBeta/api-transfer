package usecase

import (
	"errors"
	"time"

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
		SenderId:      senderId,
		ReceiverId:    receiverId,
		Total:         value,
		OperationDate: time.Now(),
	}
	sender, err := c.UserRepository.LoadUserById(newTransfer.SenderId)
	if err != nil {
		return err
	}
	if sender.Balance < value {
		return user.ErrSenderBalanceInvalid
	}
	receiver, err := c.UserRepository.LoadUserById(receiverId)
	if err != nil {
		return err
	}
	receiver.Balance += newTransfer.Total
	sender.Balance -= newTransfer.Total

	err = c.UserRepository.UpdateUserBalanceById(receiver)
	if err != nil {
		return err
	}
	err = c.UserRepository.UpdateUserBalanceById(sender)
	if err != nil {
		return err
	}
	err = c.UserRepository.CreateUserTransfer(newTransfer)
	if err != nil {
		return err
	}
	return nil
}
