package requests

//CreateWalletRequest request for create waller
type CreateWalletRequest struct {
	ClientName   string `json:"client_name"`
	CurrencyCode string `json:"currency_code"`
}

//Rules implement validator interface
func (r *CreateWalletRequest) Rules() map[string][]string {
	return map[string][]string{
		"client_name":   {"required"},
		"currency_code": {"required", "support_currency"},
	}
}

//Messages implement validator interface
func (r *CreateWalletRequest) Messages() map[string][]string {
	return map[string][]string{
		"client_name": {
			"required:client_name is required",
		},
		"currency_code": {
			"required:currency_code is required",
			"support_currency:unsupported currency",
		},
	}
}
