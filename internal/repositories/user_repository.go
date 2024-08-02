package repositories

import (
	"chat_app/internal/models"
	"errors"

	"gorm.io/gorm"
)

type AuthRepository struct {
	DB *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *AuthRepository {
	return &AuthRepository{DB: db}
}

func (repo *AuthRepository) CreateUser(user *models.User) error {
	return repo.DB.Create(&user).Error
}

func (repo *AuthRepository) UpdateUser(user *models.User) error {

	return repo.DB.Save(&user).Error
}

func (repo *AuthRepository) PhoneExists(phone string) (*uint, bool, error) {
	var user models.User
	if err := repo.DB.Where("phone = ?", phone).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, false, nil
		}
		return nil, false, err
	}
	return &user.ID, true, nil
}

func (repo *AuthRepository) GetUserInfo(ID uint) (*models.User, error) {
	var user models.User
	err := repo.DB.First(&user, ID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *AuthRepository) GetAllUsers() ([]*models.User, error) {
	var users []*models.User
	err := repo.DB.Find(&users).Error
	return users, err
}
