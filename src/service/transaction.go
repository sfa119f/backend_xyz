package service

import (
	"github.com/sfa119f/backend_xyz/src/database"
	"github.com/sfa119f/backend_xyz/src/dictionary"
)

func InsertTransaction(t dictionary.Transaction) error {
	db := database.GetDB()
	query := `
		insert into transaction 
		(customer_id, otr, admin_fee, installment_amount, interest_amount, assetname)
		values ($1, $2, $3, $4, $5, $6)
	`

	_, err := db.Exec(
		query, t.CustomerId, t.OTR, t.AdminFee, t.InstallmentAmount, t.InterestAmount, t.Assetname,
	)
	return err
}

func GetTotalTransaction(custId int64, year int64, month int64) (int64, error) {
	db := database.GetDB()
	query := `
		select coalesce(sum(otr), 0) as total from transaction 
		where customer_id = $1 and extract(year from created_at) = $2 
			and extract(month from created_at) = $3
	`

	var total int64
	if err := db.QueryRow(query, custId, year, month).Scan(&total);
	err != nil {
		return 0, err
	}
	return total, nil
}
