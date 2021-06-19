package model

import (
	"encoding/json"
)

type Mappable interface {
	ToMap() map[string]interface{}
}

func ToMap(obj Mappable) map[string]interface{} {
	var res map[string]interface{}
	objJson, _ := json.Marshal(obj)
	json.Unmarshal([]byte(objJson), &res)
	return res
}