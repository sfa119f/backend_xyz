package utils

import (
	"encoding/json"
	"net/http"
	"fmt"
	"context"
	"strings"
	"errors"
	"os"

	"github.com/sfa119f/backend_xyz/src/dictionary"
	
	jwt "github.com/golang-jwt/jwt/v4"
	_		"github.com/joho/godotenv/autoload"
)

func MiddlewareJWTAuthorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/" || path == "/src/login" || (path == "/src/customer" && r.Method == "POST") {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			JsonResp(w, 400, nil, errors.New("invalid token"))
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("Signing method invalid")
			} else if method != jwt.SigningMethodHS256 {
				return nil, errors.New("Signing method invalid")
			}
			strKey := os.Getenv("XYZ_SECRET_KEY")
			key := []byte(strKey)
			return key, nil
		})
		if err != nil {
			JsonResp(w, 400, nil, err)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			JsonResp(w, 400, nil, err)
			return
		}

		ctx := context.WithValue(context.Background(), "customerInfo", claims)
		r = r.WithContext(ctx)
		
		next.ServeHTTP(w, r)
	})
}

func JsonResp(w http.ResponseWriter, code int, data interface{}, err error) {
	w.WriteHeader(code)

	if data != nil {
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: data, 
			Error: dictionary.NoError,
		})
	} else if code == 400 {
		fmt.Println("error message:", err)
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: nil, Error: err.Error(),
		})
	} else if err != nil {
		fmt.Println("error message:", err)
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: nil, Error: dictionary.UndisclosedError,
		})
	} else {
		fmt.Println("error message:", dictionary.UndisclosedError)
		json.NewEncoder(w).Encode(dictionary.APIResponse{
			Data: nil, Error: dictionary.UndisclosedError,
		})
	}
}
