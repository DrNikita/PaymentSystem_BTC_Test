package main

import (
	"fmt"
	"github.com/DrNikita/BTC_Test/internal/model"
	log "github.com/sirupsen/logrus"
	"sync"
)

func main() {
	bankA := model.ABank{
		Accounts: make(map[string]model.Account),
		StateAccount: &model.Account{
			Client: model.Client{
				IbanNumber: model.ABankStateAccountIBAN,
				Name:       "State account",
			},
			AmountOfMoney: &model.Amount{
				Value: 0,
			},
		},
		DeadAccount: &model.Account{
			Client: model.Client{
				IbanNumber: model.ABankDeadAccountIBAN,
				Name:       "Dead account",
			},
			AmountOfMoney: &model.Amount{
				Value: 0,
			},
		},
	}

	bankB := model.BBank{
		Accounts: make(map[string]model.Account),
		StateAccount: &model.Account{
			Client: model.Client{
				IbanNumber: model.BBankStateAccountIBAN,
				Name:       "State account",
			},
			AmountOfMoney: &model.Amount{
				Value: 0,
			},
		},
		DeadAccount: &model.Account{
			Client: model.Client{
				IbanNumber: model.BBankDeadAccountIBAN,
				Name:       "Dead account",
			},
			AmountOfMoney: &model.Amount{
				Value: 0,
			},
		},
	}

	visa := model.Visa{
		RegisteredBanks: make(map[string]model.Banker),
	}
	visa.RegisterBank(model.ABankPrefix, &bankA)
	visa.RegisterBank(model.BBankPrefix, &bankB)

	// регистрация аккаунтов
	ANikita := bankA.RegisterAccount(model.Client{
		Name:       "Nikita",
		Surname:    "Chelovek",
		Patronymic: "Otchestvo",
	})
	BEgor := bankB.RegisterAccount(model.Client{
		Name:       "Egor",
		Surname:    "Chelovek",
		Patronymic: "Otchestvo",
	})
	BMorf := bankB.RegisterAccount(model.Client{
		Name:       "Morfling",
		Surname:    "Mineral",
		Patronymic: "Water",
	})
	AEmber := bankA.RegisterAccount(model.Client{
		Name:       "Ember",
		Surname:    "Spirit",
		Patronymic: "Fire",
	})

	// пополнение баланса
	ANikita.Deposit(model.Amount{
		Value: 5000,
	})
	BEgor.Deposit(model.Amount{
		Value: 15000,
	})
	BMorf.Deposit(model.Amount{
		Value: 16700,
	})
	AEmber.Deposit(model.Amount{
		Value: 19980,
	})

	var wg sync.WaitGroup
	errChan := make(chan error)

	wg.Add(4)

	go func() {
		defer wg.Done()
		amount := model.Amount{
			Value: 9700,
		}
		log.Info(fmt.Sprintf("Gorutine 1: Ember(%f)--%f-->Morf(%f)", AEmber.AmountOfMoney.Value, amount.Value, BMorf.AmountOfMoney.Value))
		err := visa.Transfer(AEmber.Client.IbanNumber, BMorf.Client.IbanNumber, amount)
		if err != nil {
			errChan <- err
		}
		log.Info(fmt.Sprintf("Gorutine 1: Ember(%f)___Morf(%f)", AEmber.AmountOfMoney.Value, BMorf.AmountOfMoney.Value))
	}()

	go func() {
		amount := model.Amount{
			Value: 9700,
		}
		log.Info(fmt.Sprintf("Gorutine 2: Nikita(%f)--%f-->Ember(%f)", ANikita.AmountOfMoney.Value, amount.Value, AEmber.AmountOfMoney.Value))
		defer wg.Done()
		err := visa.Transfer(ANikita.Client.IbanNumber, AEmber.Client.IbanNumber, amount)
		if err != nil {
			errChan <- err
		}
		log.Info(fmt.Sprintf("Gorutine 2: Nikita(%f)___Ember(%f)", ANikita.AmountOfMoney.Value, AEmber.AmountOfMoney.Value))
	}()

	go func() {
		defer wg.Done()
		amount := model.Amount{
			Value: 9700,
		}
		log.Info(fmt.Sprintf("Gorutine 3: Egor(%f)--%f-->Nikita(%f)", BEgor.AmountOfMoney.Value, amount.Value, ANikita.AmountOfMoney.Value))
		err := visa.Transfer(BEgor.Client.IbanNumber, ANikita.Client.IbanNumber, amount)
		if err != nil {
			errChan <- err
		}
		log.Info(fmt.Sprintf("Gorutine 3: Egor(%f)___Mikita(%f)", BEgor.AmountOfMoney.Value, ANikita.AmountOfMoney.Value))
	}()

	go func() {
		defer wg.Done()
		amount := model.Amount{
			Value: 9700,
		}
		log.Info(fmt.Sprintf("Gorutine 4: Morf(%f)--%f-->Egor(%f)", BMorf.AmountOfMoney.Value, amount.Value, BEgor.AmountOfMoney.Value))
		err := visa.Transfer(BMorf.Client.IbanNumber, BEgor.Client.IbanNumber, amount)
		if err != nil {
			errChan <- err
		}
		log.Info(fmt.Sprintf("Gorutine 4: Morf(%f)___Egor(%f)", BMorf.AmountOfMoney.Value, BEgor.AmountOfMoney.Value))
	}()

	select {
	case err := <-errChan:
		if err != nil {
			log.Error(err)
		}
	}

	wg.Wait()
}
