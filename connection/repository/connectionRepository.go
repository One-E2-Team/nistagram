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
	conn, _ := repo.SelectConnection(id1, id2, true)
	return conn
}

func (repo *ConnectionRepository) SelectConnection(id1, id2 uint, doCreate bool) (*model.Connection, bool){
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	conn := model.Connection{
		PrimaryProfile:    id1,
		SecondaryProfile:  id2,
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
				"RETURN e",
			conn.ToMap())
		if doCreate == false && err == nil && result.Record() != nil {
			return nil, err
		}
		if err != nil || result.Record() == nil {
			connection, err1 := transaction.Run(
				"MATCH (a:Profile), (b:Profile) \n"+
					"WHERE a.profileID = $primary AND b.profileID = $secondary \n"+
					"MERGE (a)-[e:FOLLOWS {muted: FALSE, closeFriend: FALSE, notifyPost: FALSE, notifyStory: "+
					"FALSE, notifyMessage: FALSE, notifyComment: FALSE, connectionRequest: FALSE, approved: FALSE, "+
					"messageRequest: FALSE, messageConnected: FALSE, block: FALSE}]->(b) \n"+
					"RETURN e",
				conn.ToMap())
			if err1 != nil {
				return conn, err1
			} else {
				result = connection
			}
		}
		record, _ := result.Single()
		res := record.Values[0].(dbtype.Relationship).Props
		fmt.Println(res)
		var ret = model.Connection{
			PrimaryProfile:    id1,
			SecondaryProfile:  id2,
			Muted:             res["muted"].(bool),
			CloseFriend:       res["closeFriend"].(bool),
			NotifyPost:        res["notifyPost"].(bool),
			NotifyStory:       res["notifyStory"].(bool),
			NotifyMessage:     res["notifyMessage"].(bool),
			NotifyComment:     res["notifyComment"].(bool),
			ConnectionRequest: res["connectionRequest"].(bool),
			Approved:          res["approved"].(bool),
			MessageRequest:    res["messageRequest"].(bool),
			MessageConnected:  res["messageConnected"].(bool),
			Block:             res["block"].(bool),
		}
		return ret, err
	})
	if doCreate == false && err == nil {
		return nil, true
	}
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(resultingConn)
	var ret = resultingConn.(model.Connection)
	return &ret, true
}

func (repo *ConnectionRepository) UpdateConnection(conn *model.Connection) (*model.Connection, bool){
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	resultingConn, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:FOLLOWS]->(b:Profile) \n" +
				"WHERE a.profileID = $primary AND b.profileID = $secondary \n" +
				"SET e.muted = $muted, e.closeFriend = $closeFriend, e.notifyPost = $notifyPost, " +
				"e.notifyStory = $notifyStory, e.notifyMessage = $notifyMessage, e.notifyComment = $notifyComment, " +
				"e.connectionRequest = $connectionRequest, e.approved = $approved, e.messageRequest = $messageRequest, " +
				"e.messageConnected = $messageConnected, e.block = $block \n" +
				"RETURN e\n",
			conn.ToMap())
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		record, _ := result.Single()
		res := record.Values[0].(dbtype.Relationship).Props
		fmt.Println(res)
		var ret = model.Connection{
			PrimaryProfile:    conn.PrimaryProfile,
			SecondaryProfile:  conn.SecondaryProfile,
			Muted:             res["muted"].(bool),
			CloseFriend:       res["closeFriend"].(bool),
			NotifyPost:        res["notifyPost"].(bool),
			NotifyStory:       res["notifyStory"].(bool),
			NotifyMessage:     res["notifyMessage"].(bool),
			NotifyComment:     res["notifyComment"].(bool),
			ConnectionRequest: res["connectionRequest"].(bool),
			Approved:          res["approved"].(bool),
			MessageRequest:    res["messageRequest"].(bool),
			MessageConnected:  res["messageConnected"].(bool),
			Block:             res["block"].(bool),
		}
		return ret, err
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	fmt.Println(resultingConn)
	var ret = resultingConn.(model.Connection)
	return &ret, true
}