package store

import (
	"github.com/sergazyyev/wallet/app/models"
	"github.com/sergazyyev/wallet/app/providers/db"
)

//WalletRepository interface wallet repository
type WalletRepository interface {
	Create(wallet *models.Wallet) (err error)
	AddAmount(walletID string, amount float64) (balance float64, err error)
	Transfer(walletFrom, walletTo string, amount float64) (err error)
	CurrencyIDByCode(tx db.Transactional, code string) (id int64, err error)
	WalletByCode(tx db.Transactional, code string) (*models.Wallet, error)
}
