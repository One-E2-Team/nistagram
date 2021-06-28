package service

import (
	"nistagram/connection/repository"
)

type Service struct {
	ConnectionRepository *repository.Repository
}

func contains(s *[]uint, e uint) bool {
	for _, a := range *s {
		if a == e {
			return true
		}
	}
	return false
}

