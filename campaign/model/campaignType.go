package model

type CampaignType int

const (
	ONE_TIME CampaignType = iota
	REPEATABLE
)

func (e CampaignType) ToString() string {
	switch e {
	case ONE_TIME:
		return "ONE_TIME"
	case REPEATABLE:
		return "REPEATABLE"
	default:
		return "NONE"
	}
}
