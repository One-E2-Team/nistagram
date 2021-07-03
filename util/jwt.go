package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const ExpiresIn = 86400

var TokenSecret = ""

type TokenClaims struct {
	LoggedUserId uint `json:"loggedUserId"`
	jwt.StandardClaims
}

func initPublicToken() {
	env := os.Getenv("PUBLIC_JWT_TOKEN_SECRET")
	if env == "" {
		TokenSecret = "token_secret"
	} else {
		TokenSecret = env
	}
}

func CreateToken(userId uint, issuer string) (string, error) {
	if TokenSecret == "" {
		initPublicToken()
	}
	claims := TokenClaims{LoggedUserId: userId, StandardClaims: jwt.StandardClaims{
		ExpiresAt: time.Now().Unix() + ExpiresIn,
		IssuedAt:  time.Now().Unix(),
		Issuer:    issuer,
	}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(TokenSecret))
}

func getToken(header http.Header) (string, error) {
	reqToken := header.Get("Authorization")
	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("NO_TOKEN")
	}
	return strings.TrimSpace(splitToken[1]), nil
}

func GetLoggedUserIDFromToken(r *http.Request) uint {
	if TokenSecret == "" {
		initPublicToken()
	}
	tokenString, err := getToken(r.Header)
	if err != nil {
		var err1 error
		tokenString, err1 = getTokenFromParams(r.URL.String())
		if err1 != nil {
			fmt.Println(err, err1)
			return 0
		}
	}
	return GetLoggedUserIDFromPureToken(tokenString)
}

func getTokenFromParams(s string) (string, error) {
	err := fmt.Errorf("generic error, no token in URL")
	paramsPart := strings.Split(s, "?")
	if len(paramsPart) < 2 {
		return "", err
	}
	tokenTilEnd := strings.Split(paramsPart[1],"token=")
	if len(paramsPart) < 2 {
		return "", err
	}
	tokenEscaped := strings.Split(tokenTilEnd[1], "&")
	token, err := url.QueryUnescape(tokenEscaped[0])
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	return token, nil
}


func GetLoggedUserIDFromPureToken(tok string) uint {
	if TokenSecret == "" {
		initPublicToken()
	}
	token, err := jwt.ParseWithClaims(tok, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenSecret), nil
	})
	if err != nil {
		fmt.Println(err)
		return 0
	}
	if claims, ok := token.Claims.(*TokenClaims); ok && token.Valid {
		return claims.LoggedUserId
	} else {
		fmt.Println(err)
		return 0
	}
}