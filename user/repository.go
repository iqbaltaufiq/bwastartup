package user

import "gorm.io/gorm"

type Repository interface {
	Save(user User) (User, error)
	FindByEmail(email string) (User, error)
	FindByID(ID int) (User, error)
	Update(user User) (User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Save(user User) (User, error) {
	err := r.db.Create(&user).Error
	if err != nil {
		return user, err
	}

	return user, nil
}

func (r *repository) FindByEmail(email string) (User, error) {
	var user User
	errFind := r.db.Where("email = ?", email).First(&user).Error
	if errFind != nil {
		return user, errFind
	}

	return user, nil
}

func (r *repository) FindByID(ID int) (User, error) {
	var user User
	errFind := r.db.First(&user, ID).Error
	if errFind != nil {
		return user, errFind
	}

	return user, nil
}

func (r *repository) Update(user User) (User, error) {
	if errSave := r.db.Save(&user).Error; errSave != nil {
		return user, errSave
	}

	return user, nil
}
