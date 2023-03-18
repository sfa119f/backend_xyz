package service

import (
	"github.com/sfa119f/backend_xyz/src/database"
	"github.com/sfa119f/backend_xyz/src/dictionary"
)

func InsertCustomer(customer dictionary.Customer) error {
	db := database.GetDB()
	query := `insert into customers (fullname, email, pass, salary) values ($1, $2, $3, $4)`

	_, err := db.Exec(query, customer.Fullname, customer.Email, customer.Pass, customer.Salary)

	return err
}

func Login(email string) (dictionary.Customer, error) {
	db := database.GetDB()
	query := `select id, fullname, email, pass from customers where email = $1`

	res := dictionary.Customer{}
	if err := 
		db.QueryRow(query, email).Scan(&res.Id, &res.Fullname, &res.Email, &res.Pass);
		err != nil {
			return res, err
		}
	return res, nil
}
