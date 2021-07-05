package repository

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"nistagram/connection/model"
	"nistagram/util"
)

func (repo *Repository) GetConnectedProfiles(ctx context.Context, conn model.ConnectionEdge, excludeMuted bool, direction bool) *[]uint {
	span := util.Tracer.StartSpanFromContext(ctx, "GetConnectedProfiles-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v %v\n", conn.PrimaryProfile, conn.SecondaryProfile))
	var additionalSelector string = ""
	if conn.Approved == true {
		additionalSelector += "AND e.approved = $approved "
		if conn.CloseFriend {
			additionalSelector += "AND e.closeFriend = $closeFriend "
		}
		if conn.NotifyPost {
			additionalSelector += "AND e.notifyPost = $notifyPost "
		}
		if conn.NotifyStory {
			additionalSelector += "AND e.notifyStory = $notifyStory "
		}
		if conn.NotifyComment {
			additionalSelector += "AND e.notifyComment = $notifyComment "
		}
	} else {
		return nil
	}
	if excludeMuted {
		additionalSelector += "AND e.muted = FALSE "
	}
	var targetNode, matchNode, matchDescriptor string
	if direction {
		targetNode = "b"
		matchDescriptor = "primary"
		matchNode = "a"
	} else {
		targetNode = "a"
		matchDescriptor = "secondary"
		matchNode = "b"
	}
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	profileIDs, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:FOLLOWS]->(b:Profile) \n"+
				"WHERE " + matchNode + ".profileID = $" + matchDescriptor + " AND " +
					matchNode + ".deleted = FALSE AND " + targetNode + ".deleted = FALSE " + additionalSelector + "\n"+
				"RETURN " + targetNode,
			conn.ToMap())
		var ret []uint
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

func (repo *Repository) SelectOrCreateConnection(ctx context.Context, id1, id2 uint) *model.ConnectionEdge {
	span := util.Tracer.StartSpanFromContext(ctx, "SelectOrCreateConnection-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v %v\n", id1, id2))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	conn, _ := repo.SelectConnection(nextCtx, id1, id2, true)
	return conn
}

func (repo *Repository) SelectConnection(ctx context.Context, id1, id2 uint, doCreate bool) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "SelectConnection-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v %v\n", id1, id2))
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	conn := model.ConnectionEdge{
		PrimaryProfile:    id1,
		SecondaryProfile:  id2,
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyComment:     false,
		ConnectionRequest: false,
		Approved:          false,
	}
	fmt.Println(conn.ToMap())
	resultingConn, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:FOLLOWS]->(b:Profile) \n"+
				"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n"+
				"RETURN e",
			conn.ToMap())
		record, rerr := result.Single()
		if (doCreate != false && rerr != nil) || err != nil {
			fmt.Println("inif")
			connection, err1 := transaction.Run(
				"MATCH (a:Profile), (b:Profile) \n"+
					"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n"+
					"MERGE (a)-[e:FOLLOWS {muted: FALSE, closeFriend: FALSE, notifyPost: FALSE, notifyStory: "+
					"FALSE, notifyComment: FALSE, connectionRequest: FALSE, approved: FALSE}]->(b) \n"+
					"RETURN e",
				conn.ToMap())
			if err1 != nil {
				util.Tracer.LogError(span, err1)
				return conn, err1
			} else {
				record, rerr = connection.Single()
				if rerr != nil {
					util.Tracer.LogError(span, rerr)
					return nil, rerr
				}
			}
		}
		if rerr != nil{
			util.Tracer.LogError(span, rerr)
			return nil, rerr
		}
		res := record.Values[0].(dbtype.Relationship).Props
		fmt.Println(res)
		var ret = model.ConnectionEdge{
			PrimaryProfile:    id1,
			SecondaryProfile:  id2,
			Muted:             res["muted"].(bool),
			CloseFriend:       res["closeFriend"].(bool),
			NotifyPost:        res["notifyPost"].(bool),
			NotifyStory:       res["notifyStory"].(bool),
			NotifyComment:     res["notifyComment"].(bool),
			ConnectionRequest: res["connectionRequest"].(bool),
			Approved:          res["approved"].(bool),
		}
		return ret, err
	})
	fmt.Println("resulting")
	fmt.Println(resultingConn)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err.Error())
		return nil, false
	}
	var ret = resultingConn.(model.ConnectionEdge)
	return &ret, true
}

func (repo *Repository) UpdateConnection(ctx context.Context, conn *model.ConnectionEdge) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdateConnection-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v %v\n", conn.PrimaryProfile, conn.SecondaryProfile))
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	resultingConn, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:FOLLOWS]->(b:Profile) \n"+
				"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n"+
				"SET e.muted = $muted, e.closeFriend = $closeFriend, e.notifyPost = $notifyPost, "+
				"e.notifyStory = $notifyStory, e.notifyComment = $notifyComment, "+
				"e.connectionRequest = $connectionRequest, e.approved = $approved \n"+
				"RETURN e\n",
			conn.ToMap())
		if err != nil {
			util.Tracer.LogError(span, err)
			fmt.Println(err.Error())
			return nil, err
		}
		record, _ := result.Single()
		res := record.Values[0].(dbtype.Relationship).Props
		fmt.Println(res)
		var ret = model.ConnectionEdge{
			PrimaryProfile:    conn.PrimaryProfile,
			SecondaryProfile:  conn.SecondaryProfile,
			Muted:             res["muted"].(bool),
			CloseFriend:       res["closeFriend"].(bool),
			NotifyPost:        res["notifyPost"].(bool),
			NotifyStory:       res["notifyStory"].(bool),
			NotifyComment:     res["notifyComment"].(bool),
			ConnectionRequest: res["connectionRequest"].(bool),
			Approved:          res["approved"].(bool),
		}
		return ret, err
	})
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err.Error())
		return nil, false
	}
	fmt.Println(resultingConn)
	var ret = resultingConn.(model.ConnectionEdge)
	return &ret, true
}

func (repo *Repository) DeleteConnection(ctx context.Context, followerId, profileId uint) (*model.ConnectionEdge, bool) {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteConnection-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v %v\n", followerId, profileId))
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	conn, ok := repo.SelectConnection(nextCtx, followerId, profileId, false)
	if !ok {
		util.Tracer.LogError(span, fmt.Errorf("select connection did not return a result"))
		return nil, false
	}
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		return transaction.Run(
			"MATCH (a:Profile)-[e:FOLLOWS]->(b:Profile) \n"+
				"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n"+
				"DELETE e",
			conn.ToMap())
	})
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, false
	}
	return conn, true
}

func (repo *Repository) GetAllFollowRequests(ctx context.Context, id uint) *[]uint {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllFollowRequests-repository")
	defer util.Tracer.FinishSpan(span)
	util.Tracer.LogFields(span, "repository", fmt.Sprintf("repository call for id %v\n", id))
	conn := model.ConnectionEdge{
		PrimaryProfile:    0,
		SecondaryProfile:  id,
		Muted:             false,
		CloseFriend:       false,
		NotifyPost:        false,
		NotifyStory:       false,
		NotifyComment:     false,
		ConnectionRequest: true,
		Approved:          false,
	}
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	profileIDs, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:FOLLOWS]->(b:Profile) \n"+
				"WHERE b.profileID = $secondary AND e.connectionRequest = $connectionRequest AND e.approved = $approved " +
				"AND a.deleted = FALSE AND b.deleted = FALSE \n"+
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
