package views

type (
	BalanceResponse struct {
		AccountNumber int     `json:"account_number"`
		CustomerName  string  `json:"customer_name"`
		Balance       float64 `json:"balance"`
	}

	TransferRequest struct {
		Receiver int     `json:"to_account_number" validate:"required"`
		Amount   float64 `json:"amount" validate:"required,gte=0"`
	}
)
