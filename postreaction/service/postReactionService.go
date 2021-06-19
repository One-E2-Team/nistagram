package service

import (
	"encoding/json"
	"fmt"
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
	_, err := postExists(postID)
	if err != nil {
		return err
	}
	reaction := model.Reaction{ReactionType: reactionType, PostID: postID, ProfileID: loggedUserID}
	return service.PostReactionRepository.ReactOnPost(&reaction)
}

func (service *PostReactionService) DeleteReaction(postID string, loggedUserID uint) error {
	return service.PostReactionRepository.DeleteReaction(postID, loggedUserID)
}

func (service *PostReactionService) CommentPost(commentDTO dto.CommentDTO, loggedUserID uint) error {
	_, err := postExists(commentDTO.PostID)
	if err != nil {
		return err
	}
	comment := model.Comment{PostID: commentDTO.PostID, ProfileID: loggedUserID,
		Content: commentDTO.Content, Time: time.Now()}
	return service.PostReactionRepository.CommentPost(&comment)
}

func (service *PostReactionService) ReportPost(postID string, reason string) error {
	_, err := postExists(postID)
	if err != nil {
		return err
	}
	report := model.Report{PostID: postID, Time: time.Now(), Reason: reason}
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
		util.CrossServiceProtocol+"://"+postHost+":"+postPort+"/get-collection/post",
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

	for i := 0; i < len(reports); i++{
		repIds = append(repIds, reports[i].PostID)
	}

	posts, err := getPostsByPostsIds(repIds)
	if err != nil{
		return nil, err
	}

	for i := 0; i < len(posts); i++{
		retDto := dto.ShowReportDTO{ReportID: util.GetStringIDFromMongoID(reports[i].ID),
			PostID: util.GetStringIDFromMongoID(posts[i].ID), Reason: reports[i].Reason,
		PublisherId: posts[i].PublisherId, PublisherUsername: posts[i].PublisherUsername,
		Medias: posts[i].Medias, Description: posts[i].Description}
		ret = append(ret, retDto)
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
	if err != nil{
		return nil, err
	}
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
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
func postExists(postID string) (bool, error) {
	postHost, postPort := util.GetPostHostAndPort()
	resp, err := util.CrossServiceRequest(http.MethodGet,
		util.CrossServiceProtocol+"://"+postHost+":"+postPort+"/collection/post/"+postID,
		nil, map[string]string{})
	if err != nil {
		return false, err
	}
	if resp.StatusCode != 200 {
		return false, fmt.Errorf("BAD_POST_ID")
	}
	return true, nil
}
