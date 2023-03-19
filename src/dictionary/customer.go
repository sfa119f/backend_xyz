package dictionary

import (
	"golang.org/x/crypto/bcrypt"
)

type Customer struct {
	Id				int64		`json:"id"`
	Fullname	string	`json:"fullname"`
	Email			string	`json:"email"`
	Pass			string	`json:"pass"`
}

type CustomerDetail struct {
	Id					int64		`json:"id"`
	NIK					string	`json:"nik"`
	LegalName		string	`json:"legalname"`
	PlaceBirth	string 	`json:"place_of_birth"`
	DateBirth		string	`json:"date_of_birth"`
	Salary			int64		`json:"salary"`
	KtpImg			string  `json:"ktp_img"`
	SelfieImg		string	`json:"selfie_img"`
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
