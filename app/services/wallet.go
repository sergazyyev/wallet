package services

import (
	"context"
	"github.com/google/uuid"
	"github.com/sergazyyev/wallet/app/errs"
	"github.com/sergazyyev/wallet/app/log"
	"github.com/sergazyyev/wallet/app/models"
	"github.com/sergazyyev/wallet/app/store"
	"math"
)

//WalletService service for wallet functionality
type WalletService interface {
	Create(ctx context.Context, wallet *models.Wallet) error
	Deposit(ctx context.Context, walletID string, amount float64) (balance float64, err error)
	Transfer(ctx context.Context, from, to string, amount float64) error
}

type walletSrv struct {
	st  store.Store
	log log.Logger
}

//NewWalletService constructor for walletSrv
func NewWalletService(store store.Store, logger log.Logger) WalletService {
	return &walletSrv{
		st:  store,
		log: logger,
	}
}

//Create interface implementation
func (w *walletSrv) Create(ctx context.Context, wallet *models.Wallet) error {
	wallet.Code = uuid.New().String()
	if err := w.st.WalletRepository().Create(wallet); err != nil {
		if !errs.IsCustomErr(err) {
			w.log.Errorf(`error when create wallet in db, err: %v`, err)
			return errs.New(errs.InternalServerError)
		}
		return err
	}
	return nil
}

//Deposit interface implementation
func (w *walletSrv) Deposit(ctx context.Context, walletID string, amount float64) (balance float64, err error) {
	amount = math.Round(amount*100) / 100
	balance, err = w.st.WalletRepository().AddAmount(walletID, amount)
	if err != nil {
		if !errs.IsCustomErr(err) {
			w.log.Errorf(`error when deposit wallet in db, err: %v`, err)
			return 0, errs.New(errs.InternalServerError)
		}
		return
	}
	return
}

//Transfer interface implementation
func (w *walletSrv) Transfer(ctx context.Context, from, to string, amount float64) error {
	amount = math.Round(amount*100) / 100
	if err := w.st.WalletRepository().Transfer(from, to, amount); err != nil {
		if !errs.IsCustomErr(err) {
			w.log.Errorf(`error when transfer between wallets in db, err: %v`, err)
			return errs.New(errs.InternalServerError)
		}
		return err
	}
	return nil
}
