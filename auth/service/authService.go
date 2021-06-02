package service

import (
	"fmt"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"net/smtp"
	"nistagram/auth/dto"
	"nistagram/auth/model"
	"nistagram/auth/repository"
	"os"
	"strconv"
	"time"
)

type AuthService struct {
	AuthRepository *repository.AuthRepository
}

func (service *AuthService) LogIn(dto dto.LogInDTO) (*model.User, error) {
	user, err := service.AuthRepository.GetUserByUsername(dto.Username)
	if err != nil {
		return nil, err
	}
	if user.IsDeleted {
		return nil, fmt.Errorf("USER_DELETED")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (service *AuthService) SendMail(sendTo string, subject string, mailMessage string) {
	from := os.Getenv("ISA_MAIL_USERNAME")
	password := os.Getenv("ISA_MAIL_PASSWORD")
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

func (service *AuthService) Register(dto dto.RegisterDTO) error {
	pass := hashAndSalt(dto.Password)
	user := model.User{Username: dto.Username, Password: pass,
		Email: dto.Email, Roles: nil, IsDeleted: false, ValidationExpire: time.Now()}
	err := service.AuthRepository.CreateUser(&user)
	return err
}

func (service *AuthService) GetUserByUsername(username string) *model.User {
	user, err := service.AuthRepository.GetUserByUsername(username)
	if err != nil {
		return nil
	}
	return user
}

func (service *AuthService) RequestPassRecovery(username string) error {
	var user *model.User
	user, err := service.AuthRepository.GetUserByUsername(username)
	if err != nil {
		return err
	}
	user.ValidationUid = uuid.NewString()
	user.ValidationExpire = time.Now().Add(20 * time.Minute)
	err = service.AuthRepository.UpdateUser(*user)
	if err != nil {
		return err
	}
	//TODO: change host, port and html page
	var message = "Visit this link in the next 20 minutes to change your password: http://localhost:8000/recovery.html?id=" +
		strconv.FormatUint(uint64(user.ID), 10) + "&str=" + user.ValidationUid
	go service.SendMail(user.Email, "Recovery password", message)
	return nil
}

func (service *AuthService) ChangePassword(dto dto.RecoveryDTO) error {
	user, err := service.AuthRepository.GetUserByID(dto.Id)
	if err != nil {
		return err
	}
	if user.IsDeleted || user.ValidationUid != dto.Uuid || time.Now().Unix() > user.ValidationExpire.Unix() {
		return fmt.Errorf("BAD_RECOVERY_DATA")
	}
	user.ValidationExpire = time.Now()
	user.ValidationUid = ""
	user.Password = hashAndSalt(dto.Password)
	err = service.AuthRepository.UpdateUser(*user)
	if err != nil {
		return err
	}
	return nil
}

func hashAndSalt(pass string) string {
	bytePass := []byte(pass)
	hash, err := bcrypt.GenerateFromPassword(bytePass, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return string(hash)
}
