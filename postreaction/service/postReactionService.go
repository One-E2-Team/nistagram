package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	postModel "nistagram/post/model"
	"nistagram/postreaction/dto"
	"nistagram/postreaction/model"
	"nistagram/postreaction/repository"
	"nistagram/util"
	"strings"
	"time"
)

type PostReactionService struct {
	PostReactionRepository *repository.PostReactionRepository
}

func (service *PostReactionService) ReactOnPost(ctx context.Context, reactionDto dto.ReactionDTO, loggedUserID uint, reactionType model.ReactionType) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ReactOnPost-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	post, err := getPost(nextCtx, reactionDto.PostID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if post.PostType == postModel.GetPostType("story") {
		util.Tracer.LogError(span, fmt.Errorf("cannot react on story"))
		return fmt.Errorf("CANNOT REACT ON STORY")
	}
	if post.IsDeleted {
		util.Tracer.LogError(span, fmt.Errorf("cannot react on deleted post"))
		return fmt.Errorf("CANNOT REACT ON DELETED POST")
	}
	if model.GetReactionTypeString(reactionType) == "like_reset" ||
		model.GetReactionTypeString(reactionType) == "dislike_reset" {
		err = service.DeleteReaction(nextCtx, reactionDto.PostID, loggedUserID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
	} else {
		reaction := model.Reaction{ReactionType: reactionType, PostID: reactionDto.PostID, ProfileID: loggedUserID}
		err = service.PostReactionRepository.ReactOnPost(nextCtx, &reaction)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
	}
	if post.IsCampaign {
		go func() {
			event := &dto.EventDTO{EventType: reactionDto.ReactionType, PostId: reactionDto.PostID,
				ProfileId: loggedUserID, CampaignId: reactionDto.CampaignID,
				InfluencerId: reactionDto.InfluencerID, InfluencerUsername: reactionDto.InfluencerUsername}
			if reactionDto.InfluencerID == 0 {
				err = saveToMonitoringMsInfluencer(nextCtx, event)
			}
			err = saveToMonitoringMsTargetGroup(nextCtx, event)
			if err != nil {
				util.Tracer.LogError(span, err)
				fmt.Println("Monitoring bug!")
				fmt.Println(err)
			}
		}()

	}
	return nil
}

func (service *PostReactionService) DeleteReaction(ctx context.Context, postID string, loggedUserID uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteReaction-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	return service.PostReactionRepository.DeleteReaction(nextCtx, postID, loggedUserID)
}

func (service *PostReactionService) CommentPost(ctx context.Context, commentDTO dto.CommentDTO, loggedUserID uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CommentPost-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	post, err := getPost(nextCtx, commentDTO.PostID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if post.PostType == postModel.GetPostType("story") {
		util.Tracer.LogError(span, fmt.Errorf("cannot comment on story"))
		return fmt.Errorf("CANNOT COMMENT ON STORY")
	}
	if post.IsDeleted {
		util.Tracer.LogError(span, fmt.Errorf("cannot comment on deleted post"))
		return fmt.Errorf("CANNOT COMMENT ON DELETED POST")
	}
	if strings.Contains(commentDTO.Content, "@") {
		err = canUsersBeTagged(nextCtx, post.Description, loggedUserID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
	}
	comment := model.Comment{PostID: commentDTO.PostID, ProfileID: loggedUserID,
		Content: commentDTO.Content, Time: time.Now()}
	err = service.PostReactionRepository.CommentPost(nextCtx, &comment)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	if post.IsCampaign {
		go func() {
			event := &dto.EventDTO{EventType: "COMMENT", PostId: commentDTO.PostID,
				ProfileId: loggedUserID, CampaignId: commentDTO.CampaignID,
				InfluencerId: commentDTO.InfluencerID, InfluencerUsername: commentDTO.InfluencerUsername}
			if commentDTO.InfluencerID == 0 {
				err = saveToMonitoringMsInfluencer(nextCtx, event)
			}
			err = saveToMonitoringMsTargetGroup(nextCtx, event)
			if err != nil {
				util.Tracer.LogError(span, err)
				fmt.Println("Monitoring bug!")
				fmt.Println(err)
			}
		}()
	}
	return nil
}

func (service *PostReactionService) ReportPost(ctx context.Context, postID string, reason string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ReportPost-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	post, err := getPost(nextCtx, postID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if post.IsDeleted {
		util.Tracer.LogError(span, fmt.Errorf("cannot report deleted post"))
		return fmt.Errorf("CANNOT REPORT DELETED POST")
	}
	report := model.Report{ID: primitive.NewObjectID(), PostID: postID, Time: time.Now(), Reason: reason,
		IsDeleted: false}
	return service.PostReactionRepository.ReportPost(nextCtx, &report)
}

func (service *PostReactionService) GetMyReactions(ctx context.Context, reactionType model.ReactionType, loggedUserID uint) ([]postModel.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMyReactions-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	reactions, err := service.PostReactionRepository.GetProfileReactions(nextCtx, reactionType, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	ret := make([]string, 0)
	for _, value := range reactions {
		ret = append(ret, value.PostID)
	}
	postHost, postPort := util.GetPostHostAndPort()
	postBody, _ := json.Marshal(map[string][]string{
		"ids": ret,
	})
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/posts",
		postBody, map[string]string{"Content-Type": "application/json;"})

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	if resp.StatusCode != 200 {
		util.Tracer.LogError(span, fmt.Errorf("bad post id"))
		return nil, fmt.Errorf("BAD_POST_ID")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var posts []postModel.Post
	if err = json.Unmarshal(body, &posts); err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return posts, nil
}

func (service *PostReactionService) GetReactionTypes(ctx context.Context, profileID uint, postIDs []string) []string {
	span := util.Tracer.StartSpanFromContext(ctx, "GetReactionTypes-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	ret := make([]string, 0)
	for _, value := range postIDs {
		ret = append(ret, service.PostReactionRepository.GetReactionType(nextCtx, profileID, value))
	}
	return ret
}

func (service *PostReactionService) GetAllReports(ctx context.Context) ([]dto.ShowReportDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllReports-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	var ret []dto.ShowReportDTO

	reports, err := service.PostReactionRepository.GetAllReports(nextCtx)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var repIds []string

	for i := 0; i < len(reports); i++ {
		repIds = append(repIds, reports[i].PostID)
	}

	posts, err := getPostsByPostsIds(nextCtx, repIds)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	for i := 0; i < len(posts); i++ {
		retDto := dto.ShowReportDTO{ReportID: util.GetStringIDFromMongoID(reports[i].ID),
			PostID: util.GetStringIDFromMongoID(posts[i].ID), Reason: reports[i].Reason,
			PublisherId: posts[i].PublisherId, PublisherUsername: posts[i].PublisherUsername,
			Medias: posts[i].Medias, Description: posts[i].Description}
		ret = append(ret, retDto)
	}

	return ret, nil
}

func (service *PostReactionService) DeletePostsReports(ctx context.Context, postId string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeletePostsReports-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	reports, err := service.PostReactionRepository.GetReportsByPostId(nextCtx, postId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	for i := 0; i < len(reports); i++ {
		err = service.PostReactionRepository.DeleteReport(nextCtx, reports[i].ID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
	}

	return nil
}

func (service *PostReactionService) GetAllReactions(ctx context.Context, postID string) ([]string, []string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllReactions-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	likes, dislikes, err := service.PostReactionRepository.GetAllReactions(nextCtx, postID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, nil, err
	}
	likesUsernames, err := getProfileUsernamesByIDs(nextCtx, likes)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, nil, err
	}
	dislikesUsernames, err := getProfileUsernamesByIDs(nextCtx, dislikes)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, nil, err
	}
	return likesUsernames, dislikesUsernames, nil
}

func (service *PostReactionService) GetAllComments(ctx context.Context, postID string) ([]dto.ResponseCommentDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetAllComments-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	comments, err := service.PostReactionRepository.GetAllComments(nextCtx, postID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	profileIDs := make([]uint, 0)
	for _, value := range comments {
		profileIDs = append(profileIDs, value.ProfileID)
	}
	commentUsernames, err := getProfileUsernamesByIDs(nextCtx, profileIDs)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	ret := make([]dto.ResponseCommentDTO, 0)
	if len(commentUsernames) != len(comments) {
		util.Tracer.LogError(span, fmt.Errorf("bad slice sizes"))
		return nil, fmt.Errorf("BAD_SLICE_SIZES")
	}
	for i, value := range comments {
		ret = append(ret, dto.ResponseCommentDTO{Content: value.Content, Username: commentUsernames[i]})
	}
	return ret, nil
}

func getProfileUsernamesByIDs(ctx context.Context, profileIDs []uint) ([]string, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getProfileUsernamesByIDs-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	type data struct {
		Ids []string `json:"ids"`
	}
	req := make([]string, 0)
	for _, value := range profileIDs {
		req = append(req, util.Uint2String(value))
	}
	bodyData := data{Ids: req}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-ids",
		jsonBody, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var ret []string

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err = json.Unmarshal(body, &ret); err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	return ret, nil
}

func getPostsByPostsIds(ctx context.Context, postsIds []string) ([]dto.PostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getPostsByPostsIds-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	var ret []dto.PostDTO
	type data struct {
		Ids []string `json:"ids"`
	}
	bodyData := data{Ids: postsIds}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/posts",
		jsonBody, map[string]string{})

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func getPost(ctx context.Context, postID string) (*postModel.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getPost-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/post/"+postID,
		nil, map[string]string{})
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		util.Tracer.LogError(span, fmt.Errorf("bad post id"))
		return nil, fmt.Errorf("BAD_POST_ID")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var post postModel.Post

	if err = json.Unmarshal(body, &post); err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return &post, nil
}

func canUsersBeTagged(ctx context.Context, description string, publisherId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "canUsersBeTagged-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	followingProfiles, err := getFollowingProfiles(nextCtx, publisherId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	descriptionParts := strings.Split(description, " ")
	for i := 0; i < len(descriptionParts); i++ {
		if strings.Contains(descriptionParts[i], "@") {
			taggedUsername := descriptionParts[i][1:len(descriptionParts[i])]
			var taggedProfile dto.ProfileDto
			if resp, err := getProfileByUsername(nextCtx, taggedUsername); err != nil {
				return err
			} else {
				body, _ := io.ReadAll(resp.Body)
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(resp.Body)
				if err := json.Unmarshal(body, &taggedProfile); err != nil {
					util.Tracer.LogError(span, err)
					return err
				}
			}

			if !taggedProfile.ProfileSettings.CanBeTagged {
				util.Tracer.LogError(span, fmt.Errorf("%s can't be tagged", taggedProfile.Username))
				return errors.New(taggedProfile.Username + " can't be tagged!")
			}

			if !util.IsFollowed(followingProfiles, taggedProfile.ProfileId) {
				util.Tracer.LogError(span, fmt.Errorf("%s is not followed by this profile", taggedProfile.Username))
				return errors.New(taggedProfile.Username + " is not followed by this profile!")
			}
		}
	}
	return nil
}

func getFollowingProfiles(ctx context.Context, loggedUserId uint) ([]util.FollowingProfileDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getFollowingProfiles-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/connection/following/show/"+util.Uint2String(loggedUserId),
		nil, map[string]string{})

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	var followingProfiles []util.FollowingProfileDTO

	err = json.NewDecoder(resp.Body).Decode(&followingProfiles)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	return followingProfiles, err
}

func getProfileByUsername(ctx context.Context, username string) (*http.Response, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getProfileByUsername-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get/"+username,
		nil, map[string]string{})
	return resp, err
}

func saveToMonitoringMsInfluencer(ctx context.Context, eventDto *dto.EventDTO) error {
	span := util.Tracer.StartSpanFromContext(ctx, "saveToMonitoringMsInfluencer-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	body, _ := json.Marshal(map[string]string{
		"eventType":          eventDto.EventType,
		"postId":             eventDto.PostId,
		"profileId":          util.Uint2String(eventDto.ProfileId),
		"campaignId":         util.Uint2String(eventDto.CampaignId),
		"influencerId":       util.Uint2String(eventDto.InfluencerId),
		"influencerUsername": eventDto.InfluencerUsername,
	})
	monitoringHost, monitoringPort := util.GetMonitoringHostAndPort()
	_, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+monitoringHost+":"+monitoringPort+"/influencer",
		body, map[string]string{"Content-Type": "application/json;"})
	return err
}

func saveToMonitoringMsTargetGroup(ctx context.Context, eventDto *dto.EventDTO) error {
	span := util.Tracer.StartSpanFromContext(ctx, "saveToMonitoringMsTargetGroup-service")
	defer util.Tracer.FinishSpan(span)
	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	body, _ := json.Marshal(map[string]string{
		"eventType":          eventDto.EventType,
		"postId":             eventDto.PostId,
		"profileId":          util.Uint2String(eventDto.ProfileId),
		"campaignId":         util.Uint2String(eventDto.CampaignId),
		"influencerId":       util.Uint2String(eventDto.InfluencerId),
		"influencerUsername": eventDto.InfluencerUsername,
	})
	monitoringHost, monitoringPort := util.GetMonitoringHostAndPort()
	_, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+monitoringHost+":"+monitoringPort+"/target-group",
		body, map[string]string{"Content-Type": "application/json;"})
	return err
}
