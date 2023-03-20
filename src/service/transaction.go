package service

import (
	"strconv"

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

func GetTransactionByIdCust(
	custId int64, month int64, year int64, otrMin int64, otrMax int64,
) ([]map[string]interface{}, error) {
	db := database.GetDB()

	qMonth := ""
	if month != 0 {
		qMonth = ` and extract(month from created_at) = ` + strconv.Itoa(int(month))
	}

	qYear := ""
	if month != 0 {
		qYear = ` and extract(year from created_at) = ` + strconv.Itoa(int(year))
	}

	qOtrMin := ` and otr >= ` + strconv.Itoa(int(otrMin))
	qOtrMax := ""
	if otrMax != 0 {
		qOtrMax = ` and otr <= ` + strconv.Itoa(int(otrMax))
	}

	query := 
		`select * from transaction where customer_id = ` + 
		strconv.Itoa(int(custId)) + qMonth + qYear + qOtrMin + qOtrMax

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []map[string]interface{}{}
	for rows.Next() {
		t := dictionary.Transaction{}
		if err := 
		rows.Scan(
			&t.Id, &t.CustomerId, &t.OTR, &t.AdminFee, &t.InstallmentAmount, 
			&t.InterestAmount, &t.Assetname, &t.CreatedAt,
		); err != nil {
			return res, err
		}
		res = append(
			res,
			map[string]interface{}{
				"otr": t.OTR, "admin_fee": t.AdminFee, "installment_amount": t.InstallmentAmount, 
				"interest_amount": t.InterestAmount, "assetname": t.Assetname, "created_at": t.CreatedAt,
			},
		)
	}
	return res, nil
}
