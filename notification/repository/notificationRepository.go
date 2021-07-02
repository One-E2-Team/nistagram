package repository

import "go.mongodb.org/mongo-driver/mongo"

const monitoringCollectionName = "messages"
const monitoringDbName = "notificationdb"

type NotificationRepository struct {
	Client *mongo.Client
}