package debitsam

import (
	"errors"
	"time"

	"github.com/ariefsam/eventsam"
)

type Eventsam interface {
	FetchAggregateEvent(aggregateName string, afterID, limit int) (events []eventsam.EventEntity, err error)
	Store(aggregateID string, aggregateName string, eventName string, version int64, data any) (entity eventsam.EventEntity, err error)
	Retrieve(aggregateID string, aggregateName string, afterVersion int) (events []eventsam.EventEntity, err error)
}

type Debit struct {
	SourceAccount      string  `json:"source_account"`
	Amount             float64 `json:"amount"`
	DestinationAccount string  `json:"destination_account"`
	TimeMillis         int64   `json:"time_millis"`
}

func (d *Debitsam) Debit(sourceAccount, destinationAccount string, amountDebit float64, minBalance *float64) (err error) {

	wallet, err := d.GetWallet(sourceAccount)
	if err != nil {
		return
	}

	if wallet == nil {
		err = errors.New("source wallet not found")
		return
	}

	wallet2, err := d.GetWallet(destinationAccount)
	if err != nil {
		return
	}

	if wallet2 == nil {
		err = errors.New("destination wallet not found")
		return
	}
	if minBalance != nil {
		if wallet.CreditBalance-amountDebit < *minBalance {
			err = errors.New("insufficient balance")
			return
		}
	}

	debitData := Debit{
		SourceAccount:      sourceAccount,
		Amount:             amountDebit,
		DestinationAccount: destinationAccount,
		TimeMillis:         time.Now().UnixMilli(),
	}

	_, err = d.eventsam.Store(sourceAccount, "wallet", "debit", wallet.Version+1, debitData)

	return
}

func (d *Debitsam) CreateWallet(account, name string) (err error) {
	walletEvents, err := d.eventsam.Retrieve(account, "wallet", 0)
	if err != nil {
		return
	}

	if len(walletEvents) > 0 {
		err = errors.New("wallet already exist")
		return
	}
	walletCreated := WalletCreated{
		Account: account,
		Name:    name,
	}

	_, err = d.eventsam.Store(account, "wallet", "wallet_created", 1, walletCreated)
	return
}

type WalletCreated struct {
	Account string `json:"account"`
	Name    string `json:"name"`
}
