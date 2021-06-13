package requests

//DepositWalletRequest deposit request
type DepositWalletRequest struct {
	WalletName string  `json:"wallet"`
	Amount     float64 `json:"amount"`
}

//Rules implement validator interface
func (r *DepositWalletRequest) Rules() map[string][]string {
	return map[string][]string{
		"wallet": {"required", "wallet_exists"},
		"amount": {"required", "min:0"},
	}
}

//Messages implement validator interface
func (r *DepositWalletRequest) Messages() map[string][]string {
	return map[string][]string{
		"wallet": {
			"required:wallet is required",
			"wallet_exists:wallet doesnt exists",
		},
		"amount": {
			"required:amount is required",
			"min:amount must be > 0",
		},
	}
}
