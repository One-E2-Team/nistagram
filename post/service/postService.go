package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"net/http"
	"nistagram/post/dto"
	"nistagram/post/model"
	"nistagram/post/repository"
	"nistagram/util"
	"strings"
	"time"
)

type PostService struct {
	PostRepository *repository.PostRepository
}

func (service *PostService) GetPublic(loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	blockedRelationships, err := getProfilesBlockedRelationships(loggedUserID)
	if err != nil {
		return nil, err
	}
	posts, err := service.PostRepository.GetPublic(blockedRelationships)
	if err != nil {
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetProfilesPosts(followingProfiles []util.FollowingProfileDTO, targetUsername string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	profileID, err := getProfileIDByUsername(targetUsername)
	if err != nil {
		return nil, err
	}
	sponsoredPostsDTO, err := getCampaignsWhereUserIsInfluencer(profileID)
	campaignPosts := make([]model.Post, 0)
	influencerIDs := make([]uint, 0)
	for _, sponsoredPostDTO := range sponsoredPostsDTO {
		primitiveID, err := primitive.ObjectIDFromHex(sponsoredPostDTO.PostID)
		if err != nil {
			return nil, err
		}
		post, err := service.PostRepository.Read(primitiveID)
		if err != nil {
			return nil, err
		}
		influencerIDs = append(influencerIDs, sponsoredPostDTO.InfluencerID)
		campaignPosts = append(campaignPosts, post)
	}
	initialSponsoredPosts, err := getReactionsForPosts(campaignPosts, loggedUserID)
	if err != nil {
		return nil, err
	}
	ret := make([]dto.ResponsePostDTO, 0)
	influencerUsernames, err := getProfileUsernamesByIDs(influencerIDs)
	for i, initial := range initialSponsoredPosts {
		ret = append(ret, dto.ResponsePostDTO{
			Post:               initial.Post,
			Reaction:           initial.Reaction,
			CampaignId:         sponsoredPostsDTO[i].CampaignID,
			InfluencerId:       sponsoredPostsDTO[i].InfluencerID,
			InfluencerUsername: influencerUsernames[i],
		})
	}
	posts, err := service.PostRepository.GetProfilesPosts(followingProfiles, targetUsername)
	if err != nil {
		return nil, err
	}
	postsDTO, err := getReactionsForPosts(posts, loggedUserID)
	if err != nil {
		return nil, err
	}
	for _, post := range postsDTO {
		ret = append(ret, post)
	}
	return ret, err
}

func (service *PostService) GetPublicPostByLocation(location string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	blockedRelationships, err := getProfilesBlockedRelationships(loggedUserID)
	if err != nil {
		return nil, err
	}
	posts, err := service.PostRepository.GetPublicPostByLocation(location, blockedRelationships)
	if err != nil {
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetPublicPostByHashTag(hashTag string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	blockedRelationships, err := getProfilesBlockedRelationships(loggedUserID)
	if err != nil {
		return nil, err
	}
	posts, err := service.PostRepository.GetPublicPostByHashTag(hashTag, blockedRelationships)
	if err != nil {
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetMyPosts(loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	sponsoredPostsDTO, err := getCampaignsWhereUserIsInfluencer(loggedUserID)
	campaignPosts := make([]model.Post, 0)
	influencerIDs := make([]uint, 0)
	for _, sponsoredPostDTO := range sponsoredPostsDTO {
		primitiveID, err := primitive.ObjectIDFromHex(sponsoredPostDTO.PostID)
		if err != nil {
			return nil, err
		}
		post, err := service.PostRepository.Read(primitiveID)
		if err != nil {
			return nil, err
		}
		influencerIDs = append(influencerIDs, sponsoredPostDTO.InfluencerID)
		campaignPosts = append(campaignPosts, post)
	}
	initialSponsoredPosts, err := getReactionsForPosts(campaignPosts, loggedUserID)
	if err != nil {
		return nil, err
	}
	ret := make([]dto.ResponsePostDTO, 0)
	influencerUsernames, err := getProfileUsernamesByIDs(influencerIDs)
	for i, initial := range initialSponsoredPosts {
		ret = append(ret, dto.ResponsePostDTO{
			Post:               initial.Post,
			Reaction:           initial.Reaction,
			CampaignId:         sponsoredPostsDTO[i].CampaignID,
			InfluencerId:       sponsoredPostsDTO[i].InfluencerID,
			InfluencerUsername: influencerUsernames[i],
		})
	}
	posts, err := service.PostRepository.GetMyPosts(loggedUserID)
	if err != nil {
		return nil, err
	}
	postsDTO, err := getReactionsForPosts(posts, loggedUserID)
	if err != nil {
		return nil, err
	}
	for _, post := range postsDTO {
		ret = append(ret, post)
	}
	return ret, err
}

func (service *PostService) GetPostsForHomePage(followingProfiles []util.FollowingProfileDTO, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	posts, err := service.PostRepository.GetPostsForHomePage(followingProfiles)
	if err != nil {
		return nil, err
	}
	ret := make([]dto.ResponsePostDTO, 0)
	ret, err = getReactionsForPosts(posts, loggedUserID)
	if err != nil {
		return nil, err
	}
	sponsoredPostsDTO, err := getCampaigns(loggedUserID, followingProfiles)
	if err != nil {
		return nil, err
	}
	campaignPosts := make([]model.Post, 0)
	influencerIDs := make([]uint, 0)
	for _, sponsoredPostDTO := range sponsoredPostsDTO {
		primitiveID, err := primitive.ObjectIDFromHex(sponsoredPostDTO.PostID)
		if err != nil {
			return nil, err
		}
		post, err := service.PostRepository.Read(primitiveID)
		if err != nil {
			return nil, err
		}
		influencerIDs = append(influencerIDs, sponsoredPostDTO.InfluencerID)
		campaignPosts = append(campaignPosts, post)
	}
	initialSponsoredPosts, err := getReactionsForPosts(campaignPosts, loggedUserID)
	if err != nil {
		return nil, err
	}
	influencerUsernames, err := getProfileUsernamesByIDs(influencerIDs)
	for i, initial := range initialSponsoredPosts {
		ret = append(ret, dto.ResponsePostDTO{
			Post:               initial.Post,
			Reaction:           initial.Reaction,
			CampaignId:         sponsoredPostsDTO[i].CampaignID,
			InfluencerId:       sponsoredPostsDTO[i].InfluencerID,
			InfluencerUsername: influencerUsernames[i],
		})
	}
	return ret, err
}

func (service *PostService) CreatePost(postType model.PostType, post dto.PostDto, mediaNames []string, profile dto.ProfileDto) error {
	var medias []model.Media
	for i := 0; i < len(mediaNames); i++ {
		m := model.Media{
			ID:       primitive.NewObjectID(),
			FilePath: mediaNames[i],
			WebSite:  post.Links[i],
		}
		medias = append(medias, m)
	}

	if strings.Contains(post.Description, "@") {
		err := canUsersBeTagged(post.Description, profile.ProfileId)
		if err != nil {
			return err
		}
	}

	newPost := model.Post{ID: primitive.NewObjectID(), PublisherId: profile.ProfileId, PublisherUsername: profile.Username,
		PostType: postType, Medias: medias, PublishDate: time.Now(),
		Description: post.Description, IsHighlighted: post.IsHighlighted, IsCampaign: false,
		IsCloseFriendsOnly: post.IsCloseFriendsOnly, Location: post.Location,
		HashTags: post.HashTags, IsPrivate: profile.ProfileSettings.IsPrivate, IsDeleted: false}

	return service.PostRepository.Create(&newPost)
}

func (service *PostService) ReadPost(id primitive.ObjectID) (model.Post, error) {
	return service.PostRepository.Read(id)
}

func (service *PostService) DeletePost(id primitive.ObjectID) error {
	err := service.PostRepository.Delete(id)
	if err != nil {
		return err
	}

	err = deletePostsReports(id)
	return err
}

func (service *PostService) UpdatePost(id primitive.ObjectID, post dto.PostDto) error {
	return service.PostRepository.Update(id, post)
}

func (service *PostService) DeleteUserPosts(profileId uint) error {
	posts, err := service.PostRepository.GetMyPosts(profileId)
	if err != nil {
		return err
	}
	for i := 0; i < len(posts); i++ {
		err = service.DeletePost(posts[i].ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (service *PostService) ChangeUsername(profileId uint, username string) error {
	return service.PostRepository.ChangeUsername(profileId, username)
}

func (service *PostService) ChangePrivacy(profileId uint, isPrivate bool) error {
	return service.PostRepository.ChangePrivacy(profileId, isPrivate)
}

func (service *PostService) MakeCampaign(postID string, agentID uint) error {
	id, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		return err
	}
	post, err := service.PostRepository.Read(id)
	if err != nil {
		return err
	}
	if post.PublisherId != agentID {
		return fmt.Errorf("BAD_AGENT_ID")
	}
	return service.PostRepository.MakeCampaign(id)
}

func (service *PostService) GetMediaById(mediaId string) (model.Media, error) {
	return service.PostRepository.GetMediaById(mediaId)
}

func deletePostsReports(postId primitive.ObjectID) error {
	postReactionHost, postReactionPort := util.GetPostReactionHostAndPort()
	_, err := util.CrossServiceRequest(context.Background(), http.MethodDelete,
		util.GetCrossServiceProtocol()+"://"+postReactionHost+":"+postReactionPort+"/report/"+util.GetStringIDFromMongoID(postId),
		nil, map[string]string{})
	return err
}

func canUsersBeTagged(description string, publisherId uint) error {
	var followingProfiles []uint

	resp, err := getUserFollowers(publisherId)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err = json.Unmarshal(body, &followingProfiles); err != nil {
		return err
	}

	descriptionParts := strings.Split(description, " ")
	for i := 0; i < len(descriptionParts); i++ {
		if strings.Contains(descriptionParts[i], "@") {
			taggedUsername := descriptionParts[i][1:len(descriptionParts[i])]
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

			if !util.Contains(followingProfiles, taggedProfile.ProfileId) {
				return errors.New(taggedProfile.Username + " is not followed by this profile!")
			}
		}
	}
	return nil
}

func getUserFollowers(loggedUserId uint) (*http.Response, error) {
	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/connection/following/show/"+util.Uint2String(loggedUserId),
		nil, map[string]string{})
	return resp, err
}

func getProfileByUsername(username string) (*http.Response, error) {
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get/"+username,
		nil, map[string]string{})
	return resp, err
}

func getReactionsForPosts(posts []model.Post, profileID uint) ([]dto.ResponsePostDTO, error) {
	if len(posts) == 0 {
		return make([]dto.ResponsePostDTO, 0), nil
	}
	if profileID == 0 {
		ret := make([]dto.ResponsePostDTO, len(posts))
		for i, value := range posts {
			ret[i].Reaction = "none"
			ret[i].Post = value
		}
		return ret, nil
	}
	postIDs := make([]string, 0)
	for _, value := range posts {
		postIDs = append(postIDs, util.GetStringIDFromMongoID(value.ID))
	}
	postReactionHost, postReactionPort := util.GetPostReactionHostAndPort()
	postBody, _ := json.Marshal(map[string][]string{
		"ids": postIDs,
	})
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postReactionHost+":"+postReactionPort+"/get-reaction-types/"+util.Uint2String(profileID),
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

	var reactions []string
	if err = json.Unmarshal(body, &reactions); err != nil {
		return nil, err
	}
	if len(posts) != len(reactions) {
		return nil, fmt.Errorf("BAD_LISTS")
	}
	ret := make([]dto.ResponsePostDTO, 0)

	for i, value := range posts {
		ret = append(ret, dto.ResponsePostDTO{
			Post:     value,
			Reaction: reactions[i],
		})
	}
	return ret, nil
}

func getProfilesBlockedRelationships(loggedProfileId uint) ([]uint, error) {
	connectionHost, connectionPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+connectionHost+":"+connectionPort+"/connection/block/relationships/"+util.Uint2String(loggedProfileId),
		nil, map[string]string{})

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

	var blockedRelationships []uint
	if err = json.Unmarshal(body, &blockedRelationships); err != nil {
		return nil, err
	}

	return blockedRelationships, err
}

func getCampaigns(loggedUserID uint, followingProfiles []util.FollowingProfileDTO) ([]dto.SponsoredPostsDTO, error) {
	postBody, err := json.Marshal(followingProfiles)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	campaignHost, campaignPort := util.GetCampaignHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+campaignHost+":"+campaignPort+"/available-for-profile/"+util.Uint2String(loggedUserID),
		postBody, map[string]string{"Content-Type": "application/json"})

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var ret []dto.SponsoredPostsDTO
	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func getCampaignsWhereUserIsInfluencer(userID uint) ([]dto.SponsoredPostsDTO, error) {
	campaignHost, campaignPort := util.GetCampaignHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+campaignHost+":"+campaignPort+"/accepted-by-influencer/"+util.Uint2String(userID),
		nil, map[string]string{})

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("BAD_REQUEST")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var ret []dto.SponsoredPostsDTO
	if err = json.Unmarshal(body, &ret); err != nil {
		return nil, err
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
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodPost,
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

func getProfileIDByUsername(username string) (uint, error) {
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(context.Background(), http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-username/" + username,
		nil, map[string]string{})
	if err != nil {
		return 0, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	type ProfileDTO struct {
		ID uint `json:"id"`
	}
	var ret ProfileDTO
	if err = json.Unmarshal(body, &ret); err != nil {
		return 0, err
	}

	return ret.ID, nil
}
