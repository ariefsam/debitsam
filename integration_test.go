package debitsam_test

import (
	"context"
	"debitsam"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ariefsam/eventsam/client"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

func TestIntegration(t *testing.T) {

	godotenv.Load()
	os.Remove("projection_default.db")
	log.SetFlags(log.LstdFlags | log.Llongfile)
	eventsam, err := client.NewEventsam(os.Getenv("EVENTSAM_URL"))
	assert.NoError(t, err)
	if err != nil {
		return
	}
	debitsamObj, err := debitsam.NewDebitsam(&eventsam)
	assert.NoError(t, err)
	var minBalance float64
	err = debitsamObj.Debit("internal", "w001", 10000, &minBalance)
	assert.Error(t, err)

	err = debitsamObj.CreateWallet("internal", "wallet internal")
	assert.NoError(t, err)

	err = debitsamObj.Debit("internal", "w001", 10000, &minBalance)
	assert.Error(t, err)

	minBalance = 10000
	err = debitsamObj.Debit("internal", "w001", 10000, &minBalance)
	assert.Error(t, err, "Need error w001 not found")

	err = debitsamObj.CreateWallet("w001", "wallet w001")
	assert.NoError(t, err)

	err = debitsamObj.Debit("internal", "w001", 10000, &minBalance)
	assert.NoError(t, err)

	wallet, err := debitsamObj.GetWallet("internal")
	assert.NoError(t, err)
	assert.Equal(t, float64(10000), wallet.DebitBalance)

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	debitsam.CreditProjection(ctx, &eventsam)

	wallet, err = debitsam.GetWallet("w001")
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	if wallet == nil {
		return
	}
	assert.Equal(t, float64(10000), wallet.CreditBalance)

}

func TestProjection(t *testing.T) {
	os.Remove("projection_default.db")
	log.SetFlags(log.LstdFlags | log.Llongfile)
	eventsam, err := client.NewEventsam(os.Getenv("EVENTSAM_URL"))
	assert.NoError(t, err)

	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	debitsam.CreditProjection(ctx, &eventsam)

	wallet, err := debitsam.GetWallet("w001")
	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	if wallet == nil {
		return
	}
	assert.Equal(t, float64(10000), wallet.CreditBalance)
}
