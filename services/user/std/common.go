package stduserservice

import "golang.org/x/crypto/bcrypt"

func hashed(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), 6)
	if err != nil {
		return "", err
	}

	return string(hashed), nil
}
