package service

import (
	"github.com/sfa119f/backend_xyz/src/database"
	"github.com/sfa119f/backend_xyz/src/dictionary"
)

func InsertUpdateTenor(tL dictionary.TenorLimit) error {
	db := database.GetDB()
	query := `
		insert into tenor_limit (customer_id, limit_value, month_tenor)
		values ($1, $2, $3) on conflict (customer_id, month_tenor) where customer_id > 0
		do update set limit_value = excluded.limit_value
	`

	_, err := db.Exec(query, tL.CustomerId, tL.LimitValue, tL.MonthTenor)
	return err
}
