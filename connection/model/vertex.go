package model

import "encoding/json"

type Profile struct {
	ProfileID	uint 	`json:"profileID"`
	Username	string	`json:"username"`
}

func (profile *Profile) ToMap() map[string]interface{}{
	var res map[string]interface{}
	profileJson, _ := json.Marshal(profile)
	json.Unmarshal([]byte(profileJson), &res)
	return res
}