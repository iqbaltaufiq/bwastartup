package auth

import "github.com/golang-jwt/jwt"

type Service interface {
	GenerateToken(userId int) (string, error)
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
