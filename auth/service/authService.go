package service

import (
	"errors"
	"fmt"
	"net/smtp"
	"nistagram/auth/dto"
	"nistagram/auth/model"
	"nistagram/auth/repository"
	"os"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func (service *AuthService) LogIn(dto dto.LogInDTO) error {
	credentials := model.Credentials{}
	service.AuthRepository.Database.Find(&credentials, "username", dto.Username)
	if credentials.Password != dto.Password{
		return errors.New("BAD")
	}
	return nil
}

func (service *AuthService) sendMail(sendTo string, mailMessage string) {
	from := os.Getenv("ISA_MAIL_USERNAME")
	password := os.Getenv("ISA_MAIL_PASSWORD")
	to := []string{sendTo}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	message := []byte(mailMessage)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Email Sent Successfully!")
}
