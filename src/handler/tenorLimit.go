package handler

import (
	"net/http"
	"errors"
	"strconv"

	"github.com/sfa119f/backend_xyz/src/dictionary"
	"github.com/sfa119f/backend_xyz/src/service"
	"github.com/sfa119f/backend_xyz/src/utils"
)

func MakeTenorLimit(w http.ResponseWriter, r *http.Request, salary int64, month int64) bool {
	custId := utils.GetIdCustomerInfoCtx(w, r)
	if custId == 0 {
		utils.JsonResp(w, 401, nil, errors.New(dictionary.UnauthorizedError))
		return false
	}
	if month < 0 || month > 12  {
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

func GetTenorByIdCust(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	strMonthTenor := query.Get("monthTenor")
	if strMonthTenor == "" { strMonthTenor = "0" }
	monthTenor, err := strconv.Atoi(strMonthTenor)
	if err != nil {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidParamError))
		return
	}

	custId := utils.GetIdCustomerInfoCtx(w, r)
	if custId == 0 {
		utils.JsonResp(w, 401, nil, errors.New(dictionary.UnauthorizedError))
		return
	}

	res, err := service.GetTenorByIdCust(custId, int64(monthTenor))
	if err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	utils.JsonResp(w, 200, res, nil)
}
