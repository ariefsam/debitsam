package example

import (
	"debitsam"
	"os"

	eventsam "github.com/ariefsam/eventsam/client"
)

var ds *debitsam.Debitsam

func InitDebitsam() {
	eventsamURL := os.Getenv("EVENTSAM_URL") // example: "http://localhost:8009"
	es, err := eventsam.NewEventsam(eventsamURL)
	if err != nil {
		return
	}
	ds, err = debitsam.NewDebitsam(&es)
	if err != nil {
		return
	}
}

func CreateAccount() (err error) {
	err = ds.CreateWallet("salary01", "Revenue from Salary")
	if err != nil {
		return
	}
	err = ds.CreateWallet("bank01", "My Commonthwealth Bank")
	if err != nil {
		return
	}
	return
}

func DebitAccount() (err error) {
	// bank -> account debit
	// salary -> account credit
	// so we need to debit salary account and credit to bank account
	err = ds.Debit("salary01", "bank01", 100000, nil)
	if err != nil {
		return
	}
	return
}
