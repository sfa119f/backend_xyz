package handler

import (
	"encoding/json"
	"net/http"
	"errors"
	"time"

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
