package services

import (
	"chat_app/internal/models"
	"chat_app/internal/repositories"
)

type ChatService struct {
	repo *repositories.ChatRepo
}

func NewChatService(r *repositories.ChatRepo) *ChatService {
	return &ChatService{repo: r}
}

func (c *ChatService) SendMessage(content, audioUrl, imageUrl string, userId, senderId uint) (*models.Message, error) {
	message := &models.Message{
		Content:     content,
		AudioUrl:    audioUrl,
		ImageUrl:    imageUrl,
		UserID:      userId,
		RecipientID: senderId,
	}
	err := c.repo.Create(message)
	if err != nil {
		return nil, err
	}
	return message, nil
}

func (c *ChatService) GetConversation(userId uint) ([]models.User, error) {
	users, err := c.repo.GetMessagedUsers(userId)
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (c *ChatService) GetChatPartners(userID1 uint, userID2 uint, limit int, offset int) ([]models.Message, error) {
	messages, err := c.repo.GetMessagesBetweenUsers(userID1, userID2, limit, offset)
	if err != nil {
		return nil, err
	}
	return messages, nil
}
func (c *ChatService) EditMessage(message *models.Message) error {
	return c.repo.Update(message)
}
func (c *ChatService) DeleteMessage(id uint) error {
	return c.repo.Delete(id)
}
