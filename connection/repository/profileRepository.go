package repository

import (
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/dbtype"
	"nistagram/connection/model"
)

func (repo *Repository) CreateOrUpdateProfile(profile model.ProfileVertex) *model.ProfileVertex {
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
			fmt.Println(err.Error())
			return 0, err
		}
		record, _ := result.Single()
		res := record.Values[0].(dbtype.Node).Props
		profileID, _ := res["profileID"].(float64)
		return uint(profileID), err
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	id, _ := profileID.(uint)
	ret := model.ProfileVertex{ProfileID: id}
	return &ret
}