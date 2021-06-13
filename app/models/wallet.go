package models

import "time"

//Wallet model
type Wallet struct {
	ID           int64     `json:"-" db:"id"`
	Code         string    `json:"name" db:"code"`
	CurrencyID   int64     `json:"-" db:"currency_id"`
	CurrencyCode string    `json:"currency_code" db:"currency_code"`
	ClientName   string    `json:"client_name" db:"cli_name"`
	Balance      float64   `json:"balance" db:"balance"`
	CreateAt     time.Time `json:"create_at" db:"create_at"`
}
