package models

import jwt "github.com/dgrijalva/jwt-go"

//Token struct declaration
type Token struct {
	Email string
	*jwt.StandardClaims
}
