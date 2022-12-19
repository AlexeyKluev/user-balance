package model

import "github.com/Rhymond/go-money"

type User struct {
	ID        int64
	FirstName string
	LastName  string
	Balance   *money.Money
}
