package model

type RequestStatus int

const (
	SENT RequestStatus = iota
	ACCEPTED
	DECLINED
)

func (e RequestStatus) ToString() string {
	switch e {
	case SENT:
		return "SENT"
	case ACCEPTED:
		return "ACCEPTED"
	case DECLINED:
		return "DECLINED"
	default:
		return "NONE"
	}
}
