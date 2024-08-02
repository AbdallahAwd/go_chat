package services

import (
	"chat_app/config"
	"chat_app/internal/models"
	"chat_app/internal/repositories"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type AuthService struct {
	repo   *repositories.AuthRepository
	config *config.Config
}

func NewAuthService(r *repositories.AuthRepository, conf *config.Config) *AuthService {
	return &AuthService{repo: r, config: conf}
}

func (s *AuthService) ValidatePhone(code, phone string) (string, error) {
	otp := generateRandNumbers()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"otp":   otp,
		"phone": phone,
		"exp":   time.Now().Add(time.Minute * 5).Unix(),
	})
	strToken, err := token.SignedString([]byte(s.config.JwtSecret))
	if err != nil {
		return "", err
	}
	// TODO Sending SMS Message
	return strToken, nil
}

func (s *AuthService) Verify(bodyOtp, bodyPhone, claimOtp, claimPhone string) error {
	if bodyOtp != claimOtp {
		return fmt.Errorf("invalid OTP")
	}

	if bodyOtp == claimOtp {
		if bodyPhone != claimPhone {
			return fmt.Errorf("invalid OTP not phone")
		}
		return nil
	}
	return nil
}

func (s *AuthService) CreateOrSaveUser(imageFile *multipart.FileHeader, name, phone, claimPhone, notificationToken, code string) (string, error) {
	if phone != claimPhone {
		return "", fmt.Errorf("Unauthorized")
	}
	image, err := s.SaveImage(imageFile)
	if err != nil {
		return "", err
	}
	user := &models.User{
		Name:              name,
		CountryCode:       code,
		Image:             *image,
		NotificationToken: notificationToken,
		Phone:             phone,
	}
	err = s.repo.CreateUser(user)

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {

			ID, isExist, err := s.repo.PhoneExists(phone)
			if err != nil || ID == nil || !isExist {
				return "", fmt.Errorf("phone exist error %v", err)
			}
			user.ID = *ID
			print("User ID", user.ID)
			err = s.repo.UpdateUser(user)
			if err != nil {
				return "", err
			}
			token, err := s.generateToken(user.ID)
			if err != nil {
				return "", nil
			}
			return token, nil

		}
		return "", err
	}
	token, err := s.generateToken(user.ID)
	if err != nil {
		return "", nil
	}
	return token, nil
}
func (s *AuthService) GetUserInfo(ID uint) (*models.User, error) {
	return s.repo.GetUserInfo(ID)
}
func (s *AuthService) GetAllUsers() ([]*models.User, error) {
	return s.repo.GetAllUsers()
}

func (s *AuthService) generateToken(ID uint) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"ID": ID,
	})
	strToken, err := token.SignedString([]byte(s.config.JwtSecret))
	if err != nil {
		return "", err
	}
	return strToken, nil
}

func generateRandNumbers() string {
	var num [6]int
	for i := 0; i < 6; i++ {
		num[i] = rand.Intn(10)
	}

	return fmt.Sprintf("%d%d%d%d%d%d", num[0], num[1], num[2], num[3], num[4], num[5])
}

func (as *AuthService) SaveImage(image *multipart.FileHeader) (*string, error) {
	if image == nil {
		return nil, nil
	}

	// Create the upload directory if it doesn't exist
	if err := as.createUploadDir(); err != nil {
		return nil, err
	}

	// Construct the full path to save the image
	filename := filepath.Join(as.config.UploadPath, image.Filename)
	normalizedFilename := as.normalizePath(filename)

	// Save the image file
	if err := as.saveFile(image, filename); err != nil {
		return nil, err
	}

	return &normalizedFilename, nil
}

// createUploadDir ensures that the upload directory exists
func (as *AuthService) createUploadDir() error {
	return os.MkdirAll(as.config.UploadPath, os.ModePerm)
}

// saveFile saves the uploaded image to the destination path
func (as *AuthService) saveFile(image *multipart.FileHeader, dstPath string) error {
	src, err := image.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	dst, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	_, err = io.Copy(dst, src)
	return err
}

// normalizePath converts the path to use forward slashes
func (as *AuthService) normalizePath(path string) string {
	return strings.ReplaceAll(path, "\\", "/")
}
