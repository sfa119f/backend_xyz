package handler

import (
	"encoding/json"
	"net/http"
	"errors"

	"github.com/sfa119f/backend_xyz/src/dictionary"
	"github.com/sfa119f/backend_xyz/src/service"
	"github.com/sfa119f/backend_xyz/src/utils"
)

func InsertCustomer(w http.ResponseWriter, r *http.Request) {
	customer := dictionary.Customer{}
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}

	if customer.Fullname == "" || customer.Pass == "" || customer.Email == "" {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidRequestError))
		return
	}

	if err := customer.HashPassword(customer.Pass); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}

	id, err := service.InsertCustomer(customer)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "customers_email_key"` {
			utils.JsonResp(w, 400, nil, errors.New("email already registered"))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}

	customer.Id = id
	signedToken, err := utils.MakeToken(customer)
	if err != nil {
		utils.JsonResp(w, 400, nil, err)
		return
	}

	// Success
	utils.JsonResp(w, 200, map[string]interface{}{"token": signedToken}, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	customer := dictionary.Customer{}
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	
	// Check database
	resDB, err := service.Login(0, customer.Email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.JsonResp(w, 400, nil, errors.New("invalid email or password"))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}
	
	if err := resDB.CheckPassword(customer.Pass); err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			utils.JsonResp(w, 400, nil, errors.New("invalid email or password"))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}

	signedToken, err := utils.MakeToken(resDB)
	if err != nil {
		utils.JsonResp(w, 400, nil, err)
		return
	}

	// Success
	utils.JsonResp(w, 200, map[string]interface{}{"token": signedToken}, nil)
}

func UpdateCstPass(w http.ResponseWriter, r *http.Request) {
	type ChangePass struct {
		OldPass		string		`json:"oldPass"`
		NewPass		string		`json:"newPass"`
	}
	changePass := ChangePass{}
	if err := json.NewDecoder(r.Body).Decode(&changePass); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}

	if changePass.OldPass == "" || changePass.NewPass == "" {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidRequestError))
		return
	}

	dataUpdate := dictionary.Customer{}
	dataUpdate.Pass = changePass.NewPass

	UpdateCustomer(w, r, dataUpdate, changePass.OldPass)
}

func UpdateCstExceptPass(w http.ResponseWriter, r *http.Request) {
	dataUpdate := dictionary.Customer{}
	if err := json.NewDecoder(r.Body).Decode(&dataUpdate); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}

	if dataUpdate.Email == "" || dataUpdate.Fullname == "" || dataUpdate.Pass == "" {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidRequestError))
		return
	}

	oldPass := dataUpdate.Pass 
	dataUpdate.Pass = ""

	UpdateCustomer(w, r, dataUpdate, oldPass)
}

func UpdateCustomer(
	w http.ResponseWriter, r *http.Request, dataUpdate dictionary.Customer, oldPass string,
) {
	custId := utils.GetIdCustomerInfoCtx(w, r)
	if custId == 0 {
		utils.JsonResp(w, 401, nil, errors.New(dictionary.UnauthorizedError))
		return
	}

	dataUpdate.Id = custId
	if oldPass == "" {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidRequestError))
		return
	}

	// Check database
	resDB, err := service.Login(dataUpdate.Id, "")
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.JsonResp(w, 400, nil, errors.New("invalid email or password"))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}
	
	if err := resDB.CheckPassword(oldPass); err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			utils.JsonResp(w, 400, nil, errors.New("invalid username or password"))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}

	if (dataUpdate.Pass == "") {
		dataUpdate.Pass = resDB.Pass
	} else {
		if err := dataUpdate.HashPassword(dataUpdate.Pass); err != nil {
			utils.JsonResp(w, 500, nil, err)
			return
		}
		dataUpdate.Email = resDB.Email
		dataUpdate.Fullname = resDB.Fullname
	}

	if dataUpdate.Email == "" || dataUpdate.Fullname == "" || dataUpdate.Pass == "" {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidRequestError))
		return
	}

	if err := service.UpdateCustomer(dataUpdate); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}

	// Success
	utils.JsonResp(w, 200, map[string]string{"message": "success"}, nil)
}

func InsertUpdateCstDetails(w http.ResponseWriter, r *http.Request) {
	custId := utils.GetIdCustomerInfoCtx(w, r)
	if custId == 0 {
		utils.JsonResp(w, 401, nil, errors.New(dictionary.UnauthorizedError))
		return
	}

	cstDetails := dictionary.CustomerDetail{}
	if err := json.NewDecoder(r.Body).Decode(&cstDetails); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	cstDetails.Id = custId

	if (cstDetails.NIK == "" || cstDetails.LegalName == "" || 
		cstDetails.PlaceBirth == "" || cstDetails.DateBirth == "" || 
		cstDetails.KtpImg == "" || cstDetails.SelfieImg == "" ){
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidRequestError))
		return
	}

	err := service.InsertUpdateCstDetails(cstDetails)
	if err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "customer_details_nik_key"` {
			utils.JsonResp(w, 400, nil, errors.New("NIK already registered"))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}

	countMakeTenorLimit := 0
	for i := 1; i <= 4; i++ {
		if ok := MakeTenorLimit(w, r, cstDetails.Salary, int64(i)); ok {
			countMakeTenorLimit++
		}
	}

	// Success
	resMap := map[string]interface{}{"message": "success", "countMakeTenorLimit": countMakeTenorLimit}
	utils.JsonResp(w, 200, resMap, nil)
}
