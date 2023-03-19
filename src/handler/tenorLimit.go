package handler

import (
	"net/http"
	"errors"

	"github.com/sfa119f/backend_xyz/src/dictionary"
	"github.com/sfa119f/backend_xyz/src/service"
	"github.com/sfa119f/backend_xyz/src/utils"
)

func MakeTenorLimit(w http.ResponseWriter, r *http.Request, salary int64, month int64) bool {
	custId := utils.GetIdCustomerInfoCtx(w, r)
	if custId == 0 || month < 0 || month > 12  {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidRequestError))
		return false
	}

	tL := dictionary.TenorLimit{ CustomerId: custId, MonthTenor: month }
	if err := tL.MakeLimit(salary); err != nil {
		utils.JsonResp(w, 400, nil, err)
		return false
	}

	if err := service.InsertUpdateTenor(tL); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return false
	}
	return true
}
