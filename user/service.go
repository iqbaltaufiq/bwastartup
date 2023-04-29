package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input UserLoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(id int, location string) (User, error)
	GetUserById(userId int) (User, error)
	GetAllUsers() ([]User, error)
	UpdateUser(input FormUpdateUserInput) (User, error)
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

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {
	user, errFind := s.repository.FindByEmail(input.Email)
	if user.Id == 0 {
		return true, nil
	}
	if errFind != nil {
		return false, errFind
	}
	return false, nil
}

func (s *service) SaveAvatar(id int, location string) (User, error) {
	user, errFind := s.repository.FindByID(id)
	if errFind != nil {
		return user, errFind
	}

	user.AvatarFileName = location

	updatedUser, errUpdate := s.repository.Update(user)
	if errUpdate != nil {
		return updatedUser, errUpdate
	}

	return user, nil
}

func (s *service) GetUserById(userId int) (User, error) {
	user, errFind := s.repository.FindByID(userId)
	if errFind != nil {
		return user, errFind
	}
	if user.Id == 0 {
		return user, errors.New("no user found for that id")
	}

	return user, nil
}

func (s *service) GetAllUsers() ([]User, error) {
	users, err := s.repository.FindAll()
	if err != nil {
		return users, err
	}

	return users, nil
}

func (s *service) UpdateUser(input FormUpdateUserInput) (User, error) {
	user, err := s.repository.FindByID(input.Id)
	if err != nil {
		return user, err
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation

	updatedUser, err := s.repository.Update(user)
	if err != nil {
		return updatedUser, err
	}

	return updatedUser, nil
}
