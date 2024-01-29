package model

type SystemPayer interface {
	RegisterBank(bankPrefix string, bank Banker)
	Transfer(ibanA string, ibanB string, amount Amount) error
}

type Visa struct {
	RegisteredBanks map[string]Banker
}

type Payment struct {
	Sender    *Account
	Recipient *Account
	Amount    Amount
}

type Amount struct {
	Value float64
}

func (ps *Visa) RegisterBank(bankPrefix string, bank Banker) {
	ps.RegisteredBanks[bankPrefix] = bank
}

func (ps *Visa) Transfer(ibanA string, ibanB string, amount Amount) error {
	ibanABankPrefix := ibanA[:8]
	ibanBBankPrefix := ibanB[:8]
	clientA := ps.RegisteredBanks[ibanABankPrefix].GetAccount(ibanA)
	clientB := ps.RegisteredBanks[ibanBBankPrefix].GetAccount(ibanB)
	if err := clientA.Withdraw(amount); err != nil {
		return err
	}
	clientB.Deposit(amount)
	return nil
}
