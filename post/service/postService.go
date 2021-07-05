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

func (service *PostService) GetPublic(ctx context.Context, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPublic-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	blockedRelationships, err := getProfilesBlockedRelationships(nextCtx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	posts, err := service.PostRepository.GetPublic(nextCtx, blockedRelationships)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(nextCtx, posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetProfilesPosts(ctx context.Context, followingProfiles []util.FollowingProfileDTO, targetUsername string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetProfilesPosts-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	profileID, err := getProfileIDByUsername(nextCtx, targetUsername)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	sponsoredPostsDTO, err := getCampaignsWhereUserIsInfluencer(nextCtx, profileID)
	campaignPosts := make([]model.Post, 0)
	influencerIDs := make([]uint, 0)
	for _, sponsoredPostDTO := range sponsoredPostsDTO {
		primitiveID, err := primitive.ObjectIDFromHex(sponsoredPostDTO.PostID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		post, err := service.PostRepository.Read(nextCtx, primitiveID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		influencerIDs = append(influencerIDs, sponsoredPostDTO.InfluencerID)
		campaignPosts = append(campaignPosts, post)
	}
	initialSponsoredPosts, err := getReactionsForPosts(nextCtx, campaignPosts, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	ret := make([]dto.ResponsePostDTO, 0)
	influencerUsernames, err := getProfileUsernamesByIDs(nextCtx, influencerIDs)
	for i, initial := range initialSponsoredPosts {
		ret = append(ret, dto.ResponsePostDTO{
			Post:               initial.Post,
			Reaction:           initial.Reaction,
			CampaignId:         sponsoredPostsDTO[i].CampaignID,
			InfluencerId:       sponsoredPostsDTO[i].InfluencerID,
			InfluencerUsername: influencerUsernames[i],
		})
	}
	posts, err := service.PostRepository.GetProfilesPosts(nextCtx, followingProfiles, targetUsername)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	postsDTO, err := getReactionsForPosts(nextCtx, posts, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	for _, post := range postsDTO {
		ret = append(ret, post)
	}
	return ret, err
}

func (service *PostService) GetPublicPostByLocation(ctx context.Context, location string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPublicPostByLocation-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	blockedRelationships, err := getProfilesBlockedRelationships(nextCtx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	posts, err := service.PostRepository.GetPublicPostByLocation(nextCtx, location, blockedRelationships)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(nextCtx, posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetPublicPostByHashTag(ctx context.Context, hashTag string, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPublicPostByHashTag-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	blockedRelationships, err := getProfilesBlockedRelationships(nextCtx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	posts, err := service.PostRepository.GetPublicPostByHashTag(nextCtx, hashTag, blockedRelationships)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	responseDTO, err := getReactionsForPosts(nextCtx, posts, loggedUserID)
	return responseDTO, err
}

func (service *PostService) GetMyPosts(ctx context.Context, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMyPosts-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	sponsoredPostsDTO, err := getCampaignsWhereUserIsInfluencer(nextCtx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
	}
	campaignPosts := make([]model.Post, 0)
	influencerIDs := make([]uint, 0)
	for _, sponsoredPostDTO := range sponsoredPostsDTO {
		primitiveID, err := primitive.ObjectIDFromHex(sponsoredPostDTO.PostID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		post, err := service.PostRepository.Read(nextCtx, primitiveID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		influencerIDs = append(influencerIDs, sponsoredPostDTO.InfluencerID)
		campaignPosts = append(campaignPosts, post)
	}
	initialSponsoredPosts, err := getReactionsForPosts(nextCtx, campaignPosts, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	ret := make([]dto.ResponsePostDTO, 0)
	influencerUsernames, err := getProfileUsernamesByIDs(nextCtx, influencerIDs)
	if err != nil {
		util.Tracer.LogError(span, err)
	}
	for i, initial := range initialSponsoredPosts {
		ret = append(ret, dto.ResponsePostDTO{
			Post:               initial.Post,
			Reaction:           initial.Reaction,
			CampaignId:         sponsoredPostsDTO[i].CampaignID,
			InfluencerId:       sponsoredPostsDTO[i].InfluencerID,
			InfluencerUsername: influencerUsernames[i],
		})
	}
	posts, err := service.PostRepository.GetMyPosts(nextCtx, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	postsDTO, err := getReactionsForPosts(nextCtx, posts, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	for _, post := range postsDTO {
		ret = append(ret, post)
	}
	return ret, err
}

func (service *PostService) GetPostsForHomePage(ctx context.Context, followingProfiles []util.FollowingProfileDTO, loggedUserID uint) ([]dto.ResponsePostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetPostsForHomePage-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	posts, err := service.PostRepository.GetPostsForHomePage(nextCtx, followingProfiles)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	ret := make([]dto.ResponsePostDTO, 0)
	ret, err = getReactionsForPosts(nextCtx, posts, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	sponsoredPostsDTO, err := getCampaigns(nextCtx, loggedUserID, followingProfiles)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	campaignPosts := make([]model.Post, 0)
	influencerIDs := make([]uint, 0)
	for _, sponsoredPostDTO := range sponsoredPostsDTO {
		primitiveID, err := primitive.ObjectIDFromHex(sponsoredPostDTO.PostID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		post, err := service.PostRepository.Read(nextCtx, primitiveID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return nil, err
		}
		influencerIDs = append(influencerIDs, sponsoredPostDTO.InfluencerID)
		campaignPosts = append(campaignPosts, post)
	}
	initialSponsoredPosts, err := getReactionsForPosts(nextCtx, campaignPosts, loggedUserID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	influencerUsernames, err := getProfileUsernamesByIDs(nextCtx, influencerIDs)
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

func (service *PostService) CreatePost(ctx context.Context, postType model.PostType, post dto.PostDto, mediaNames []string, profile dto.ProfileDto) error {
	span := util.Tracer.StartSpanFromContext(ctx, "CreatePost-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

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
		err := canUsersBeTagged(nextCtx, post.Description, profile.ProfileId)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
	}

	newPost := model.Post{ID: primitive.NewObjectID(), PublisherId: profile.ProfileId, PublisherUsername: profile.Username,
		PostType: postType, Medias: medias, PublishDate: time.Now(),
		Description: post.Description, IsHighlighted: post.IsHighlighted, IsCampaign: false,
		IsCloseFriendsOnly: post.IsCloseFriendsOnly, Location: post.Location,
		HashTags: post.HashTags, IsPrivate: profile.ProfileSettings.IsPrivate, IsDeleted: false}

	return service.PostRepository.Create(nextCtx, &newPost)
}

func (service *PostService) ReadPost(ctx context.Context, id primitive.ObjectID) (model.Post, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "ReadPost-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.PostRepository.Read(nextCtx, id)
}

func (service *PostService) DeletePost(ctx context.Context, id primitive.ObjectID) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeletePost-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	err := service.PostRepository.Delete(nextCtx, id)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	err = deletePostsReports(nextCtx, id)
	return err
}

func (service *PostService) UpdatePost(ctx context.Context, id primitive.ObjectID, post dto.PostDto) error {
	span := util.Tracer.StartSpanFromContext(ctx, "UpdatePost-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.PostRepository.Update(nextCtx, id, post)
}

func (service *PostService) DeleteUserPosts(ctx context.Context, profileId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "DeleteUserPosts-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	posts, err := service.PostRepository.GetMyPosts(nextCtx, profileId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	for i := 0; i < len(posts); i++ {
		err = service.DeletePost(nextCtx, posts[i].ID)
		if err != nil {
			util.Tracer.LogError(span, err)
			return err
		}
	}
	return nil
}

func (service *PostService) ChangeUsername(ctx context.Context, profileId uint, username string) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ChangeUsername-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	return service.PostRepository.ChangeUsername(nextCtx, profileId, username)
}

func (service *PostService) ChangePrivacy(ctx context.Context, profileId uint, isPrivate bool) error {
	span := util.Tracer.StartSpanFromContext(ctx, "ChangePrivacy-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	return service.PostRepository.ChangePrivacy(nextCtx, profileId, isPrivate)
}

func (service *PostService) MakeCampaign(ctx context.Context, postID string, agentID uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "MakeCampaign-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	id, err := primitive.ObjectIDFromHex(postID)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	post, err := service.PostRepository.Read(nextCtx, id)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	if post.PublisherId != agentID {
		util.Tracer.LogError(span, fmt.Errorf("bad agent id"))
		return fmt.Errorf("BAD_AGENT_ID")
	}
	return service.PostRepository.MakeCampaign(nextCtx, id)
}

func (service *PostService) GetMediaById(ctx context.Context, mediaId string) (model.Media, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "GetMediaById-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	return service.PostRepository.GetMediaById(nextCtx, mediaId)
}

func deletePostsReports(ctx context.Context, postId primitive.ObjectID) error {
	span := util.Tracer.StartSpanFromContext(ctx, "deletePostsReports-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	postReactionHost, postReactionPort := util.GetPostReactionHostAndPort()
	_, err := util.CrossServiceRequest(nextCtx, http.MethodDelete,
		util.GetCrossServiceProtocol()+"://"+postReactionHost+":"+postReactionPort+"/report/"+util.GetStringIDFromMongoID(postId),
		nil, map[string]string{})
	return err
}

func canUsersBeTagged(ctx context.Context, description string, publisherId uint) error {
	span := util.Tracer.StartSpanFromContext(ctx, "canUsersBeTagged-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	var followingProfiles []uint

	resp, err := getUserFollowers(nextCtx, publisherId)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if err = json.Unmarshal(body, &followingProfiles); err != nil {
		util.Tracer.LogError(span, err)
		return err
	}

	descriptionParts := strings.Split(description, " ")
	for i := 0; i < len(descriptionParts); i++ {
		if strings.Contains(descriptionParts[i], "@") {
			taggedUsername := descriptionParts[i][1:len(descriptionParts[i])]
			var taggedProfile dto.ProfileDto
			if resp, err := getProfileByUsername(nextCtx, taggedUsername); err != nil {
				util.Tracer.LogError(span, err)
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
				util.Tracer.LogError(span, fmt.Errorf("%s can't be tagged",taggedProfile.Username))
				return errors.New(taggedProfile.Username + " can't be tagged!")
			}

			if !util.Contains(followingProfiles, taggedProfile.ProfileId) {
				util.Tracer.LogError(span, fmt.Errorf("%s is not followed by this profile",taggedProfile.Username))
				return errors.New(taggedProfile.Username + " is not followed by this profile!")
			}
		}
	}
	return nil
}

func getUserFollowers(ctx context.Context, loggedUserId uint) (*http.Response, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getUserFollowers-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	connHost, connPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+connHost+":"+connPort+"/connection/following/show/"+util.Uint2String(loggedUserId),
		nil, map[string]string{})
	return resp, err
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

func getReactionsForPosts(ctx context.Context, posts []model.Post, profileID uint) ([]dto.ResponsePostDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getReactionsForPosts-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

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
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+postReactionHost+":"+postReactionPort+"/get-reaction-types/"+util.Uint2String(profileID),
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

	var reactions []string
	if err = json.Unmarshal(body, &reactions); err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	if len(posts) != len(reactions) {
		util.Tracer.LogError(span, fmt.Errorf("bad lists"))
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

func getProfilesBlockedRelationships(ctx context.Context, loggedProfileId uint) ([]uint, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getProfilesBlockedRelationships-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	connectionHost, connectionPort := util.GetConnectionHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+connectionHost+":"+connectionPort+"/connection/block/relationships/"+util.Uint2String(loggedProfileId),
		nil, map[string]string{})

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

	var blockedRelationships []uint
	if err = json.Unmarshal(body, &blockedRelationships); err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}

	return blockedRelationships, err
}

func getCampaigns(ctx context.Context, loggedUserID uint, followingProfiles []util.FollowingProfileDTO) ([]dto.SponsoredPostsDTO, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getCampaigns-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	postBody, err := json.Marshal(followingProfiles)
	if err != nil {
		util.Tracer.LogError(span, err)
		fmt.Println(err)
		return nil, err
	}
	campaignHost, campaignPort := util.GetCampaignHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodPost,
		util.GetCrossServiceProtocol()+"://"+campaignHost+":"+campaignPort+"/available-for-profile/"+util.Uint2String(loggedUserID),
		postBody, map[string]string{"Content-Type": "application/json"})

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		util.Tracer.LogError(span, fmt.Errorf("bad request"))
		return nil, fmt.Errorf("BAD_REQUEST")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var ret []dto.SponsoredPostsDTO
	if err = json.Unmarshal(body, &ret); err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	return ret, nil
}

func getCampaignsWhereUserIsInfluencer(ctx context.Context, userID uint) ([]dto.SponsoredPostsDTO, error){
	span := util.Tracer.StartSpanFromContext(ctx, "getCampaignsWhereUserIsInfluencer-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)

	campaignHost, campaignPort := util.GetCampaignHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+campaignHost+":"+campaignPort+"/accepted-by-influencer/"+util.Uint2String(userID),
		nil, map[string]string{})

	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	if resp.StatusCode != 200 {
		util.Tracer.LogError(span, fmt.Errorf("bad request"))
		return nil, fmt.Errorf("BAD_REQUEST")
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)
	var ret []dto.SponsoredPostsDTO
	if err = json.Unmarshal(body, &ret); err != nil {
		util.Tracer.LogError(span, err)
		return nil, err
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

func getProfileIDByUsername(ctx context.Context, username string) (uint, error) {
	span := util.Tracer.StartSpanFromContext(ctx, "getProfileIDByUsername-service")
	defer util.Tracer.FinishSpan(span)

	nextCtx := util.Tracer.ContextWithSpan(ctx, span)
	profileHost, profilePort := util.GetProfileHostAndPort()
	resp, err := util.CrossServiceRequest(nextCtx, http.MethodGet,
		util.GetCrossServiceProtocol()+"://"+profileHost+":"+profilePort+"/get-by-username/" + username,
		nil, map[string]string{})
	if err != nil {
		util.Tracer.LogError(span, err)
		return 0, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		util.Tracer.LogError(span, err)
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
		util.Tracer.LogError(span, err)
		return 0, err
	}

	return ret.ID, nil
}
