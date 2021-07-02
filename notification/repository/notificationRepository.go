package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/notification/model"
)

const messagesCollectionName = "messages"
const notificationDbName = "notificationdb"

type NotificationRepository struct {
	Client *mongo.Client
}

func (repo *NotificationRepository) CreateMessage(message *model.Message) error {
	collection := repo.getCollection()
	_, err := collection.InsertOne(context.TODO(), message)
	return err
}

func (repo *NotificationRepository) Seen(messageId string) error {
	collection := repo.getCollection()
	filter := bson.D{{"_id", primitive.ObjectIDFromHex(messageId)}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return err
	}

	for cursor.Next(context.TODO()) {
		var message model.Message
		err = cursor.Decode(&message)
		if err != nil{
			return err
		}
		message.Seen = true
	}

	return err
}

func (repo *NotificationRepository) GetAllMessages(firstId uint, secondId uint) ([]model.Message,error) {
	var ret []model.Message

	collection := repo.getCollection()
	filter := bson.D{{"senderid", firstId}, {"receiverid", secondId}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		return nil,err
	}

	for cursor.Next(context.TODO()) {
		var message model.Message
		err = cursor.Decode(&message)
		if err != nil{
			return nil,err
		}
		ret = append(ret, message)
	}

	filter = bson.D{{"senderid", secondId}, {"receiverid", firstId}}

	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		return nil,err
	}

	for cursor.Next(context.TODO()) {
		var message model.Message
		err = cursor.Decode(&message)
		if err != nil{
			return nil,err
		}
		ret = append(ret, message)
	}

	return ret,err
}

func (repo *NotificationRepository) getCollection() *mongo.Collection {
	return repo.Client.Database(notificationDbName).Collection(messagesCollectionName)
}