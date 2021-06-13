package pgstore

import (
	"github.com/sergazyyev/wallet/app/models"
	"github.com/sergazyyev/wallet/app/providers/db"
	"time"
)

type reportRepo struct {
	conn db.SQLExecutor
}

func (r *reportRepo) TransactionHistoryByDate(dateFrom, dateTo time.Time, operationType string) (res []*models.Transaction, err error) {
	if err = r.conn.Select(&res, `select t.id, t.wallet_id, w.code as wallet_code, t.date, t.amount, t.operation_type from transactions t join wallet w on w.id = t.wallet_id where t.date::date between $1::date and $2::date and operation_type = $3 order by date desc`,
		dateFrom, dateTo, operationType); err != nil {
		return
	}
	return
}
