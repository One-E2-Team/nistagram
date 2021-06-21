package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"nistagram/connection/model"
)

func (repo *Repository) CreateBlock(id1, id2 uint) (*model.BlockEdge, bool) {
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	block := model.BlockEdge{
		PrimaryProfile:    id1,
		SecondaryProfile:  id2,
	}
	resultingBlock, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(
				"MATCH (a:Profile), (b:Profile) \n" +
					"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n" +
					"MERGE (a)-[e:BLOCKED]->(b) \n" +
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
		res := record.Values[0].(dbtype.Relationship)
		fmt.Println(res)
		var ret = model.BlockEdge{
			PrimaryProfile:    id1,
			SecondaryProfile:  id2,
		}
		return ret, err
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	var ret = resultingBlock.(model.BlockEdge)
	return &ret, true
}

func (repo *Repository) SelectBlock(id1, id2 uint) (*model.BlockEdge, bool) {
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	block := model.BlockEdge{
		PrimaryProfile:    id1,
		SecondaryProfile:  id2,
	}
	resultingBlock, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:BLOCKED]->(b:Profile) \n"+
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
		res := record.Values[0].(dbtype.Relationship)
		fmt.Println(res)
		var ret = model.BlockEdge{
			PrimaryProfile:    id1,
			SecondaryProfile:  id2,
		}
		return ret, err
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	var ret = resultingBlock.(model.BlockEdge)
	return &ret, true
}

func (repo *Repository) DeleteBlock(followerId, profileId uint) (*model.BlockEdge, bool) {
	block, ok := repo.SelectBlock(followerId, profileId)
	if !ok {
		return nil, false
	}
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		return transaction.Run(
			"MATCH (a:Profile)-[e:BLOCKED]->(b:Profile) \n"+
				"WHERE a.profileID = $primary AND b.profileID = $secondary AND a.deleted = FALSE AND b.deleted = FALSE \n"+
				"DELETE e",
			block.ToMap())
	})
	if err != nil {
		return nil, false
	}
	return block, true
}

func (repo *Repository) GetBlockedProfiles(id uint, directed bool) *[]uint {
	block := model.BlockEdge{
		PrimaryProfile:   id,
		SecondaryProfile: 0,
	}
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	profileIDs, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:BLOCKED]->(b:Profile) \n" +
				"WHERE a.profileID = $primary AND a.deleted = FALSE AND b.deleted = FALSE \n" +
				"RETURN b",
			block.ToMap())
		var ret []uint
		if err != nil {
			fmt.Println(err.Error())
			return ret, err
		}
		for result.Next() {
			ret = append(ret, uint(result.Record().Values[0].(dbtype.Node).Props["profileID"].(float64)))
		}
		return ret, err
	})
	var ret []uint
	if err != nil {
		fmt.Println(err.Error())
	} else {
		temp := profileIDs.([]uint)
		ret = temp
	}
	if !directed{
		profileIDsInv, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(
				"MATCH (a:Profile)-[e:BLOCKED]->(b:Profile) \n" +
					"WHERE b.profileID = $primary AND a.deleted = FALSE AND b.deleted = FALSE \n" +
					"RETURN a",
				block.ToMap())
			var ret []uint
			if err != nil {
				fmt.Println(err.Error())
				return ret, err
			}
			for result.Next() {
				ret = append(ret, uint(result.Record().Values[0].(dbtype.Node).Props["profileID"].(float64)))
			}
			return ret, err
		})
		if err != nil {
			fmt.Println(err.Error())
		} else {
			ret = append(ret, profileIDsInv.([]uint)...)
		}
	}
	return &ret
}


