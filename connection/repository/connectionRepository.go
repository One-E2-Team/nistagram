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
	fmt.Println(profile.ToMap())
	profileID, err := session.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"CREATE (Profile) SET Profile.profileID = $profileID RETURN Profile",
			profile.ToMap())
		if err != nil {
			fmt.Println(err.Error())
			return 0, err
		}

		//if result.Next() {
		//	return result.Record().Values[0], nil
		//}

		//return nil, result.Err()
		record, _ := result.Single()
		res := record.GetByIndex(0).(dbtype.Node).Props
		profileID, _ := res["profileID"].(float64)
		//profileID, _ := strconv.ParseUint(profileIDstr,10,32)
		fmt.Println(profileID)
		return uint(profileID), err
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	id, _ := profileID.(uint)
	ret := model.Profile{ProfileID: id}
	return &ret
}

/*var res interface{
	Id,Labels,Props
}*/