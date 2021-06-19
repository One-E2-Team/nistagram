package service


/*
func (service *Service) MessageConnect(followerId, profileId uint) (*model.Connection, bool) {
	connection, ok := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	conn2, ok2 := service.ConnectionRepository.SelectConnection(profileId, followerId, false)
	if !connection.MessageRequest || (!ok || !ok2) {
		return nil, false
	}
	connection.MessageRequest = false
	connection.MessageConnected = true
	conn2.MessageRequest = false
	conn2.MessageConnected = true
	service.ConnectionRepository.UpdateConnection(connection)
	resConnection, ok1 := service.ConnectionRepository.UpdateConnection(conn2)
	if ok1 {
		return resConnection, true
	} else {
		return conn2, false
	}
}

func (service *Service) MessageRequest(followerId, profileId uint) (*model.Connection, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
	connection := service.ConnectionRepository.SelectOrCreateConnection(followerId, profileId)
	if connection.MessageConnected {
		return nil, false
	}
	connection.MessageRequest = true
	conn2 := service.ConnectionRepository.SelectOrCreateConnection(profileId, followerId)
	if conn2.MessageConnected {
		return nil, false
	}
	if !conn2.Approved {
		conn2.MessageRequest = false
	}
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	service.ConnectionRepository.UpdateConnection(conn2)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}
*/

/*
func (service *Service) ToggleNotifyMessage(followerId, profileId uint) (*model.Connection, bool) {
	if service.IsInBlockingRelationship(followerId, profileId) {
		return nil, false
	}
	connection, okSelect := service.ConnectionRepository.SelectConnection(followerId, profileId, false)
	if okSelect && connection == nil {
		return connection, false
	}
	connection.NotifyMessage = !connection.NotifyMessage
	resConnection, ok := service.ConnectionRepository.UpdateConnection(connection)
	if ok {
		return resConnection, true
	} else {
		return connection, false
	}
}
*/