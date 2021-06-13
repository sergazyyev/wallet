package models

import (
	"github.com/sergazyyev/wallet/pkg/convert"
	"time"
)

//Transaction model for transaction log
type Transaction struct {
	ID            int64     `json:"id" db:"id"`
	WalletID      int64     `json:"-" db:"wallet_id"`
	WalletName    string    `json:"wallet_name" db:"wallet_code"`
	Date          time.Time `json:"date" db:"date"`
	Amount        float64   `json:"amount" db:"amount"`
	OperationType string    `json:"operation_type" db:"operation_type"`
}

//CsvRecord record for csv
func (t *Transaction) CsvRecord() []string {
	return []string{
		convert.String(t.ID),
		t.WalletName,
		t.Date.Format(time.RFC3339),
		convert.String(t.Amount),
		t.OperationType,
	}
}

//Transactions wrapper type
type Transactions []*Transaction

//CsvRecords return csv records from transactions history
func (trs Transactions) CsvRecords() [][]string {
	res := [][]string{
		{"id", "wallet_name", "date", "amount", "operation_type"},
	}
	for _, tr := range trs {
		res = append(res, tr.CsvRecord())
	}
	return res
}
