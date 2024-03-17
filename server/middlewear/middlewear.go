package middlewear

import (
	"net/http"

	st "VK_app/server/structures"

	"log"

	"github.com/golang-jwt/jwt"
)

// CheckToken is a function that takes an http.Handler and returns an http.Handler.
//
// It checks the validity of the token in the request header and calls the next handler if the token is valid.
func CheckToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return st.Secret, nil
		})
		if err != nil || !token.Valid {
			log.Println(err)
			http.Error(w, "Unauthorized user", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// CheckTokenAdmin is a function that checks the token for admin access.
//
// It takes in a http.Handler as a parameter and returns a http.Handler.
func CheckTokenAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return st.Secret, nil
		})
		if err != nil || !token.Valid {
			log.Println(err)
			http.Error(w, "Unauthorized admin", http.StatusUnauthorized)
			return
		}
		if role, ok := token.Claims.(jwt.MapClaims)["role"].(int); ok && role != 1  {
			log.Println(err)
			http.Error(w, "Unauthorized admin role", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
