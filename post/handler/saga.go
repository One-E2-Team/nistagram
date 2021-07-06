package handler

import (
	"context"
	"github.com/go-redis/redis/v8"
	"nistagram/util/saga"
)

func (handler *Handler) ChangePrivacyHandler(client *redis.Client, m saga.Message) {
	if m.Action == saga.ActionStart {
		switch m.Functionality {
		case saga.ChangeProfilesPrivacy:
			//test for rollback: sendToReplyChannel(client, &m, saga.ActionError, saga.ProfileService, saga.PostService)
			err := handler.PostService.ChangePrivacy(context.Background(), m.Profile.ID, m.Profile.ProfileSettings.IsPrivate)
			if err != nil{
				saga.SendToReplyChannel(client, &m, saga.ActionError, saga.ProfileService, saga.PostService)
			}else{
				saga.SendToReplyChannel(client, &m, saga.ActionDone, saga.ProfileService, saga.PostService)
			}
		}
	}
}