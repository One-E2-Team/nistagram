package util

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"time"
)

const AgentExpiresIn = 86400000

var AgentTokenSecret = ""

type AgentTokenClaims struct {
	AgentID uint `json:"agentID"`
	jwt.StandardClaims
}

func initAgentToken() {
	env := os.Getenv("AGENT_JWT_TOKEN_SECRET")
	if env == "" {
		AgentTokenSecret = "token_secret"
	} else {
		AgentTokenSecret = env
	}
}

func CreateAgentToken(id uint) (string, error) {
	if AgentTokenSecret == "" {
		initAgentToken()
	}
	claims := AgentTokenClaims{
		AgentID: id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Unix() + AgentExpiresIn,
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(AgentTokenSecret))
}

func parseAgentTokenClaim(r *http.Request) (*AgentTokenClaims, bool) {
	if AgentTokenSecret == "" {
		initAgentToken()
	}
	tokenString, err := getToken(r.Header)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	token, err := jwt.ParseWithClaims(tokenString, &AgentTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(AgentTokenSecret), nil
	})
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	if token.Valid {
		return token.Claims.(*AgentTokenClaims), true
	} else {
		return nil, false
	}
}

func ValidateAgentToken(r *http.Request) bool {
	if claims, ok := parseAgentTokenClaim(r); ok {
		if agentHasAPIAccessPrivilege(claims.AgentID) {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
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
		if ValidateAgentToken(request) {
			finalHandler(true)(writer, request)
		} else {
			finalHandler(false)(writer, request)
		}
	}
}