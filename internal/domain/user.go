package user

import "strconv"

type User struct {
	Id      string
	Name    string
	Balance int
}

func (u *User) GetBalanceFormated() string {
	balance := strconv.Itoa(u.Balance)
	if u.Balance == 0 {
		return "0.00"
	}
	if u.Balance < 100 {
		return "0." + balance
	}
	balance = balance[:len(balance)-2] + "." + balance[len(balance)-2:]
	return balance
}
