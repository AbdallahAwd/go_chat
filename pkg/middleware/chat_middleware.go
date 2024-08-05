package middlewares

// func ChatMiddleware(jwtSecret string) func(http.Handler) http.Handler {
// 	return func(next http.Handler) http.Handler {
// 		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			authToken := r.Header.Get("token")
// 			if authToken == "" {
// 				utils.ErrorJSON(w, "Unauthorized")
// 				return
// 			}
// 			authToken, ok := strings.CutPrefix(authToken, "Bearer ")
// 			if !ok {
// 				utils.ErrorJSON(w, "Invalid Token Format")
// 				return
// 			}
// 			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
// 				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 				}
// 				return []byte(jwtSecret), nil
// 			})

// 			if err != nil || !token.Valid {
// 				utils.ErrorJSON(w, "Invalid Token")
// 				return
// 			}
// 			if claims, ok := token.Claims.(jwt.MapClaims); ok {
// 				log.Printf("Claims: %+v", claims)
// 				ID, ok := claims["ID"]

// 				if !ok {
// 					utils.ErrorJSON(w, "Invalid claims data", http.StatusInternalServerError)
// 					return
// 				}

// 				ctx := context.WithValue(r.Context(), utils.ID, ID)

// 				next.ServeHTTP(w, r.WithContext(ctx))
// 			} else {
// 				utils.ErrorJSON(w, "Invalid token claims", http.StatusUnauthorized)
// 			}
// 			next.ServeHTTP(w, r)
// 		})
// 	}
// }
