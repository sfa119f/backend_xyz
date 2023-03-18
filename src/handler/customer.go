package handler

import (
	"encoding/json"
	"net/http"
	"fmt"

	"github.com/sfa119f/backend_xyz/src/dictionary"
	"github.com/sfa119f/backend_xyz/src/service"
)

func InsertCustomer(w http.ResponseWriter, r *http.Request) {
	customer := dictionary.Customer{}
	json.NewDecoder(r.Body).Decode(&customer)

	// fmt.Println(customer.Salary)
	if customer.Fullname == "" || customer.Pass == "" || customer.Email == "" {
		fmt.Println("err insert product:", "Body must have fullname, email, and pass")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: nil, Error: "body must have fullname, email, and pass",
		})
		return
	}

	if err := customer.HashPassword(customer.Pass); err != nil {
		w.WriteHeader(500)
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: nil, Error: dictionary.UndisclosedError,
		})
	}

	if err := service.InsertCustomer(customer); err != nil {
		fmt.Println("err insert product:", err)
		if (err.Error() == `pq: duplicate key value violates unique constraint "customers_email_key"`) {
			w.WriteHeader(400)
			json.NewEncoder(w).Encode(dictionary.APIResponse{
				Data: nil, Error: "email already registered",
			})
		} else {
			w.WriteHeader(500)
			json.NewEncoder(w).Encode(dictionary.APIResponse{
				Data: nil, Error: dictionary.UndisclosedError,
			})
		}
		return
	}

	w.WriteHeader(201)
	json.NewEncoder(w).Encode(dictionary.APIResponse{
		Data: map[string]string{"message": "Success"}, 
		Error: dictionary.NoError,
	})
}