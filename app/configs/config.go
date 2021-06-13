package configs

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

//AppCfg app's configuration
type AppCfg struct {
	BindAddr string `env:"APP_PORT,required"`
	DB       *PostgresCfg
}

//PostgresCfg db configuration
type PostgresCfg struct {
	Host          string `env:"DB_HOST,required"`
	Port          int    `env:"DB_PORT,required"`
	DbName        string `env:"DB_NAME,required"`
	SSLMode       string `env:"DB_SSLMODE,required"`
	Username      string `env:"DB_USERNAME,required"`
	Password      string `env:"DB_PASSWORD,required"`
	MaxIdleConns  int    `env:"DB_MAX_IDLE_CONNS,required"`
	MaxOpenConns  int    `env:"DB_MAX_OPEN_CONNS,required"`
	IsAutoMigrate bool   `env:"IS_AUTOMIGRATE,required"`
}

//ConnectionString postgres connection string
func (p *PostgresCfg) ConnectionString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s sslmode=%s user=%s password=%s",
		p.Host,
		p.Port,
		p.DbName,
		p.SSLMode,
		p.Username,
		p.Password)
}

//NewConfig returns app config
func NewConfig() (*AppCfg, error) {
	cfg := &AppCfg{DB: &PostgresCfg{}}
	if err := env.Parse(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
