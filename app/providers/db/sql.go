package db

import (
	"database/sql/driver"
	"github.com/jmoiron/sqlx"
	"io"
)

//Selector interface for convenient function Select from sqlx
type Selector interface {
	Select(dest interface{}, query string, args ...interface{}) error
}

//StartTransaction interface create transaction in db
type StartTransaction interface {
	Start() (Transactional, error)
}

//Transactional interface for sqlx.Tx
type Transactional interface {
	sqlx.Execer
	sqlx.Queryer
	Selector
	driver.Tx
}

//SQLExecutor interface for sql commands executor
type SQLExecutor interface {
	io.Closer
	sqlx.Execer
	sqlx.Queryer
	Selector
	StartTransaction
	Migrator
}

//Migrator interface for sql db
type Migrator interface {
	MigrateUP() error
	MigrateDown() error
}
