package service

import (
	"github.com/sfa119f/backend_xyz/src/database"
	"github.com/sfa119f/backend_xyz/src/dictionary"
)

func InsertCustomer(customer dictionary.Customer) (int64, error) {
	db := database.GetDB()
	query := `insert into customers (fullname, email, pass) values ($1, $2, $3) returning id`

	var id int64
	if err := db.QueryRow(query, customer.Fullname, customer.Email, customer.Pass).Scan(&id);
	err != nil {
		return 0, err
	}
	return id, nil
}

func Login(id int64, email string) (dictionary.Customer, error) {
	db := database.GetDB()
	
	query := `
		select * from customers 
		where case 
			when $1 = 0 then email = $2
			else id = $1
		end
	`

	res := dictionary.Customer{}
	if err := 
		db.QueryRow(query, id, email).Scan(&res.Id, &res.Fullname, &res.Email, &res.Pass);
		err != nil {
			return res, err
		}
	return res, nil
}

func UpdateCustomer(customer dictionary.Customer) error {
	db := database.GetDB()

	query := `
		update customers set fullname = $2, email = $3, pass = $4
		where id = $1 returning id
	`

	var id int64
	if err := 
		db.QueryRow(
			query, customer.Id, customer.Fullname, customer.Email, customer.Pass,
		).Scan(&id);
		err != nil {
			return err
		}
	return nil
}
