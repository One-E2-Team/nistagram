package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/monitoring/model"
	"nistagram/util"
)

const monitoringCollectionName = "events"
const monitoringDbName = "monitoringdb"

type MonitoringRepository struct {
	Client *mongo.Client
}

func (repo *MonitoringRepository) Create(ctx context.Context, event *model.Event) error {
	span := util.Tracer.StartSpanFromContext(ctx, "Create-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for event creation"))
	collection := repo.getCollection()
	_, err := collection.InsertOne(context.TODO(), event)
	return err
}

func (repo *MonitoringRepository) GetEventsByCampaignId(ctx context.Context, campaignId uint) ([]model.Event,error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetEventsByCampaignId-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", campaignId))
	collection := repo.getCollection()
	filter := bson.D{{"campaignid", campaignId}}
	cursor, err := collection.Find(context.TODO(), filter)

	var events []model.Event

	for cursor.Next(context.TODO()) {
		var event model.Event

		err = cursor.Decode(&event)
		if err != nil{
			util.Tracer.LogError(span, err)
			return events,err
		}

		events = append(events, event)
	}

	return events,err
}

func (repo *MonitoringRepository) getCollection() *mongo.Collection {
	return repo.Client.Database(monitoringDbName).Collection(monitoringCollectionName)
}
