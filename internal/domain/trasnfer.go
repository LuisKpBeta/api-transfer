package user

import (
	"errors"
	"strconv"
	"strings"
	"time"
)

type Transfer struct {
	SenderId      string
	ReceiverId    string
	Total         int
	OperationDate time.Time
}

var (
	ErrSenderBalanceInvalid = errors.New("sender don't have enought balance")
)

func GetBalanceValueFromString(value string) (int, error) {
	valueWithouDot := strings.ReplaceAll(value, ".", "")
	valueWithouDot = strings.ReplaceAll(valueWithouDot, ",", "")
	return strconv.Atoi(valueWithouDot)
}
