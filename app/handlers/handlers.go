package handlers

import (
	"github.com/sergazyyev/wallet/app/handlers/requests"
	"github.com/sergazyyev/wallet/app/models"
	"net/http"
)

func (s *Server) handleCreateWallet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(requests.CreateWalletRequest)
		err := s.validateRequest(r, req)
		if err != nil {
			s.errorResponse(w, err)
			return
		}
		model := &models.Wallet{
			ClientName:   req.ClientName,
			CurrencyCode: req.CurrencyCode,
		}
		err = s.wallet.Create(r.Context(), model)
		if err != nil {
			s.errorResponse(w, err)
			return
		}
		s.response(w, model, "")
	}
}

func (s *Server) handleDepositWallet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(requests.DepositWalletRequest)
		err := s.validateRequest(r, req)
		if err != nil {
			s.errorResponse(w, err)
			return
		}
		bal, err := s.wallet.Deposit(r.Context(), req.WalletName, req.Amount)
		if err != nil {
			s.errorResponse(w, err)
			return
		}
		s.response(w, map[string]float64{"balance": bal}, "")
	}
}

func (s *Server) handleTransferBetweenWallets() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := new(requests.TransferInWalletsRequests)
		err := s.validateRequest(r, req)
		if err != nil {
			s.errorResponse(w, err)
			return
		}
		err = s.wallet.Transfer(r.Context(), req.FromWallet, req.ToWallet, req.Amount)
		if err != nil {
			s.errorResponse(w, err)
			return
		}
		s.response(w, map[string]string{"message": "ok"}, "")
	}
}

func (s *Server) handleReportWalletTransactions() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &requests.ReportTransactionHistoryRequest{
			DateFrom:      r.URL.Query().Get("date_from"),
			DateTo:        r.URL.Query().Get("date_to"),
			OperationType: r.URL.Query().Get("operation_type"),
			ResponseType:  r.URL.Query().Get("response_type"),
		}
		err := s.validateRequest(r, req)
		if err != nil {
			s.errorResponse(w, err)
			return
		}
		res, err := s.report.TransactionHistory(req.GetDateFrom(), req.GetDateTo(), req.OperationType, req.ResponseType)
		if err != nil {
			s.errorResponse(w, err)
			return
		}
		contentType := "application/json"
		if req.ResponseType == "csv" {
			contentType = "text/csv"
		}
		s.response(w, res, contentType)
	}
}
