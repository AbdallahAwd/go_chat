package middlewares

import (
	"chat_app/pkg/utils"
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
)

func PhoneMiddleware(jwtSecret string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			otpToken := r.Header.Get("Authorization")

			if otpToken == "" {
				utils.ErrorJSON(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}
			otpToken, ok := strings.CutPrefix(otpToken, "Bearer ")
			if !ok {
				utils.ErrorJSON(w, "Invalid Token Format", http.StatusUnauthorized)
				return
			}

			log.Printf("Received token: %s", otpToken)

			token, err := jwt.Parse(otpToken, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(jwtSecret), nil
			})
			if err != nil {
				log.Printf("Error parsing token: %v", err)
				utils.ErrorJSON(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				log.Printf("Claims: %+v", claims)
				phone, ok1 := claims["phone"].(string)
				otp, ok2 := claims["otp"].(string)
				if !ok1 || !ok2 {
					utils.ErrorJSON(w, "Invalid claims data", http.StatusInternalServerError)
					return
				}
				values := map[string]string{"phone": phone, "otp": otp}
				ctx := context.WithValue(r.Context(), utils.PhoneOTP, values)
				next.ServeHTTP(w, r.WithContext(ctx))
			} else {
				utils.ErrorJSON(w, "Invalid token claims", http.StatusUnauthorized)
			}
		})
	}
}
