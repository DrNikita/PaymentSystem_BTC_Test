package model

import (
	"crypto/rand"
)

const (
	BBankPrefix           = "BY05BNKB"
	BBankStateAccountIBAN = "BY05BNKB36029000000000000000"
	BBankDeadAccountIBAN  = "BY05BNKB36029111111111111111"
)

type BBank struct {
	Accounts     map[string]Account
	StateAccount *Account
	DeadAccount  *Account
}

func (b *BBank) GetAccount(iban string) Account {
	return b.Accounts[iban]
}

func (b *BBank) RegisterAccount(client Client) Account {
	iban := b.GenerateIban()
	client.IbanNumber = iban
	account := Account{
		Client: client,
		AmountOfMoney: &Amount{
			Value: 0,
		},
	}
	b.Accounts[iban] = account
	return b.Accounts[iban]
}

func (b *BBank) transfer(payment *Payment) error {
	if err := payment.Sender.Withdraw(payment.Amount); err != nil {
		return err
	}
	payment.Recipient.Deposit(payment.Amount)
	return nil
}

func (b *BBank) IssueMoney(amount Amount) {
	b.StateAccount.AmountOfMoney.Value += amount.Value
}

func (b *BBank) DestroyMoney(account Account, amount Amount) error {
	payment := Payment{
		Sender:    &account,
		Recipient: b.DeadAccount,
		Amount:    amount,
	}
	if err := b.transfer(&payment); err != nil {
		return err
	}
	return nil
}

func (b *BBank) GenerateIban() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	const length = 28

	randomBytes := make([]byte, length)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}

	for i := 0; i < length; i++ {
		randomBytes[i] = charset[int(randomBytes[i])%len(charset)]
	}

	randomString := BBankPrefix + string(randomBytes[7:])

	return randomString
}
