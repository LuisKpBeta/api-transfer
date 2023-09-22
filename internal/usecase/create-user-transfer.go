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
	trx := c.UserRepository.StartTransaction()
	defer c.UserRepository.AbortTransaction(trx)
	sender, err := c.UserRepository.TrxLoadUserById(trx, newTransfer.SenderId)
	if err != nil {
		return err
	}
	if sender == nil {
		return user.ErrSenderNotFound
	}
	if sender.Balance < value {
		return user.ErrSenderBalanceInvalid
	}
	receiver, err := c.UserRepository.TrxLoadUserById(trx, receiverId)
	if err != nil {
		return err
	}
	if receiver == nil {
		return user.ErrReceiverNotFound
	}

	if err = c.UserRepository.UpdateUserBalanceById(trx, newTransfer.ReceiverId, newTransfer.Total); err != nil {
		return err
	}
	if err = c.UserRepository.UpdateUserBalanceById(trx, newTransfer.SenderId, -newTransfer.Total); err != nil {
		return err
	}
	if err = c.UserRepository.CreateUserTransfer(trx, newTransfer); err != nil {
		return err
	}
	c.UserRepository.EndTransaction(trx)
	return nil
}
