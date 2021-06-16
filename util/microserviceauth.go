package util

import (
	"bytes"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

var msJwt string

const MSExpiresIn = 86400000

var MSTokenSecret = os.Getenv("MICROSERVICE_JWT_TOKEN_SECRET")

type MSTokenClaims struct {
	Microservice string `json:"microservice"`
	jwt.StandardClaims
}

func SetupMSAuth(ms string) error {
	claims := MSTokenClaims{Microservice: ms, StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + MSExpiresIn,
		IssuedAt:  time.Now().Unix(),
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	var err error
	msJwt, err = token.SignedString([]byte(MSTokenSecret))
	return err
}

func ValidateMSToken(r *http.Request, ms []string) bool {
	tokenString, err := getToken(r.Header)
	if err != nil {
		fmt.Println(err)
		return false
	}
	token, err := jwt.ParseWithClaims(tokenString, &MSTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(MSTokenSecret), nil
	})
	if err != nil {
		fmt.Println(err)
		return false
	}
	if claims, ok := token.Claims.(*MSTokenClaims); ok && token.Valid {
		for _, value := range ms {
			if claims.Microservice == value {
				return true
			}
		}
		return false
	} else {
		fmt.Println(err)
		return false
	}
}

func CrossServiceRequest(method string, path string, data []byte, headers map[string]string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(method, path, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+msJwt)
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	return client.Do(req)
}

func MSAuth(handler func(http.ResponseWriter, *http.Request), microservices []string) func(http.ResponseWriter, *http.Request) {

	finalHandler := func(pass bool) func(http.ResponseWriter, *http.Request) {
		if pass {
			return handler
		} else {
			return func(writer http.ResponseWriter, request *http.Request) {
				writer.WriteHeader(http.StatusOK)
				writer.Header().Set("Content-Type", "application/json")
				_, _ = writer.Write([]byte("{\"status\":\"fail\", \"reason\":\"unauthorized\"}"))
			}
		}
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		if ValidateMSToken(request, microservices) {
			finalHandler(true)(writer, request)
		} else {
			finalHandler(false)(writer, request)
		}
	}
}
