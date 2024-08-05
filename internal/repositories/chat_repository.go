package repositories

import (
	"chat_app/internal/models"

	"gorm.io/gorm"
)

type ChatRepo struct {
	db *gorm.DB
}

func NewChatRepo(DB *gorm.DB) *ChatRepo {
	return &ChatRepo{
		db: DB,
	}
}

func (c *ChatRepo) Create(message *models.Message) error {
	return c.db.Create(&message).Error
}

func (r *ChatRepo) GetMessagedUsers(userID uint) ([]models.User, error) {
	var users []models.User

	recipientSubQuery := r.db.Model(&models.Message{}).
		Where("user_id = ?", userID).
		Select("distinct recipient_id")

	senderSubQuery := r.db.Model(&models.Message{}).
		Where("recipient_id = ?", userID).
		Select("distinct user_id")

	if err := r.db.Where("id IN (?) OR id IN (?)", recipientSubQuery, senderSubQuery).
		Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *ChatRepo) GetMessagesBetweenUsers(userID1, userID2 uint, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.Where("(user_id = ? AND recipient_id = ?) OR (user_id = ? AND recipient_id = ?)",
		userID1, userID2, userID2, userID1).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *ChatRepo) Update(message *models.Message) error {
	return r.db.Save(message).Error
}

func (r *ChatRepo) Delete(id uint) error {
	return r.db.Delete(&models.Message{}, id).Error
}
