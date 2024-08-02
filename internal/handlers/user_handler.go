package handlers

import (
	"chat_app/internal/services"
	"chat_app/pkg/utils"
	"encoding/json"
	"log"
	"net/http"
)

type AuthHandler struct {
	service *services.AuthService
}

func NewAuthHandler(s *services.AuthService) *AuthHandler {
	return &AuthHandler{service: s}
}

func (h *AuthHandler) ValidatePhone(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Code  string `json:"code" `
		Phone string `json:"phone" `
	}
	// validate := validator.New()
	// err := validate.Struct(req)
	// if err != nil {
	// 	utils.ResponseHandler(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
	// 	return
	// }
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, "Decode Error"+err.Error(), http.StatusBadRequest)
		return
	}

	token, err := h.service.ValidatePhone(req.Code, req.Phone)
	if err != nil {
		utils.ErrorJSON(w, "Service Error "+err.Error(), http.StatusBadRequest)
		return
	}
	utils.ResponseHandler(w, map[string]string{"message": "Sent An SMS message", "otp_token": token})

}

func (h *AuthHandler) VerifyPhone(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Otp   string `json:"otp"`
		Phone string `json:"phone"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.ErrorJSON(w, "Decode Error"+err.Error())
		return
	}
	claims, ok := r.Context().Value(utils.PhoneOTP).(map[string]string)
	if !ok || claims == nil {
		utils.ErrorJSON(w, "Failed to get claims from context", http.StatusInternalServerError)
		return
	}

	log.Printf("Claims in handler: %+v", claims)

	otpClaim := claims["otp"]
	phoneClaim := claims["phone"]

	err := h.service.Verify(req.Otp, req.Phone, otpClaim, phoneClaim)

	if err != nil {
		utils.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send a success response
	utils.ResponseHandler(w, map[string]string{"message": "Phone verified successfully"})
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	phone := r.FormValue("phone")
	notificationToken := r.FormValue("notification_token")
	_, image, err := r.FormFile("image")
	if err != nil {
		utils.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	}
	claim, ok := r.Context().Value(utils.PhoneOTP).(map[string]string)
	if !ok {
		utils.ErrorJSON(w, "Cannot Get Claims", http.StatusInternalServerError)
		return
	}
	if token, err := h.service.CreateOrSaveUser(image, name, phone, claim["phone"], notificationToken, "+2"); err != nil {
		utils.ErrorJSON(w, err.Error(), http.StatusBadRequest)
		return
	} else {
		utils.ResponseHandler(w, map[string]string{"messsage": "success", "token": token})
	}

}

func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	id, ok := r.Context().Value(utils.ID).(float64)
	if !ok {
		utils.ErrorJSON(w, "Invalid claims data", http.StatusInternalServerError)
		return
	}

	user, err := h.service.GetUserInfo(uint(id))
	if err != nil {
		utils.ErrorJSON(w, "User Error"+err.Error(), http.StatusBadRequest)
		return
	}
	utils.ResponseHandler(w, user)
}

func (h *AuthHandler) GetAllUser(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		utils.ErrorJSON(w, "Users Error"+err.Error(), http.StatusInternalServerError)
		return
	}
	utils.ResponseHandler(w, users)
}
