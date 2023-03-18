package handler

import (
	"encoding/json"
	"net/http"
	"errors"
	"strconv"

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

	if err := service.InsertCustomer(customer); err != nil {
		if err.Error() == `pq: duplicate key value violates unique constraint "customers_email_key"` {
			utils.JsonResp(w, 400, nil, errors.New("email already registered"))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}

	// Success
	utils.JsonResp(w, 200, map[string]string{"message": "success"}, nil)
}

func Login(w http.ResponseWriter, r *http.Request) {
	token := ""
	if err := json.NewDecoder(r.Body).Decode(&token); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}

	if token == "" {
		utils.JsonResp(w, 400, nil, errors.New(dictionary.InvalidRequestError))
		return
	}

	reqByte, err := utils.DecryptAES(token)
	if err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}

	customer := dictionary.Customer{}
	if err := json.Unmarshal(reqByte, &customer); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	
	// Check database
	resDB, err := service.Login(customer.Email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.JsonResp(w, 400, nil, errors.New(dictionary.NotFoundError))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}
	
	if err := resDB.CheckPassword(customer.Pass); err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			utils.JsonResp(w, 400, nil, errors.New(dictionary.NotFoundError))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}

	mapTokenRes := map[string]string{"id": strconv.Itoa(resDB.Id)}
	strTokenRes, err := json.Marshal(mapTokenRes)
	if err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	
	tokenRes, err := utils.EncryptAES(strTokenRes)
	if err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	
	mapResult := map[string]string{
		"token": tokenRes,
		"email": resDB.Email,
		"fullname": resDB.Fullname,
	}

	// Success
	utils.JsonResp(w, 200, mapResult, nil)
}