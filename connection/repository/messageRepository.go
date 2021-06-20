package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"nistagram/connection/model"
)

func (repo *Repository) CreateOrUpdateMessageRelationship(message model.MessageEdge) (*model.MessageEdge, bool) {
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	resultingBlock, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile), (b:Profile) \n" +
				"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n" +
				"MERGE (a)-[e:MESSAGE]->(b) " +
				"	ON CREATE SET e += { approved: $approved, notifyMessage: $notifyMessage } \n" +
				"	ON MATCH SET e += { approved: $approved, notifyMessage: $notifyMessage } \n" +
				"RETURN e",
			message.ToMap())
		var record *neo4j.Record
		if err != nil {
			return nil, err
		} else {
			record, err = result.Single()
			if err != nil {
				return nil, err
			}
		}
		res := record.Values[0].(dbtype.Relationship).Props
		fmt.Println(res)
		var ret = model.MessageEdge{
			PrimaryProfile:		message.PrimaryProfile,
			SecondaryProfile:	message.SecondaryProfile,
			Approved:			res["approved"].(bool),
			NotifyMessage:		res["notifyMessage"].(bool),
		}
		return ret, err
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	var ret = resultingBlock.(model.MessageEdge)
	return &ret, true
}

func (repo *Repository) SelectMessage(id1, id2 uint) (*model.MessageEdge, bool) {
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	block := model.MessageEdge{
		PrimaryProfile:    id1,
		SecondaryProfile:  id2,
	}
	resultingBlock, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:MESSAGE]->(b:Profile) \n"+
				"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n"+
				"RETURN e",
			block.ToMap())
		var record *neo4j.Record
		if err != nil {
			return nil, err
		} else {
			record, err = result.Single()
			if err != nil {
				return nil, err
			}
		}
		res := record.Values[0].(dbtype.Relationship).Props
		fmt.Println(res)
		var ret = model.MessageEdge{
			PrimaryProfile:		id1,
			SecondaryProfile:	id2,
			Approved: 			res["approved"].(bool),
			NotifyMessage:		res["notifyMessage"].(bool),
		}
		return ret, err
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	var ret = resultingBlock.(model.MessageEdge)
	return &ret, true
}

func (repo *Repository) DeleteMessage(followerId, profileId uint) (*model.MessageEdge, bool) {
	message, ok := repo.SelectMessage(followerId, profileId)
	if !ok {
		return nil, false
	}
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		return transaction.Run(
			"MATCH (a:Profile)-[e:MESSAGE]->(b:Profile) \n"+
				"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n"+
				"DELETE e",
			message.ToMap())
	})
	if err != nil {
		return nil, false
	}
	return message, true
}