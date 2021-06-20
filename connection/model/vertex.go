package model

type ProfileVertex struct {
	ProfileID	uint `json:"profileID"`
	Deleted		bool `json:"deleted"`
}

func (profile *ProfileVertex) ToMap() map[string]interface{} {
	return ToMap(profile)
}
