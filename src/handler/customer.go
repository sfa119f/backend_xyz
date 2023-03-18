package handler

import (
	"encoding/json"
	"net/http"
	"errors"
	"os"
	"time"

	"github.com/sfa119f/backend_xyz/src/dictionary"
	"github.com/sfa119f/backend_xyz/src/service"
	"github.com/sfa119f/backend_xyz/src/utils"
	
	jwt "github.com/golang-jwt/jwt/v4"
	_		"github.com/joho/godotenv/autoload"
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
	customer := dictionary.Customer{}
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		utils.JsonResp(w, 500, nil, err)
		return
	}
	
	// Check database
	resDB, err := service.Login(customer.Email)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			utils.JsonResp(w, 400, nil, errors.New("invalid username or password"))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}
	
	if err := resDB.CheckPassword(customer.Pass); err != nil {
		if err.Error() == "crypto/bcrypt: hashedPassword is not the hash of the given password" {
			utils.JsonResp(w, 400, nil, errors.New("invalid username or password"))
		} else {
			utils.JsonResp(w, 500, nil, err)
		}
		return
	}

	appName := os.Getenv("APP_NAME")
	claims := dictionary.JwtClaims{
    StandardClaims: jwt.StandardClaims{
			Issuer: appName,
			ExpiresAt: time.Now().Add(time.Duration(10) * time.Minute).Unix(),
    },
    Id: resDB.Id,
    Fullname: resDB.Fullname,
    Email: resDB.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	strKey := os.Getenv("XYZ_SECRET_KEY")
	key := []byte(strKey)

	signedToken, err := token.SignedString(key)
	if err != nil {
		utils.JsonResp(w, 400, nil, err)
		return
	}

	// Success
	utils.JsonResp(w, 200, map[string]interface{}{"token": signedToken}, nil)
}
