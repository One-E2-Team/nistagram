package service

import (
	//"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"net/smtp"
	"nistagram/auth/dto"
	"nistagram/auth/model"
	"nistagram/auth/repository"
	"os"
	"time"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func (service *AuthService) LogIn(dto dto.LogInDTO) error {
	user := model.User{}
	service.AuthRepository.Database.Find(&user, "username", dto.Username)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return nil
}

func (service *AuthService) SendMail(sendTo string, mailMessage string) {
	from := os.Getenv("ISA_MAIL_USERNAME")
	password := os.Getenv("ISA_MAIL_PASSWORD")
	to := []string{sendTo}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	message := []byte(mailMessage)
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("Email Sent Successfully!")
}

func (service *AuthService) Register(dto dto.RegisterDTO) error {
	pass := hashAndSalt([]byte(dto.Password))
	user := model.User{Username: dto.Username, Password: pass,
		Email: dto.Email, Roles: nil, IsDeleted: false, ValidationExpire: time.Now()}
	err := service.AuthRepository.CreateUser(&user)
	if err != nil {
		fmt.Println(err.Error())
	}
	return err
}

func hashAndSalt(pass []byte) string {
	hash, err := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(hash)
}
