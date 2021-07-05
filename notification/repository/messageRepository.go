package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"nistagram/notification/model"
	"nistagram/util"
)

const messagesCollectionName = "messages"
const notificationDbName = "notificationdb"

func (repo *Repository) CreateMessage(ctx context.Context, message *model.Message) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateMessage-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v %v\n", message.SenderId, message.ReceiverId))
	collection := repo.getCollection()
	_, err := collection.InsertOne(context.TODO(), message)
	return err
}

func (repo *Repository) Seen(ctx context.Context, messageId string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "Seen-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", messageId))
	collection := repo.getCollection()
	idPrimitive, err := primitive.ObjectIDFromHex(messageId)
	if err != nil{
		util.Tracer.LogError(span, err)
		return err
	}
	filter := bson.D{{"_id", idPrimitive}}
	update := bson.D{
		{"$set", bson.D{
			{"seen", true},
		}},
	}

	_, err = collection.UpdateOne(context.TODO(), filter, update)

	return err
}

func (repo *Repository) GetAllMessages(ctx context.Context, firstId uint, secondId uint) ([]model.Message,error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllMessages-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v %v\n", firstId, secondId))
	var ret []model.Message

	collection := repo.getCollection()
	filter := bson.D{{"senderid", firstId}, {"receiverid", secondId}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil,err
	}

	for cursor.Next(context.TODO()) {
		var message model.Message
		err = cursor.Decode(&message)
		if err != nil{
			util.Tracer.LogError(span, err)
			return nil,err
		}
		ret = append(ret, message)
	}

	filter = bson.D{{"senderid", secondId}, {"receiverid", firstId}}

	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil,err
	}

	for cursor.Next(context.TODO()) {
		var message model.Message
		err = cursor.Decode(&message)
		if err != nil{
			util.Tracer.LogError(span, err)
			return nil,err
		}
		ret = append(ret, message)
	}

	return ret,err
}

func (repo *Repository) GetConnectedProfileIds(ctx context.Context, profileId uint) ([]uint,error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetConnectedProfileIds-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", profileId))
	var ret []uint
	collection := repo.getCollection()
	filter := bson.D{{"senderid", profileId}}
	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil,err
	}
	for cursor.Next(context.TODO()) {
		var message model.Message
		err = cursor.Decode(&message)
		if err != nil{
			util.Tracer.LogError(span, err)
			return nil,err
		}
		if !util.Contains(ret, message.ReceiverId){
			ret = append(ret, message.ReceiverId)
		}
	}
	filter = bson.D{{"receiverid", profileId}}
	cursor, err = collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil,err
	}
	for cursor.Next(context.TODO()) {
		var message model.Message
		err = cursor.Decode(&message)
		if err != nil{
			util.Tracer.LogError(span, err)
			return nil,err
		}
		if !util.Contains(ret, message.SenderId){
			ret = append(ret, message.SenderId)
		}
	}

	return ret,nil
}

func (repo *Repository) GetMessageById(ctx context.Context, messageId string) (model.Message,error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMessageById-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", messageId))
	var ret model.Message

	messId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil{
		util.Tracer.LogError(span, err)
		return ret, err
	}

	collection := repo.getCollection()
	filter := bson.D{{"_id", messId}}

	cursor, err := collection.Find(context.TODO(), filter)
	if err != nil {
		util.Tracer.LogError(span, err)
		return ret,err
	}

	for cursor.Next(context.TODO()) {
		err = cursor.Decode(&ret)
	}

	return ret,err
}

func (repo *Repository) DeleteMessage(ctx context.Context, messageId string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteMessage-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", messageId))

	messId, err := primitive.ObjectIDFromHex(messageId)
	if err != nil{
		util.Tracer.LogError(span, err)
		return  err
	}

	collection := repo.getCollection()
	filter := bson.D{{"_id", messId}}

	_, err = collection.DeleteOne(context.TODO(), filter)

	return err
}

func (repo *Repository) getCollection() *mongo.Collection {
	return repo.Client.Database(notificationDbName).Collection(messagesCollectionName)
}
