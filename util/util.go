package util

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"net/http"
	"net/smtp"
	"nistagram/profile/model"
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

func GetProfile(ctx context.Context, id uint) *model.Profile {
	span := Tracer.StartSpanFromContext(ctx, "GetProfile-util")
	defer Tracer.FinishSpan(span)
	Tracer.LogFields(span, "service", fmt.Sprintf("servicing id %v\n", id))
	nextCtx := Tracer.ContextWithSpan(ctx, span)
	var p model.Profile
	profileHost, profilePort := GetProfileHostAndPort()
	resp, err := CrossServiceRequest(nextCtx, http.MethodGet,
		GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-id/"+Uint2String(id),
		nil, map[string]string{})
	if err != nil {
		Tracer.LogError(span, err)
		fmt.Println(err)
		return nil
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			Tracer.LogError(span, err)
		}
	}(resp.Body)
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		Tracer.LogError(span, err1)
		fmt.Println(err1)
		return nil
	}
	err = json.Unmarshal(body, &p)
	if err != nil {
		Tracer.LogError(span, err)
		fmt.Println(err)
		return nil
	}
	return &p
}

func GetUserPrivileges(id uint) ([]string, bool) {
	authHost, authPort := GetAuthHostAndPort()
	resp, err := CrossServiceRequest(context.Background() ,http.MethodGet,
		GetCrossServiceProtocol()+"://"+authHost+":"+authPort+"/privileges/"+Uint2String(id),
		nil, map[string]string{})
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)
	body, err1 := ioutil.ReadAll(resp.Body)
	if err1 != nil {
		fmt.Println(err1)
		return nil, false
	}
	var privileges []string
	err = json.Unmarshal(body, &privileges)
	if err != nil {
		fmt.Println(err)
		return nil, false
	}
	return privileges, true
}