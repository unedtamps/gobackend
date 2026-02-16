package utils

import "golang.org/x/crypto/bcrypt"

func GenereateHash(password string) (string, error) {
	res, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(res), nil
}

func CompareHash(plain, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
}

// HashPassword hashes a password using bcrypt
func HashPassword(password string) (string, error) {
	return GenereateHash(password)
}

// CheckPassword verifies a password against a hash
func CheckPassword(password, hash string) bool {
	return CompareHash(password, hash) == nil
}
