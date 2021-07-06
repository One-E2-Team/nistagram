package handler

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"nistagram/util/saga"
)

func (handler *Handler) ChangePrivacyRollbackHandler(client *redis.Client, m saga.Message) {
	if m.Action == saga.ActionRollback {
		switch m.Functionality{
		case saga.ChangeProfilesPrivacy:
			profile := m.Profile
			profile.ProfileSettings.IsPrivate = !profile.ProfileSettings.IsPrivate
			err := handler.ProfileService.ProfileRepository.UpdateProfileSettings(context.Background(), profile.ProfileSettings)
			if err != nil{
				fmt.Println(err)
			}
			saga.SendToReplyChannel(client, &m, saga.ActionError, saga.ProfileChannel, saga.ProfileChannel)
		}
	}
}