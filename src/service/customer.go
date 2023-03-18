package service

import (
	"github.com/sfa119f/backend_xyz/src/database"
	"github.com/sfa119f/backend_xyz/src/dictionary"
)

func InsertCustomer(customer dictionary.Customer) error {
	db := database.GetDB()
	query := `insert into customers (fullname, email, pass) values ($1, $2, $3)`

	_, err := db.Exec(query, customer.Fullname, customer.Email, customer.Pass)
	return err
}

func Login(email string) (dictionary.Customer, error) {
	db := database.GetDB()
	query := `select * from customers where email = $1`

	res := dictionary.Customer{}
	if err := 
		db.QueryRow(query, email).Scan(&res.Id, &res.Fullname, &res.Email, &res.Pass);
		err != nil {
			return res, err
		}
	return res, nil
}
