package pgstore

import (
	"github.com/sergazyyev/wallet/app/providers/db"
	"github.com/sergazyyev/wallet/app/store"
)

type pgStore struct {
	conn       db.SQLExecutor
	walletRepo store.WalletRepository
	reportRepo store.ReportRepository
}

//WalletRepository return postgresql repository for wallet
func (s *pgStore) WalletRepository() store.WalletRepository {
	if s.walletRepo != nil {
		return s.walletRepo
	}
	s.walletRepo = &walletRepo{conn: s.conn}
	return s.walletRepo
}

//ReportRepository return postgresql repository for report
func (s *pgStore) ReportRepository() store.ReportRepository {
	if s.reportRepo != nil {
		return s.reportRepo
	}
	s.reportRepo = &reportRepo{conn: s.conn}
	return s.reportRepo
}

//New return store implementation
func New(executor db.SQLExecutor) store.Store {
	return &pgStore{
		conn: executor,
	}
}
