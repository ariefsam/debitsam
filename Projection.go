package debitsam

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var dbProjection *gorm.DB

func Init() {
	var err error
	filepath := os.Getenv("DB_PROJECTION_FILEPATH")
	if filepath == "" {
		filepath = "projection_default.db"
	}
	dbProjection, err = gorm.Open(sqlite.Open(filepath), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	dbProjection.AutoMigrate(&WalletGorm{})
	dbProjection.AutoMigrate(&Cursor{})
}

type Cursor struct {
	gorm.Model
	Cursor uint
}

type WalletGorm struct {
	gorm.Model
	Wallet
}

func CreditProjection(ctx context.Context, eventsam Eventsam) {
	Init()
	for {
		if ctx.Err() != nil {
			return
		}
		events, err := eventsam.FetchAggregateEvent("wallet", 0, 100)

		if err != nil {
			time.Sleep(10 * time.Second)
			log.Println(err)
			continue
		}

		if err != nil {
			time.Sleep(10 * time.Second)
			log.Println(err)
			continue
		}
		for _, event := range events {
			switch event.EventName {
			case "wallet_created":
				walletCreated := WalletCreated{}

				err = json.Unmarshal([]byte(event.Data), &walletCreated)
				if err != nil {
					log.Println(err)
				}

				wallet := WalletGorm{}
				wallet.Name = walletCreated.Name
				wallet.Account = walletCreated.Account
				if err = dbProjection.Create(&wallet).Error; err != nil {
					log.Println(err)
				}

			case "debit":
				debit := Debit{}
				err = json.Unmarshal([]byte(event.Data), &debit)
				if err != nil {
					log.Println(err)
				}
				gormWallet := WalletGorm{}
				if err = dbProjection.Where("account = ?", debit.SourceAccount).First(&gormWallet).Error; err != nil {
					log.Println(err)
				}
				gormWallet.DebitBalance += debit.Amount
				if err = dbProjection.Save(&gormWallet).Error; err != nil {
					log.Println(err)
				}

				gormWallet = WalletGorm{}
				if err = dbProjection.Where("account = ?", debit.DestinationAccount).First(&gormWallet).Error; err != nil {
					log.Println(err)
				}
				gormWallet.CreditBalance += debit.Amount
				if err = dbProjection.Save(&gormWallet).Error; err != nil {
					log.Println(err)
				}
			}
			cursor := Cursor{}
			cursor.Cursor = event.ID
			if err = dbProjection.Create(&cursor).Error; err != nil {
				log.Println(err)
			}

		}

		time.Sleep(1 * time.Second)
	}
}

func GetWallet(account string) (wallet *Wallet, err error) {
	gormWallet := WalletGorm{}
	if err = dbProjection.Where("account = ?", account).First(&gormWallet).Error; err != nil {
		log.Println(err)
		return
	}
	wallet = &gormWallet.Wallet
	return
}
