package model

type ProfileType int

const (
	ADMIN ProfileType = iota
	REGULAR
	AGENT
)
