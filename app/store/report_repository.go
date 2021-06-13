package store

import (
	"github.com/sergazyyev/wallet/app/models"
	"time"
)

//ReportRepository interface for report repository
type ReportRepository interface {
	TransactionHistoryByDate(dateFrom, dateTo time.Time, operationType string) ([]*models.Transaction, error)
}
