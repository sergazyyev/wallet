package db

import (
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

var (
	defaultMaxConns  = 100
	defaultIdleConns = 3
)

//Connection sqlx wrapper
type Connection struct {
	*sqlx.DB
}

//Option function for set options
type Option func(*Connection)

//New create new connection to postgresql db
func New(connectionStr string, options ...Option) (SQLExecutor, error) {
	db, err := sqlx.Open("pgx", connectionStr)
	if err != nil {
		return nil, err
	}
	conn := &Connection{db}
	conn.SetMaxOpenConns(defaultMaxConns)
	conn.SetMaxIdleConns(defaultIdleConns)
	for _, op := range options {
		op(conn)
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	return conn, err
}

//WithMaxOpenConns sets max opening connections for db
func WithMaxOpenConns(n int) Option {
	return func(c *Connection) {
		c.SetMaxOpenConns(n)
	}
}

//WithMaxIdleConns sets max idle connections for db
func WithMaxIdleConns(n int) Option {
	return func(c *Connection) {
		c.SetMaxIdleConns(n)
	}
}

//Start implement StartTransaction interface
func (d *Connection) Start() (Transactional, error) {
	return d.Beginx()
}

//MigrateUP migrate db schema up
func (d *Connection) MigrateUP() error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	_, err := migrate.Exec(d.DB.DB, "postgres", migrations, migrate.Up)
	if err != nil {
		return err
	}
	return nil
}

//MigrateDown migrate db schema down
func (d *Connection) MigrateDown() error {
	migrations := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	_, err := migrate.Exec(d.DB.DB, "postgres", migrations, migrate.Down)
	if err != nil {
		return err
	}
	return nil
}
