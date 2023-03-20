package dictionary

type Transaction struct {
	Id								int64		`json:"id"`
	CustomerId				int64		`json:"customer_id"`
	OTR								int64		`json:"otr"`
	AdminFee					int64		`json:"admin_fee"`					// default 5000
	InstallmentAmount	int64		`json:"installment_amount"`	// in Month default 1
	InterestAmount	 	int64		`json:"interest_amount"`		// in Percent(%) default 8
	Assetname					string	`json:"assetname"`
	CreatedAt					string	`json:"created_at"`
}
