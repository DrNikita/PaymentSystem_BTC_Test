package model

type Banker interface {
	GetAccount(iban string) Account
	RegisterAccount(client Client) Account
	IssueMoney(amount Amount)
	DestroyMoney(account Account, amount Amount) error
	GenerateIban() string
}
