package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v4"
)

type Service interface {
	GenerateToken(userId int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type service struct {
}

func NewService() *service {
	return &service{}
}

var JWT_KEY string = "JWT_TOKEN_KEY"

func (s *service) GenerateToken(userId int) (string, error) {
	claim := jwt.MapClaims{}
	claim["user_id"] = userId

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	signedToken, errSign := token.SignedString([]byte(JWT_KEY))
	if errSign != nil {
		return signedToken, errSign
	}

	return signedToken, nil
}

func (s *service) ValidateToken(encodedToken string) (*jwt.Token, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, errors.New("invalid token")
		}

		return []byte(JWT_KEY), nil
	}

	parsedToken, errParse := jwt.Parse(encodedToken, keyFunc)

	if errParse != nil {
		return parsedToken, errParse
	}

	return parsedToken, nil
}
