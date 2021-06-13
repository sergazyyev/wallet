package handlers

import (
	"github.com/gorilla/mux"
	"github.com/sergazyyev/wallet/app/services"
	"net/http"
)

//Server app's server
type Server struct {
	mux    *mux.Router
	wallet services.WalletService
	report services.ReportService
}

//NewServer constructor
func NewServer(mux *mux.Router, wallet services.WalletService, report services.ReportService) *Server {
	return &Server{
		mux:    mux,
		wallet: wallet,
		report: report,
	}
}

//RegisterHandlers registers endpoints
func (s *Server) RegisterHandlers() {
	s.mux.HandleFunc(`/wallet`, s.handleCreateWallet()).Methods(http.MethodPost)
	s.mux.HandleFunc(`/wallet/deposit`, s.handleDepositWallet()).Methods(http.MethodPut)
	s.mux.HandleFunc(`/wallet/transfer`, s.handleTransferBetweenWallets()).Methods(http.MethodPut)
	s.mux.HandleFunc(`/report/transactions`, s.handleReportWalletTransactions()).Methods(http.MethodGet)
}

//Stop tries to gracefully shutdown server
func (s *Server) Stop() {
}

//ServeHTTP implement http.Handler interface
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
