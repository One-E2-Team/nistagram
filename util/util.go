package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/smtp"
	"nistagram/auth/model"
	"os"
	"strconv"
)

func SendMail(sendTo string, subject string, mailMessage string) {
	from := os.Getenv("MAIL_USERNAME")
	password := os.Getenv("MAIL_PASSWORD")
	to := []string{sendTo}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	msg := []byte("To: " + sendTo + "\r\n" + "Subject: " + subject + "\r\n" + "\r\n" + mailMessage + "\r\n")
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func Uint2String(input uint) string{
	return strconv.FormatUint(uint64(input), 10)
}

func String2Uint(input string) uint {
	u64, err := strconv.ParseUint(input, 10, 32)
	if err != nil {
		fmt.Println(err)
		return 0
	}
	return uint(u64)
}

func userHasPermission(userId uint) func(string) bool {
	return func(permission string) bool{
		var user model.User
		authHost, authPort := GetProfileHostAndPort()
		resp, err := http.Get("http://"+authHost+":"+authPort+"/get-by-id/" + Uint2String(userId))
		if err != nil {
			fmt.Println(err)
			return false
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil{
			fmt.Println(err)
		}
		fmt.Println("Body: ", body)
		defer resp.Body.Close()
		err = json.Unmarshal(body, &user)
		fmt.Println(user.Roles)
		return false
	}
}

func LoggedUserHasPermission(permission string, r *http.Request) bool{
	loggedUserHasPermission:= userHasPermission(GetLoggedUserIDFromToken(r))
	return loggedUserHasPermission(permission)
}
