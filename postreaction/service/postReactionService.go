package service

import (
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	postModel "nistagram/post/model"
	"nistagram/postreaction/dto"
	"nistagram/postreaction/model"
	"nistagram/postreaction/repository"
	"nistagram/util"
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
		return fmt.Errorf("Cannot react on story!")
	}
	reaction := model.Reaction{ReactionType: reactionType, PostID: postID, ProfileID: loggedUserID}
	return service.PostReactionRepository.ReactOnPost(&reaction)
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
		return fmt.Errorf("Cannot comment on story!")
	}
	comment := model.Comment{PostID: commentDTO.PostID, ProfileID: loggedUserID,
		Content: commentDTO.Content, Time: time.Now()}
	return service.PostReactionRepository.CommentPost(&comment)
}

func (service *PostReactionService) ReportPost(postID string, reason string) error {
	_, err := getPost(postID)
	if err != nil {
		return err
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
		util.CrossServiceProtocol+"://"+postHost+":"+postPort+"/posts",
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
		util.CrossServiceProtocol+"://"+profileHost+":"+profilePort+"/get-by-ids",
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
		util.CrossServiceProtocol+"://"+postHost+":"+postPort+"/posts",
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
		util.CrossServiceProtocol+"://"+postHost+":"+postPort+"/post/"+postID,
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
