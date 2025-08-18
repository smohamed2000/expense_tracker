package utils

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string ) (string, error) {
	hashpass, err := bcrypt.GenerateFromPassword([]byte(password),bcrypt.DefaultCost)
	if err != nil {
		return "",err
	}
	// we need to convert byte to a string becasue it is in bye format
	return string(hashpass), err
}
//compare password

func ComparePassword(hashPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword),[]byte(password))
	if err !=nil {
		return err

	}
	return nil
}