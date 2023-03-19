package dictionary

import (
	"errors"
	"math"
)

type TenorLimit struct {
	Id					int64		`json:"id"`
	CustomerId	int64		`json:"customer_id"`
	LimitValue	int64		`json:"limit_value"`
	MonthTenor	int64		`json:"month_tenor"`
}

func (tL *TenorLimit) MakeLimit(salary int64) error {
	if tL.MonthTenor > 4 {
		return errors.New("month tenor should not be more than 4")
	}

	for i := 0; i < 4; i++ {
		limit_salary := 1000000 * math.Pow(10, float64(i))
		ratio := 0.10 + (0.04 * float64(i))
		temp := false
		for j := 1; j < 5; j++ {
			if tL.MonthTenor == int64(j) && salary < int64(limit_salary) {
				temp = true
				break
			}
			ratio += 0.01
		}
		if temp { 
			tL.LimitValue = int64(float64(salary) * ratio)
			break 
		}
	}
	return nil
}
