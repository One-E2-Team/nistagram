package model

type RequestStatus int

const (
	SENT RequestStatus = iota
	ACCEPTED
)

func (e RequestStatus) ToString() string {
	switch e {
	case SENT:
		return "SENT"
	case ACCEPTED:
		return "ACCEPTED"
	default:
		return "NONE"
	}
}
