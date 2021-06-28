package repository

import "go.mongodb.org/mongo-driver/mongo"

type MonitoringRepository struct {
	Client *mongo.Client
}
