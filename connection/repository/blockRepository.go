package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"nistagram/connection/model"
)

type BlockRepository struct {
	DatabaseDriver *neo4j.Driver
}

func (repo *BlockRepository) CreateBlock(id1, id2 uint) (*model.Block, bool) {
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	block := model.Block{
		PrimaryProfile:    id1,
		SecondaryProfile:  id2,
	}
	resultingBlock, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
			result, err := transaction.Run(
				"MATCH (a:Profile), (b:Profile) \n" +
					"WHERE a.profileID = $primary AND b.profileID = $secondary \n" +
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
		var ret = model.Block{
			PrimaryProfile:    id1,
			SecondaryProfile:  id2,
		}
		return ret, err
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	var ret = resultingBlock.(model.Block)
	return &ret, true
}

func (repo *BlockRepository) SelectBlock(id1, id2 uint) (*model.Block, bool) {
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	block := model.Block{
		PrimaryProfile:    id1,
		SecondaryProfile:  id2,
	}
	resultingBlock, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:BLOCKED]->(b:Profile) \n"+
				"WHERE a.profileID = $primary AND b.profileID = $secondary \n"+
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
		var ret = model.Block{
			PrimaryProfile:    id1,
			SecondaryProfile:  id2,
		}
		return ret, err
	})
	if err != nil {
		fmt.Println(err.Error())
		return nil, false
	}
	var ret = resultingBlock.(model.Block)
	return &ret, true
}

func (repo *BlockRepository) DeleteBlock(followerId, profileId uint) (*model.Block, bool) {
	block, ok := repo.SelectBlock(followerId, profileId)
	if !ok {
		return nil, false
	}
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()
	_, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		return transaction.Run(
			"MATCH (a:Profile)-[e:BLOCKED]->(b:Profile) \n"+
				"WHERE a.profileID = $primary AND b.profileID = $secondary \n"+
				"DELETE e",
			block.ToMap())
	})
	if err != nil {
		return nil, false
	}
	return block, true
}

func (repo *BlockRepository) GetBlockedProfiles(id uint, directed bool) *[]uint {
	block := model.Block{
		PrimaryProfile:   id,
		SecondaryProfile: 0,
	}
	session := (*repo.DatabaseDriver).NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead})
	defer session.Close()
	profileIDs, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"MATCH (a:Profile)-[e:BLOCKED]->(b:Profile) \n" +
				"WHERE a.profileID = $primary \n" +
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
					"WHERE b.profileID = $primary \n" +
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
		if err != nil {
			fmt.Println(err.Error())
		} else {
			ret = append(ret, profileIDsInv.([]uint)...)
		}
	}
	return &ret
}
