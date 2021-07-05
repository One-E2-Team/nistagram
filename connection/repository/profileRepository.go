package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"nistagram/connection/model"
	"nistagram/util"
)

func (repo *Repository) CreateOrUpdateProfile(ctx context.Context, profile model.ProfileVertex) *model.ProfileVertex {
	span := util.Tracer.StartSpanFromContext(ctx, "CreateOrUpdateProfile-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", profile.ProfileID))
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	profileID, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			//"CREATE (n:Profile) SET Profile.profileID = $profileID RETURN Profile",
			"MERGE (n:Profile {profileID: $profileID}) \n" +
				"	ON CREATE SET n += { deleted: $deleted} \n" +
				"	ON MATCH SET n.deleted = $deleted \n" +
				"RETURN n", //kreira ako ne postoji
			profile.ToMap())
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err.Error())
			return 0, err
		}
		record, _ := result.Single()
		res := record.Values[0].(dbtype.Node).Props
		profileID, _ := res["profileID"].(float64)
		return uint(profileID), err
	})
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err.Error())
	}
	id, _ := profileID.(uint)
	ret := model.ProfileVertex{ProfileID: id}
	return &ret
}

func (repo *Repository) FindConnectionDegreeForRecommendation(ctx context.Context, id uint) (*map[uint]*[]uint, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "FindConnectionDegreeForRecommendation-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", id))
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	profile := model.ProfileVertex{ProfileID: id}
	ret, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH path = (a:Profile {profileID: $profileID, deleted: FALSE})-[:FOLLOWS*2..4]->(b:Profile) \n" +
				"WHERE b.deleted = FALSE and not exists( (a)-[:FOLLOWS]-(b) ) \n" +
				"RETURN length(path) as Degree, b",
			profile.ToMap())
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err.Error())
			return nil, err
		}
		var ret = make(map[uint]*[]uint)
		for result.Next() {
			record := result.Record()
			id := uint(record.Values[1].(dbtype.Node).Props["profileID"].(float64))
			deg := uint(record.Values[0].(int64))
			if ret[id] == nil {
				temp := make([]uint, 0)
				ret[id] = &temp
			}
			temp := append(*(ret[id]), deg)
			ret[id] = &temp
		}
		return &ret, nil
	})
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err.Error())
		return nil, false
	}
	if ret == nil {
		util.Tracer.LogError(span, fmt.Errorf("repository anon func returned nill"))
		return nil, false
	}
	data, ok := ret.(*map[uint]*[]uint)
	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("casting error"))
		return nil, false
	}
	return data, true
}