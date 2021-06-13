package app

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sergazyyev/wallet/app/configs"
	"github.com/sergazyyev/wallet/app/env"
	"github.com/sergazyyev/wallet/app/handlers"
	"github.com/sergazyyev/wallet/app/providers/db"
	"github.com/sergazyyev/wallet/app/services"
	"github.com/sergazyyev/wallet/app/store/pgstore"
	"github.com/sergazyyev/wallet/app/validators"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

//Start configure and starts app
func Start() (err error) {
	err = env.Load()
	if err != nil {
		return fmt.Errorf(`error when load env from file, err: %w`, err)
	}
	cfg, err := configs.NewConfig()
	if err != nil {
		return fmt.Errorf(`error when parse config, err: %w`, err)
	}
	sqlExec, err := db.New(cfg.DB.ConnectionString(), db.WithMaxIdleConns(cfg.DB.MaxIdleConns), db.WithMaxOpenConns(cfg.DB.MaxOpenConns))
	if err != nil {
		return fmt.Errorf(`error configure sql provider, err: %w`, err)
	}
	defer sqlExec.Close()
	if cfg.DB.IsAutoMigrate {
		if err := sqlExec.MigrateUP(); err != nil {
			return fmt.Errorf(`error when migrate db schema UP, err: %w`, err)
		}
	}
	st := pgstore.New(sqlExec)
	validators.InitCustomValidators(st)
	logger := logrus.New()
	walletSrv := services.NewWalletService(st, logger)
	reportSrv := services.NewReportService(st, logger)
	server := handlers.NewServer(mux.NewRouter(), walletSrv, reportSrv)
	server.RegisterHandlers()

	errCh := make(chan error, 2)
	go func() {
		logger.Infof(`started server on %s`, cfg.BindAddr)
		if err = http.ListenAndServe(cfg.BindAddr, server); err != nil {
			errCh <- err
		}
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT)
		stopSignal := <-c
		server.Stop()
		errCh <- fmt.Errorf("%s", stopSignal)
	}()

	err = <-errCh
	logger.Infof(`stoped app with %s`, err)
	return
}
