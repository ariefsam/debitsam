package debitsam

import (
	"encoding/json"
	"errors"
	"log"

	"github.com/ariefsam/eventsam"
)

type Wallet struct {
	Name          string  `json:"name"`
	Account       string  `gorm:"unique" json:"account"`
	DebitBalance  float64 `json:"debit_balance"`
	CreditBalance float64 `json:"credit_balance"`
	Version       int64   `json:"version"`
}

func (wallet *Wallet) ApplyWallet(events []eventsam.EventEntity) (Rwallet *Wallet, err error) {
	for _, event := range events {
		switch event.EventName {
		case "wallet_created":
			var walletCreated WalletCreated
			err = json.Unmarshal([]byte(event.Data), &walletCreated)
			if err != nil {
				log.Println(err)
				continue
			}
			wallet.Version = event.Version
		case "debit":
			debit := Debit{}
			err = json.Unmarshal([]byte(event.Data), &debit)
			if err != nil {
				log.Println(err)
				continue
			}
			wallet.DebitBalance += debit.Amount
			wallet.Version = event.Version
		}
	}
	Rwallet = wallet
	return
}

func (d *Debitsam) GetWallet(accountID string) (w *Wallet, err error) {
	walletEvents, err := d.eventsam.Retrieve(accountID, "wallet", 0)
	if len(walletEvents) == 0 {
		err = errors.New("wallet not found")
		return
	}

	wallet, err := new(Wallet).ApplyWallet(walletEvents)
	if err != nil {
		return
	}

	w = wallet

	return
}
