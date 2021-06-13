package services

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"github.com/sergazyyev/wallet/app/errs"
	"github.com/sergazyyev/wallet/app/log"
	"github.com/sergazyyev/wallet/app/models"
	"github.com/sergazyyev/wallet/app/store"
	"time"
)

const (
	csvResponse  = `csv`
	jsonResponse = `json`
)

//ReportService service for report functionality
type ReportService interface {
	TransactionHistory(dateFrom, dateTo time.Time, operationType, responseType string) ([]byte, error)
}

type reportSrv struct {
	st  store.Store
	log log.Logger
}

//NewReportService constructor for service
func NewReportService(store store.Store, logger log.Logger) ReportService {
	return &reportSrv{
		st:  store,
		log: logger,
	}
}

//TransactionHistory interface implementation
func (r *reportSrv) TransactionHistory(dateFrom, dateTo time.Time, operationType, responseType string) ([]byte, error) {
	data, err := r.st.ReportRepository().TransactionHistoryByDate(dateFrom, dateTo, operationType)
	if err != nil {
		if !errs.IsCustomErr(err) {
			r.log.Errorf(`error when retrieve data from db for transactions report, err: %v`, err)
			return nil, errs.New(errs.InternalServerError)
		}
		return nil, err
	}
	if len(data) < 1 {
		return nil, errs.New(errs.EmptyReport)
	}

	buffer := new(bytes.Buffer)
	if responseType == csvResponse {
		err = csv.NewWriter(buffer).WriteAll(models.Transactions(data).CsvRecords())
	} else {
		err = json.NewEncoder(buffer).Encode(data)
	}
	if err != nil {
		r.log.Errorf(`error when make %s response for transactions report, err: %v`, responseType, err)
		return nil, errs.New(errs.InternalServerError)
	}
	return buffer.Bytes(), nil
}
