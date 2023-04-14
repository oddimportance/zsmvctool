package lib

import (
	//	"log"
	"unicode"
)

type HandlePassword struct {
}

func (h *HandlePassword) ConfirmPassword(password string, comparePassword string) bool {
	return password == comparePassword
}

func (h *HandlePassword) ValidatePassword(password string) bool {
	var (
		hasMinLen         = false
		hasUpper          = false
		hasLower          = false
		hasNumber         = false
		hasSpecial        = false
		passwordMinLength = 8
	)
	if len(password) >= passwordMinLength {
		//		log.Println("Length")
		hasMinLen = true
	}

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			//			log.Println("Is uppder")
			hasUpper = true
		case unicode.IsLower(char):
			//			log.Println("has lower")
			hasLower = true
		case unicode.IsNumber(char):
			//			log.Println("has number")
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			//			log.Println("has spechial")
			hasSpecial = true
		}
	}

	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}
