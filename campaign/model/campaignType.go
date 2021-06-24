package model

type CampaignType int

const (
	ONE_TIME CampaignType = iota
	REPEATABLE
)
