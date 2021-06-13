package requests

//TransferInWalletsRequests request struct
type TransferInWalletsRequests struct {
	FromWallet string  `json:"from"`
	ToWallet   string  `json:"to"`
	Amount     float64 `json:"amount"`
}

//Rules implement validator interface
func (r *TransferInWalletsRequests) Rules() map[string][]string {
	return map[string][]string{
		"from":   {"required", "wallet_exists"},
		"to":     {"required", "wallet_exists"},
		"amount": {"required", "min:0"},
	}
}

//Messages implement validator interface
func (r *TransferInWalletsRequests) Messages() map[string][]string {
	return map[string][]string{
		"from": {
			"required:from is required",
			"wallet_exists:from wallet doesnt exists",
		},
		"to": {
			"required:to is required",
			"wallet_exists:to wallet doesnt exists",
		},
		"amount": {
			"required:amount is required",
			"min:amount must be > 0",
		},
	}
}
