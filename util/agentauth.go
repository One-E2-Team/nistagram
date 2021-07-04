package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strconv"
	"time"
)

const AgentExpiresIn = 86400000


func CreateAgentToken(id uint) (string, error) {
	if TokenSecret == "" {
		initPublicToken()
	}
	claims := TokenClaims{
		LoggedUserId: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + AgentExpiresIn,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(TokenSecret))
}

func parseAgentTokenClaim(r *http.Request) (*TokenClaims, bool) {
	if TokenSecret == "" {
		initPublicToken()
	}
	tokenString, err := getToken(r.Header)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenSecret), nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	if token.Valid {
		return token.Claims.(*TokenClaims), true
	} else {
		return nil, false
	}
}

func ValidateAgentToken(r *http.Request) (bool, uint) {
	if claims, ok := parseAgentTokenClaim(r); ok {
		if agentHasAPIAccessPrivilege(claims.LoggedUserId) {
			return true, claims.LoggedUserId
		}
	}
	return false, 0
}

func agentHasAPIAccessPrivilege(id uint) bool {
	privileges, ok := GetUserPrivileges(id)
	if !ok {
		return false
	}
	for _, val := range privileges {
		if val == "AGENT_API_ACCESS" {
			return true
		}
	}
	return false
}

func AgentAuth(handler func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {

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
		if check, id := ValidateAgentToken(request); check {
			writer.Header().Set("initiator", "API_" + strconv.Itoa(int(id)))
			finalHandler(true)(writer, request)
		} else {
			writer.Header().Set("initiator", "API_UNAUTHORIZED")
			finalHandler(false)(writer, request)
		}
	}
}