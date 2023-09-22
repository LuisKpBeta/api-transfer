package user

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Transfer struct {
	Id            string
	SenderId      string
	ReceiverId    string
	Total         int
	OperationDate time.Time
}

var (
	ErrSenderBalanceInvalid = errors.New("sender don't have enought balance")
	ErrReceiverNotFound     = errors.New("receiver not found for informed id")
	ErrSenderNotFound       = errors.New("receiver not found for informed id")
)

func GetBalanceValueFromString(value string) (int, error) {
	valueWithouDot := strings.ReplaceAll(value, ".", "")
	valueWithouDot = strings.ReplaceAll(valueWithouDot, ",", "")
	return strconv.Atoi(valueWithouDot)
}
