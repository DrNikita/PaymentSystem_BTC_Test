package model

import (
	"errors"
	"sync"
)

type Account struct {
	rw            sync.RWMutex
	Client        Client
	AmountOfMoney *Amount
}

func (a *Account) Withdraw(amount Amount) error {
	a.rw.Lock()
	defer a.rw.Unlock()
	if a.AmountOfMoney.Value < amount.Value {
		return errors.New("not enough money")
	}
	a.AmountOfMoney.Value = a.AmountOfMoney.Value - amount.Value
	return nil
}

func (a *Account) Deposit(amount Amount) {
	a.rw.Lock()
	defer a.rw.Unlock()
	a.AmountOfMoney.Value = a.AmountOfMoney.Value + amount.Value
}
