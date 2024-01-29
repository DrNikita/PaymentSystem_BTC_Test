package model

import (
	"crypto/rand"
)

const (
	ABankPrefix           = "BY04BNKA"
	ABankStateAccountIBAN = "BY04BNKA36029000000000000000"
	ABankDeadAccountIBAN  = "BY04BNKA36029111111111111111"
)

type ABank struct {
	Accounts     map[string]Account
	StateAccount *Account
	DeadAccount  *Account
}

func (b *ABank) GetAccount(iban string) Account {
	return b.Accounts[iban]
}

func (b *ABank) RegisterAccount(client Client) Account {
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

func (b *ABank) Transfer(payment *Payment) {
	payment.Sender.Withdraw(payment.Amount)
	payment.Recipient.Deposit(payment.Amount)
}

func (b *ABank) IssueMoney(amount Amount) {
	b.StateAccount.AmountOfMoney.Value += amount.Value
}

func (b *ABank) DestroyMoney(account Account, amount Amount) {
	payment := Payment{
		Sender:    &account,
		Recipient: b.DeadAccount,
		Amount:    amount,
	}
	b.Transfer(&payment)
}

func (b *ABank) GenerateIban() string {
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

	randomString := ABankPrefix + string(randomBytes[7:])

	return randomString
}
