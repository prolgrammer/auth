package pkg

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

type tokenService struct {
	signCert []byte
}

var (
	ErrExpired = errors.New("token is expired")
	ErrInvalid = errors.New("token is invalid")
)

type TokenService interface {
	CreateAccessToken(claims map[string]interface{}) (string, error)
	ParseAccessToken(token string) (map[string]interface{}, error)

	CreateRefreshToken() (string, error)
}

func NewTokenService(signCert []byte) TokenService {
	return &tokenService{signCert: signCert}
}

func (j *tokenService) CreateAccessToken(claim map[string]interface{}) (string, error) {
	var mapClaims jwt.MapClaims
	mapClaims = claim
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, mapClaims)
	tokenString, err := token.SignedString(j.signCert)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (j *tokenService) ParseAccessToken(token string) (map[string]interface{}, error) {
	claims := jwt.MapClaims{}
	parsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return j.signCert, nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrExpired
		}
		return nil, errors.Join(ErrInvalid, err)
	}

	claims = parsed.Claims.(jwt.MapClaims)

	return claims, err
}

func (j *tokenService) CreateRefreshToken() (string, error) {
	rndBytes := make([]byte, 32)
	if _, err := rand.Read(rndBytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(rndBytes), nil
}
