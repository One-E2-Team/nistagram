package service

import (
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

func (service *PostReactionService) ReactOnPost(postID string, loggedUserID uint, reactionType model.ReactionType) error {
	post, err := getPost(postID)
	if err != nil {
		return err
	}
	if post.PostType == postModel.GetPostType("story") {
		return fmt.Errorf("CANNOT REACT ON STORY")
	}
	if post.IsDeleted {
		return fmt.Errorf("CANNOT REACT ON DELETED POST")
	}
	reaction := model.Reaction{ReactionType: reactionType, PostID: postID, ProfileID: loggedUserID}
	err = service.PostReactionRepository.ReactOnPost(&reaction)
	if err != nil{
		return err
	}
	if post.IsCampaign{
		go func() {
			event := &dto.EventDTO{Type: model.GetReactionTypeString(reactionType),
				PostId: post.ID.Hex(), ProfileId: loggedUserID}
			err = saveToMonitoringMs(event)
			if err != nil {
				fmt.Println("Monitoring bug!")
				fmt.Println(err)
			}
		}()

	}
	return nil
}

func (service *PostReactionService) DeleteReaction(postID string, loggedUserID uint) error {
	return service.PostReactionRepository.DeleteReaction(postID, loggedUserID)
}

func (service *PostReactionService) CommentPost(commentDTO dto.CommentDTO, loggedUserID uint) error {
	post, err := getPost(commentDTO.PostID)
	if err != nil {
		return err
	}
	if post.PostType == postModel.GetPostType("story") {
		return fmt.Errorf("CANNOT COMMENT ON STORY")
	}
	if post.IsDeleted {
		return fmt.Errorf("CANNOT COMMENT ON DELETED POST")
	}
	if strings.Contains(commentDTO.Content, "@") {
		err = canUsersBeTagged(post.Description, loggedUserID)
		if err != nil {
			return err
		}
	}
	comment := model.Comment{PostID: commentDTO.PostID, ProfileID: loggedUserID,
		Content: commentDTO.Content, Time: time.Now()}
	return service.PostReactionRepository.CommentPost(&comment)
}

func (service *PostReactionService) ReportPost(postID string, reason string) error {
	post, err := getPost(postID)
	if err != nil {
		return err
	}
	if post.IsDeleted {
		return fmt.Errorf("CANNOT REPORT DELETED POST")
	}
	report := model.Report{ID: primitive.NewObjectID(), PostID: postID, Time: time.Now(), Reason: reason,
		IsDeleted: false}
	return service.PostReactionRepository.ReportPost(&report)
}

func (service *PostReactionService) GetMyReactions(reactionType model.ReactionType, loggedUserID uint) ([]postModel.Post, error) {
	reactions, err := service.PostReactionRepository.GetProfileReactions(reactionType, loggedUserID)
	if err != nil {
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
	resp, err := util.CrossServiceRequest(http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/posts",
		postBody, map[string]string{"Content-Type": "application/json;"})

	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("BAD_POST_ID")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var posts []postModel.Post
	if err = json.Unmarshal(body, &posts); err != nil {
		return nil, err
	}
	return posts, nil
}

func (service *PostReactionService) GetReactionTypes(profileID uint, postIDs []string) []string {
	ret := make([]string, 0)
	for _, value := range postIDs {
		ret = append(ret, service.PostReactionRepository.GetReactionType(profileID, value))
	}
	return ret
}

func (service *PostReactionService) GetAllReports() ([]dto.ShowReportDTO, error) {
	var ret []dto.ShowReportDTO

	reports, err := service.PostReactionRepository.GetAllReports()
	if err != nil {
		return nil, err
	}

	var repIds []string

	for i := 0; i < len(reports); i++ {
		repIds = append(repIds, reports[i].PostID)
	}

	posts, err := getPostsByPostsIds(repIds)
	if err != nil {
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

func (service *PostReactionService) DeletePostsReports(postId string) error {

	reports, err := service.PostReactionRepository.GetReportsByPostId(postId)
	if err != nil {
		return err
	}

	for i := 0; i < len(reports); i++ {
		err = service.PostReactionRepository.DeleteReport(reports[i].ID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (service *PostReactionService) GetAllReactions(postID string) ([]string, []string, error) {
	likes, dislikes, err := service.PostReactionRepository.GetAllReactions(postID)
	if err != nil {
		return nil, nil, err
	}
	likesUsernames, err := getProfileUsernamesByIDs(likes)
	if err != nil {
		return nil, nil, err
	}
	dislikesUsernames, err := getProfileUsernamesByIDs(dislikes)
	if err != nil {
		return nil, nil, err
	}
	return likesUsernames, dislikesUsernames, nil
}

func (service *PostReactionService) GetAllComments(postID string) ([]dto.ResponseCommentDTO, error) {
	comments, err := service.PostReactionRepository.GetAllComments(postID)
	if err != nil {
		return nil, err
	}
	profileIDs := make([]uint, 0)
	for _, value := range comments {
		profileIDs = append(profileIDs, value.ProfileID)
	}
	commentUsernames, err := getProfileUsernamesByIDs(profileIDs)
	if err != nil {
		return nil, err
	}
	ret := make([]dto.ResponseCommentDTO, 0)
	if len(commentUsernames) != len(comments) {
		return nil, fmt.Errorf("BAD_SLICE_SIZES")
	}
	for i, value := range comments {
		ret = append(ret, dto.ResponseCommentDTO{Content: value.Content, Username: commentUsernames[i]})
	}
	return ret, nil
}

func getProfileUsernamesByIDs(profileIDs []uint) ([]string, error) {
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
		return nil, err
	}
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-ids",
		jsonBody, map[string]string{"Content-Type": "application/json;"})
	if err != nil {
		return nil, err
	}

	var ret []string

	body, err := io.ReadAll(resp.Body)
	if err != nil {
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

func getPostsByPostsIds(postsIds []string) ([]dto.PostDTO, error) {
	var ret []dto.PostDTO
	type data struct {
		Ids []string `json:"ids"`
	}
	bodyData := data{Ids: postsIds}
	jsonBody, err := json.Marshal(bodyData)
	if err != nil {
		return nil, err
	}
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/posts",
		jsonBody, map[string]string{})

	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
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
func getPost(postID string) (*postModel.Post, error) {
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+postHost+":"+postPort+"/post/"+postID,
		nil, map[string]string{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("BAD_POST_ID")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	var post postModel.Post

	if err = json.Unmarshal(body, &post); err != nil {
		return nil, err
	}
	return &post, nil
}

func canUsersBeTagged(description string, publisherId uint) error {

	followingProfiles, err := getFollowingProfiles(publisherId)
	if err != nil {
		return err
	}

	descriptionParts := strings.Split(description, " ")
	for i := 0; i < len(descriptionParts); i++ {
		if strings.Contains(descriptionParts[i], "@") {
			taggedUsername := descriptionParts[i][1 : len(descriptionParts[i])]
			var taggedProfile dto.ProfileDto
			if resp, err := getProfileByUsername(taggedUsername); err != nil {
				return err
			} else {
				body, _ := io.ReadAll(resp.Body)
				defer func(Body io.ReadCloser) {
					_ = Body.Close()
				}(resp.Body)
				if err := json.Unmarshal(body, &taggedProfile); err != nil {
					return err
				}
			}

			if !taggedProfile.ProfileSettings.CanBeTagged {
				return errors.New(taggedProfile.Username + " can't be tagged!")
			}

			if !util.IsFollowed(followingProfiles, taggedProfile.ProfileId) {
				return errors.New(taggedProfile.Username + " is not followed by this profile!")
			}
		}
	}
	return nil
}

func getFollowingProfiles(loggedUserId uint) ([]util.FollowingProfileDTO, error) {
	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/connection/following/show/"+util.Uint2String(loggedUserId),
		nil, map[string]string{})

	if err != nil {
		return nil, err
	}

	var followingProfiles []util.FollowingProfileDTO

	err = json.NewDecoder(resp.Body).Decode(&followingProfiles)
	if err != nil{
		return nil, err
	}

	return followingProfiles, err
}

func getProfileByUsername(username string) (*http.Response, error) {
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get/"+username,
		nil, map[string]string{})
	return resp, err
}

func saveToMonitoringMs(eventDto *dto.EventDTO) error{
	body, _ := json.Marshal(map[string]string{
		"type": 		 eventDto.Type,
		"postId":  		 eventDto.PostId,
		"profileId":     util.Uint2String(eventDto.ProfileId),
	})
	monitoringHost, monitoringPort := util.GetMonitoringHostAndPort()
	_, err := util.CrossServiceRequest(http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+monitoringHost+":"+monitoringPort+"/",
		body, map[string]string{"Content-Type": "application/json;"})
	return err
}