package util

import (
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/smtp"
	"os"
	"strconv"
	"strings"
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

func Uint2String(input uint) string {
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

func Contains(array []uint, el uint) bool {
	for _, a := range array {
		if a == el {
			return true
		}
	}
	return false
}

func GetStringIDFromMongoID(mongoID primitive.ObjectID) string {
	return strings.Split(mongoID.String(), "\"")[1]
}

func IsFollowed(array []FollowingProfileDTO, el uint) bool {
	for _, a := range array {
		if a.ProfileID == el {
			return true
		}
	}
	return false
}

func IsCloseFriend(array []FollowingProfileDTO, el uint) bool {
	for _, a := range array {
		if a.ProfileID == el && a.CloseFriend{
			return true
		}
	}
	return false
}

type FollowingProfileDTO struct {
	ProfileID        uint     `json:"profileID"`
	CloseFriend 	 bool     `json:"closeFriend"`
}
