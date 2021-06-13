package pgstore

import (
	"fmt"
	"github.com/sergazyyev/wallet/app/errs"
	"github.com/sergazyyev/wallet/app/models"
	"github.com/sergazyyev/wallet/app/providers/db"
)

const (
	lockTimeOut     = `3s`
	debitOperation  = `D`
	creditOperation = `C`
)

type walletBalance struct {
	ID      int64   `db:"id"`
	Code    string  `db:"code"`
	Balance float64 `json:"balance"`
}

type operationResult struct {
	transaction   *models.Transaction
	actualBalance float64
}

type walletRepo struct {
	conn db.SQLExecutor
}

//Create creates wallet
func (w *walletRepo) Create(wallet *models.Wallet) (err error) {
	wallet.CurrencyID, err = w.CurrencyIDByCode(nil, wallet.CurrencyCode)
	if err != nil {
		return
	}
	if err = w.conn.QueryRowx(`insert into wallet(code, currency_id, cli_name) values ($1, $2, $3) returning id, balance, create_at`,
		wallet.Code, wallet.CurrencyID, wallet.ClientName).Scan(&wallet.ID, &wallet.Balance, &wallet.CreateAt);
		err != nil {
		return
	}
	return err
}

//AddAmount adds amount of money to wallet
func (w *walletRepo) AddAmount(walletID string, amount float64) (balance float64, err error) {
	var (
		tx  db.Transactional
		bal *walletBalance
		res *operationResult
	)
	tx, err = w.conn.Start()
	if err != nil {
		return
	}

	//lock row and retrieve balance
	bal, err = w.getBalanceWithLock(tx, walletID)
	if err != nil {
		err = w.error(tx.Rollback, err)
		return
	}
	res, err = w.addBalanceAndLog(tx, bal, amount)
	if err != nil {
		err = w.error(tx.Rollback, err)
		return
	}
	err = w.error(tx.Commit, err)
	balance = res.actualBalance
	return
}

//Transfer transfers money between two wallets
func (w *walletRepo) Transfer(walletFrom, walletTo string, amount float64) (err error) {
	var (
		tx             db.Transactional
		balFrom, balTo *walletBalance
	)
	tx, err = w.conn.Start()
	if err != nil {
		return
	}

	// debit
	balFrom, err = w.getBalanceWithLock(tx, walletFrom)
	if err != nil {
		return w.error(tx.Rollback, err)
	}
	if balFrom.Balance < amount {
		return w.error(tx.Rollback, errs.New(errs.NotEnoughBalance))
	}
	_, err = w.minusBalanceAndLog(tx, balFrom, amount)
	if err != nil {
		return w.error(tx.Rollback, err)
	}

	//credit
	balTo, err = w.getBalanceWithLock(tx, walletTo)
	if err != nil {
		return w.error(tx.Rollback, err)
	}
	_, err = w.addBalanceAndLog(tx, balTo, amount)
	if err != nil {
		return w.error(tx.Rollback, err)
	}
	err = w.error(tx.Commit, err)
	return
}

//CurrencyIDByCode interface implementation
func (w *walletRepo) CurrencyIDByCode(tx db.Transactional, code string) (id int64, err error) {
	if tx != nil {
		err = tx.QueryRowx(`select id from currencies where code = $1 and archive_fl = false`, code).Scan(&id)
		return
	}
	err = w.conn.QueryRowx(`select id from currencies where code = $1 and archive_fl = false`, code).Scan(&id)
	return
}

//WalletByCode interface implementation
func (w *walletRepo) WalletByCode(tx db.Transactional, code string) (res *models.Wallet, err error) {
	res = new(models.Wallet)
	if tx != nil {
		err = tx.QueryRowx(`select w.id, w.code, w.currency_id,  c.code as currency_code, w.cli_name, w.balance,  w.create_at from wallet w join currencies c on c.id = w.currency_id where w.code = $1`,
			code).StructScan(res)
		return
	}
	err = w.conn.QueryRowx(`select w.id, w.code, w.currency_id,  c.code as currency_code, w.cli_name, w.balance,  w.create_at from wallet w join currencies c on c.id = w.currency_id where w.code = $1`,
		code).StructScan(res)
	return
}

func (w *walletRepo) addBalanceAndLog(tx db.Transactional, wallet *walletBalance, amount float64) (result *operationResult, err error) {
	result = new(operationResult)
	result.actualBalance, err = w.addBalance(tx, wallet.Code, amount)
	if err != nil {
		return
	}
	result.transaction = &models.Transaction{
		WalletID:      wallet.ID,
		Amount:        amount,
		OperationType: creditOperation,
	}
	err = w.addTransactionLog(tx, result.transaction)
	return
}

func (w *walletRepo) minusBalanceAndLog(tx db.Transactional, wallet *walletBalance, amount float64) (result *operationResult, err error) {
	result = new(operationResult)
	result.actualBalance, err = w.minusBalance(tx, wallet.Code, amount)
	if err != nil {
		return
	}
	result.transaction = &models.Transaction{
		WalletID:      wallet.ID,
		Amount:        amount,
		OperationType: debitOperation,
	}
	err = w.addTransactionLog(tx, result.transaction)
	return
}

func (w *walletRepo) addBalance(tx db.Transactional, walletID string, amount float64) (balance float64, err error) {
	if err = tx.QueryRowx(`update wallet set balance = balance + $1 where code = $2 returning balance`, amount, walletID).
		Scan(&balance); err != nil {
		return
	}
	return
}

func (w *walletRepo) minusBalance(tx db.Transactional, walletID string, amount float64) (balance float64, err error) {
	if err = tx.QueryRowx(`update wallet set balance = balance - $1 where code = $2 returning balance`, amount, walletID).
		Scan(&balance); err != nil {
		return
	}
	return
}

//retrieve and lock row for get actual balance
func (w *walletRepo) getBalanceWithLock(tx db.Transactional, walletID string) (balance *walletBalance, err error) {
	_, err = tx.Exec(fmt.Sprintf(`SET LOCAL lock_timeout = '%s'`, lockTimeOut))
	if err != nil {
		return
	}
	balance = new(walletBalance)
	if err = tx.QueryRowx(`select id, code, balance from wallet where code = $1 for update`, walletID).StructScan(balance);
		err != nil {
		return
	}
	return
}

func (w *walletRepo) addTransactionLog(tx db.Transactional, transaction *models.Transaction) (err error) {
	if err = tx.QueryRowx(`insert into transactions (wallet_id, amount, operation_type) values ($1, $2, $3) returning id, amount, date`,
		transaction.WalletID, transaction.Amount, transaction.OperationType).Scan(&transaction.ID, &transaction.Amount, &transaction.Date);
		err != nil {
		return
	}
	return
}

func (w *walletRepo) error(fn func() error, origin error) error {
	if err := fn(); err != nil {
		return err
	}
	return origin
}
