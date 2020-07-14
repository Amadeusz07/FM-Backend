package auth

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var mySigningKey = []byte("secretKey")

func IsAuthorized(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if header := r.Header["Authorization"]; header != nil {
			splitValue := strings.Fields(header[0])
			if splitValue[0] == "Bearer" {
				token, err := jwt.Parse(splitValue[1], func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("Error on extracting auth token")
					}
					return mySigningKey, nil
				})

				if err != nil {
					fmt.Println(err)
				}

				if token.Valid {
					claims := token.Claims.(jwt.MapClaims)
					if !claims.VerifyNotBefore(time.Now().Unix(), true) {
						r.Header.Add("userId", claims["userId"].(string))
						next.ServeHTTP(w, r)
					}
				} else {
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

		} else {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	})
}

func GenerateJWT(userId primitive.ObjectID, username string, expDate time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["userId"] = userId.Hex()
	claims["user"] = username
	claims["exp"] = expDate.Unix()

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("JWT token generation went wrong")
		return "", err
	}

	return tokenString, err
}

func GenerateJWTWithProject(userId primitive.ObjectID, username string, expDate time.Time, projectId primitive.ObjectID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["userId"] = userId.Hex()
	claims["user"] = username
	claims["exp"] = expDate.Unix()
	claims["projectId"] = projectId

	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Errorf("JWT token generation went wrong")
		return "", err
	}

	return tokenString, err
}
