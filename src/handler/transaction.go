package handler

import (
	"encoding/json"
	"net/http"
	"errors"
	"time"
	"strconv"

	"github.com/sfa119f/backend_xyz/src/dictionary"
	"github.com/sfa119f/backend_xyz/src/service"
	"github.com/sfa119f/backend_xyz/src/utils"
)

func InsertTransaction(w http.ResponseWriter, r *http.Request) {
	custId := utils.GetIdCustomerInfoCtx(w, r)
	if custId == 0 {
		utils.JsonResp(w, 401, nil, errors.New(dictionary.UnauthorizedError))
		return
	}

	trans := dictionary.Transaction{}
	if err := json.NewDecoder(r.Body).Decode(&trans); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	if trans.Assetname == "" {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidRequestError))
		return
	}

	trans.CustomerId = custId
	if trans.AdminFee == 0 { trans.AdminFee = 5000 }
	if trans.InstallmentAmount == 0 { trans.InstallmentAmount = 1 }
	if trans.InterestAmount == 0 { trans.InterestAmount = 8 }

	// Check valid transaction -> total in the same month not more than tenor limit
	arrTenor, err := service.GetTenorByIdCust(custId, int64(trans.InstallmentAmount))
	if err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	var limitValue int64 = arrTenor[0]["limit_value"]

	currentTime := time.Now()
	year := int64(currentTime.Year())
	month := int64(currentTime.Month())
	totalInMonth, err := service.GetTotalTransaction(custId, year, month)
	if err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}

	if totalInMonth + trans.OTR > limitValue {
		utils.JsonResp(w, 400, nil, errors.New("transaction exceeds the tenor limit"))
		return
	}

	// insert to database
	if err := service.InsertTransaction(trans); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	utils.JsonResp(w, 201, map[string]string{"message": "success"}, nil)
}

func GetTransactionByIdCust(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	strMonth := query.Get("month")
	if strMonth == "" { strMonth = "0" }
	month, err := strconv.Atoi(strMonth)
	if err != nil {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidParamError))
		return
	}

	strYear := query.Get("year")
	if strYear == "" { strYear = "0" }
	year, err := strconv.Atoi(strYear)
	if err != nil {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidParamError))
		return
	}

	strOtrMin := query.Get("otrMin")
	if strOtrMin == "" { strOtrMin = "0" }
	otrMin, err := strconv.Atoi(strOtrMin)
	if err != nil {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidParamError))
		return
	}

	strOtrMax := query.Get("otrMax")
	if strOtrMax == "" { strOtrMax = "0" }
	otrMax, err := strconv.Atoi(strOtrMax)
	if err != nil {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidParamError))
		return
	}

	custId := utils.GetIdCustomerInfoCtx(w, r)
	if custId == 0 {
		utils.JsonResp(w, 401, nil, errors.New(dictionary.UnauthorizedError))
		return
	}

	res, err := service.GetTransactionByIdCust(
		custId, int64(month), int64(year), int64(otrMin), int64(otrMax),
	)
	if err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	utils.JsonResp(w, 200, res, nil)
}
