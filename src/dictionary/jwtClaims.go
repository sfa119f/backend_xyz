package dictionary

import (
	jwt "github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	jwt.StandardClaims
	Id				int 		`json:"id"`
	Fullname	string	`json:"fullname"`
	Email			string	`json:"email"`
}