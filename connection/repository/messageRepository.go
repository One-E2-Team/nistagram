package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"nistagram/connection/model"
	"nistagram/util"
)

func (repo *Repository) CreateOrUpdateMessageRelationship(ctx context.Context, message model.MessageEdge) (*model.MessageEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateOrUpdateMessageRelationship-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", message.PrimaryProfile))
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	resultingBlock, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile), (b:Profile) \n" +
				"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n" +
				"MERGE (a)-[e:MESSAGE]->(b) " +
				"	ON CREATE SET e += { approved: $approved, notifyMessage: $notifyMessage } \n" +
				"	ON MATCH SET e.approved = $approved , e.notifyMessage = $notifyMessage  \n" +
				"RETURN e",
			message.ToMap())
		var record *neo4j.Record
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		} else {
			record, err = result.Single()
			if err != nil {
				util.Tracer.LogError(span, err)
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
		util.Tracer.LogError(span, err)
		fmt.Println(err.Error())
		return nil, false
	}
	var ret = resultingBlock.(model.MessageEdge)
	return &ret, true
}

func (repo *Repository) SelectMessage(ctx context.Context, id1, id2 uint) (*model.MessageEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "SelectMessage-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for ids %v %v\n", id1, id2))
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
			util.Tracer.LogError(span, err)
			return nil, err
		} else {
			record, err = result.Single()
			if err != nil {
				util.Tracer.LogError(span, err)
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
		util.Tracer.LogError(span, err)
		fmt.Println(err.Error())
		return nil, false
	}
	var ret = resultingBlock.(model.MessageEdge)
	return &ret, true
}

func (repo *Repository) DeleteMessage(ctx context.Context, followerId, profileId uint) (*model.MessageEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteMessage-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for ids %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	message, ok := repo.SelectMessage(nextCtx, followerId, profileId)
	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("select message did not return a result"))
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
		util.Tracer.LogError(span, err)
		return nil, false
	}
	return message, true
}

func (repo *Repository) GetAllMessageRequests(ctx context.Context, id uint) *[]uint {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllMessageRequests-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", id))
	conn := model.MessageEdge{
		PrimaryProfile:    0,
		SecondaryProfile:  id,
		Approved:          false,
	}
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	profileIDs, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:MESSAGE]->(b:Profile) \n"+
				"WHERE b.profileID = $secondary AND e.approved = $approved AND a.deleted = FALSE AND b.deleted = FALSE \n"+
				"RETURN a",
			conn.ToMap())
		var ret []uint = make([]uint, 0)
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err.Error())
			return ret, err
		}

		for result.Next() {
			ret = append(ret, uint(result.Record().Values[0].(dbtype.Node).Props["profileID"].(float64)))
		}

		return ret, err
	})
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err.Error())
	}
	ret := profileIDs.([]uint)
	return &ret
}