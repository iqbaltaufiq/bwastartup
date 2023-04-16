package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input UserLoginInput) (User, error)
}

type service struct {
	repository Repository
}

func NewService(repository Repository) *service {
	return &service{
		repository: repository,
	}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {
	user := User{}
	user.Name = input.Name
	user.Occupation = input.Occupation
	user.Email = input.Email

	hash, errHash := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if errHash != nil {
		return user, errHash
	}
	user.PasswordHash = string(hash)
	user.Role = "user"

	newUser, errSave := s.repository.Save(user)
	if errSave != nil {
		return user, errSave
	}

	return newUser, nil
}

func (s *service) Login(input UserLoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, errFind := s.repository.FindByEmail(email)
	if errFind != nil {
		return user, errFind
	}

	if user.Id == 0 {
		return user, errors.New("no user found for that email")
	}

	errCompare := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if errCompare != nil {
		return user, errCompare
	}

	return user, nil
}
