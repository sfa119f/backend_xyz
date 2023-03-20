package service

import (
	"strconv"

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

func GetTenorByIdCust(custId int64, monthTenor int64) ([]map[string]int64, error) {
	db := database.GetDB()

	qMonthTenor := ``
	if monthTenor != 0 {
		qMonthTenor = ` and month_tenor = ` + strconv.Itoa(int(monthTenor))
	}
	query := 
		`select * from tenor_limit where customer_id = ` + strconv.Itoa(int(custId)) + qMonthTenor

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	res := []map[string]int64{}
	for rows.Next() {
		tL := dictionary.TenorLimit{}
		if err := rows.Scan(&tL.Id, &tL.CustomerId, &tL.LimitValue, &tL.MonthTenor); err != nil {
			return res, err
		}
		res = append(
			res,
			map[string]int64{"limit_value": tL.LimitValue, "month_tenor": tL.MonthTenor},
		)
	}
	return res, nil
}
