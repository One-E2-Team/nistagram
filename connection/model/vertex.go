package model

type ProfileVertex struct {
	ProfileID uint `json:"profileID"`
}

func (profile *ProfileVertex) ToMap() map[string]interface{} {
	return ToMap(profile)
}
