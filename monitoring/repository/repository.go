package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/monitoring/model"
)

const monitoringCollectionName = "events"
const monitoringDbName = "monitoringdb"

type MonitoringRepository struct {
	Client *mongo.Client
}

func (repo *MonitoringRepository) Create(event *model.Event) error {
	collection := repo.getCollection()
	_, err := collection.InsertOne(context.TODO(), event)
	return err
}

func (repo *MonitoringRepository) getCollection() *mongo.Collection {
	return repo.Client.Database(monitoringDbName).Collection(monitoringCollectionName)
}
