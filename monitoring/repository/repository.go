package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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

func (repo *MonitoringRepository) GetEventsByCampaignId(campaignId uint) ([]model.Event,error) {
	collection := repo.getCollection()
	filter := bson.D{{"campaignid", campaignId}}
	cursor, err := collection.Find(context.TODO(), filter)

	var events []model.Event

	for cursor.Next(context.TODO()) {
		var event model.Event

		err = cursor.Decode(&event)
		if err != nil{
			return events,err
		}

		events = append(events, event)
	}

	return events,err
}

func (repo *MonitoringRepository) getCollection() *mongo.Collection {
	return repo.Client.Database(monitoringDbName).Collection(monitoringCollectionName)
}
