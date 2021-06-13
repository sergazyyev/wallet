package requests

import (
	"github.com/sergazyyev/wallet/app/validators"
	"time"
)

//ReportTransactionHistoryRequest request struct
type ReportTransactionHistoryRequest struct {
	DateFrom      string `json:"date_from"`
	DateTo        string `json:"date_to"`
	OperationType string `json:"operation_type"`
	ResponseType  string `json:"response_type"`
}

//Rules implement validator interface
func (r *ReportTransactionHistoryRequest) Rules() map[string][]string {
	return map[string][]string{
		"date_from":      {"required", "date"},
		"date_to":        {"required", "date"},
		"operation_type": {"required", "in:C,D"},
		"response_type":  {"required", "in:json,csv"},
	}
}

//Messages implement validator interface
func (r *ReportTransactionHistoryRequest) Messages() map[string][]string {
	return map[string][]string{
		"date_from": {
			"required:date_from is required",
			"date:date_from format must be YYYY-MM-DD",
		},
		"date_to": {
			"required:date_to is required",
			"date:date_to format must be YYYY-MM-DD",
		},
		"operation_type": {
			"required:operation_type is required",
			"in:operation_type must be C (credit) or D (debit)",
		},
		"response_type": {
			"required:response_type is required",
			"in:response_type must be json or csv",
		},
	}
}

//GetDateFrom converts request string date to time.Time
func (r *ReportTransactionHistoryRequest) GetDateFrom() (date time.Time) {
	date, _ = time.Parse(validators.DateFormat, r.DateFrom)
	return
}

//GetDateTo converts request string date to time.Time
func (r *ReportTransactionHistoryRequest) GetDateTo() (date time.Time) {
	date, _ = time.Parse(validators.DateFormat, r.DateTo)
	return
}
