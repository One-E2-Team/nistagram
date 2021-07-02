package service

import "nistagram/notification/repository"

type NotificationService struct {
	NotificationRepository *repository.NotificationRepository
}