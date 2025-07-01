package pkg

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	hashCost = 10
)

type bcryptHashService struct {
	hashCost int
}

type HashService interface {
	GenerateHash(stringToHash string) ([]byte, error)
	CompareStringAndHash(stringToCompare string, hashedString string) bool
}

func NewBcryptHashService() HashService {
	return &bcryptHashService{hashCost: hashCost}
}

func (p *bcryptHashService) GenerateHash(stringToHash string) ([]byte, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(stringToHash), p.hashCost)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (p *bcryptHashService) CompareStringAndHash(stringToCompare string, hashedString string) bool {
	passwordMatched := bcrypt.CompareHashAndPassword([]byte(hashedString), []byte(stringToCompare))
	if passwordMatched != nil {
		return false
	}

	return true
}
