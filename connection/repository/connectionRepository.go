package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"nistagram/connection/model"
)

type ConnectionRepository struct {
	DatabaseDriver *neo4j.Driver
}

func (repo *ConnectionRepository) CreateProfile(profile model.Profile) *model.Profile {
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	profileID, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			//"CREATE (n:Profile) SET Profile.profileID = $profileID RETURN Profile",
			"MERGE (n:Profile {profileID: $profileID}) RETURN n", //kreira ako ne postoji
			profile.ToMap())
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}

		//if result.Next() {
		//	return result.Record().Values[0], nil
		//}

		record, _ := result.Single()
		res := record.Values[0].(dbtype.Node).Props
		profileID, _ := res["profileID"].(float64)
		return uint(profileID), err
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	id, _ := profileID.(uint)
	ret := model.Profile{ProfileID: id}
	return &ret
}

func (repo *ConnectionRepository) SelectOrCreateConnection(id1, id2 uint) *model.Connection{
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	conn := model.Connection{
		PrimaryProfile:    model.Profile{ProfileID: id1},
		SecondaryProfile:  model.Profile{ProfileID: id2},
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyMessage:     false,
		NotifyComment:     false,
		ConnectionRequest: false,
		Approved:          false,
		MessageRequest:    false,
		MessageConnected:  false,
		Block:             false,
	}
	fmt.Println(conn.ToMap())
	resultingConn, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:FOLLOWS]->(b:Profile) \n" +
				"WHERE a.profileID = $primary AND b.profileID = $secondary \n" +
				"RETURN e\n",
			conn.ToMap())
		if err != nil || result.Record() == nil {
			connection, err1 := transaction.Run(
				"MATCH (a:Profile), (b:Profile) \n" +
					"WHERE a.profileID = $primary AND b.profileID = $secondary \n" +
					"MERGE (a)-[e:FOLLOWS {muted: FALSE, closeFriend: FALSE, notifyPost: FALSE, notifyStory: " +
					"FALSE, notifyMessage: FALSE, notifyComment: FALSE, connectionRequest: FALSE, approved: FALSE, " +
					"messageRequest: FALSE, messageConnected: FALSE, block: FALSE}]->(b) \n" +
					"RETURN e",
				conn.ToMap())
			if err1 != nil {
				return conn, err1
			} else {
				result = connection
			}
		}
		record, _ := result.Single()
		res := record.Values[0].(dbtype.Node).Props
		fmt.Println(res)
		return conn, err
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resultingConn)
	//id, _ := profileID.(uint)
	//fmt.Println(id)
	ret := model.Connection{}
	return &ret
}