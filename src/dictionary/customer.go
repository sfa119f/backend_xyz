package dictionary

import (
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	Id				int 		`json:"id"`
	Fullname	string	`json:"fullname"`
	Email			string	`json:"email"`
	Pass			string	`json:"pass"`
	Salary		int			`json:"salary"`
}

func (customer *Customer) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	customer.Pass = string(bytes)
	return nil
}

func (customer *Customer) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(customer.Pass), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
