package validators

import (
	"database/sql"
	"errors"
	"github.com/sergazyyev/wallet/app/store"
	"github.com/thedevsaddam/govalidator"
)

const (
	//DateFormat for report requests
	DateFormat = `2006-01-02`
)

type customValidator struct {
	st store.Store
}

//InitCustomValidators register custom
//validators in govalidator
func InitCustomValidators(store store.Store) {
	validator := &customValidator{st: store}
	validator.addSupportCurrencyValidator()
	validator.addWalletExistsValidator()
}

func (v *customValidator) addSupportCurrencyValidator() {
	govalidator.AddCustomRule("support_currency", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		id, err := v.st.WalletRepository().CurrencyIDByCode(nil, val)
		if id == 0 || err == sql.ErrNoRows {
			return errors.New(message)
		}
		return nil
	})
}

func (v *customValidator) addWalletExistsValidator() {
	govalidator.AddCustomRule("wallet_exists", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		wallet, err := v.st.WalletRepository().WalletByCode(nil, val)
		if wallet == nil || err == sql.ErrNoRows {
			return errors.New(message)
		}
		return nil
	})
}
