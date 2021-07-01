package model

import (
	"encoding/xml"
	"strings"
	"time"
)

type CampaignReport struct{
	XMLName        			xml.Name       				`json:"campaignReport" xml:"campaign_report"`
	BasicInformation 		BasicInformation 			`json:"basicInformation" xml:"basic_information"`
	OverallStatistics		OverallStatistics			`json:"overallStatistics" xml:"overall_statistics"`
	ParametersStatistics	[]ParametersStatistics		`json:"parametersStatistics" xml:"parameters_statistics"`
}

type BasicInformation struct {
	CampaignId		   uint				        `json:"campaignId" xml:"campaign_id"`
	PostID             string               	`json:"postId" xml:"post_id"`
	AgentID            uint                 	`json:"agentId" xml:"agent_id"`
	CampaignType       string               	`json:"campaignType" xml:"campaign_type"`
	Start              time.Time            	`json:"start" xml:"start"`
	End                time.Time            	`json:"end" xml:"end"`
}

type OverallStatistics struct {
	Stats 				Stats 			`json:"stats" xml:"stats"`
}

type ParametersStatistics struct {
	Start            			 time.Time            `json:"start" xml:"start"`
	End              			 time.Time            `json:"end" xml:"end"`
	Timestamps  	 			 []time.Time	      `json:"timestamps" xml:"timestamps"`
	InfluencerStats	 			 []InfluencerStats	  `json:"influencerStats" xml:"influencer_stats"`
	InfluencerWhoDidNotAccept	 []string	  		  `json:"influencerWhoDidNotAccept" xml:"influencer_who_did_not_accept"`
	TargetGroupsStats	 		 []TargetGroupsStats  `json:"targetGroupsStats" xml:"target_groups_stats"`
}

type InfluencerStats struct {
	Username			string		    `json:"username" xml:"username"`
	Stats 				Stats 			`json:"stats" xml:"stats"`
}

type TargetGroupsStats struct {
	Interest			string		    `json:"interest" xml:"interest"`
	Stats 				Stats 			`json:"stats" xml:"stats"`
}

type Stats struct{
	Likes				uint 			`json:"likes" xml:"likes"`
	Dislikes 			uint			`json:"dislikes" xml:"dislikes"`
	LikeResets			uint			`json:"likeResets" xml:"likeResets"`
	DislikeResets		uint			`json:"dislikeResets" xml:"dislike_resets"`
	Comments			uint			`json:"comments" xml:"comments"`
	TotalSiteVisits		uint			`json:"totalSiteVisits" xml:"total_site_visits"`
	SpecificSiteVisits	[]SiteVisit		`json:"specificSiteVisits" xml:"specific_site_visits"`
}

func (stats *Stats) AddSpecificSite(site string){
	for _, sv := range stats.SpecificSiteVisits{
		if sv.WebSite == site{
			sv.Visits += 1
			return
		}
	}
	sv := SiteVisit{WebSite: site, Visits: 1}
	stats.SpecificSiteVisits = append(stats.SpecificSiteVisits, sv)
}

type SiteVisit struct{
	WebSite 		string		`json:"webSite"`
	Visits			uint		`json:"visits"`
}

func (ps *ParametersStatistics) AddEventForInf(username string,eventType string, website string){
	for _, infStat := range ps.InfluencerStats{
		if username == infStat.Username{
			switch strings.ToLower(eventType){
			case "like":
				infStat.Stats.Likes += 1
			case "dislike":
				infStat.Stats.Dislikes += 1
			case "like_reset":
				infStat.Stats.LikeResets += 1
			case "dislike_reset":
				infStat.Stats.DislikeResets += 1
			case "comment":
				infStat.Stats.Comments += 1
			case "visit":
				infStat.Stats.TotalSiteVisits += 1
				infStat.Stats.AddSpecificSite(website)
			}
			return
		}
	}
	stats := Stats{}
	switch strings.ToLower(eventType){
	case "like":
		stats.Likes += 1
	case "dislike":
		stats.Dislikes += 1
	case "like_reset":
		stats.LikeResets += 1
	case "dislike_reset":
		stats.DislikeResets += 1
	case "comment":
		stats.Comments += 1
	case "visit":
		stats.TotalSiteVisits += 1
		stats.AddSpecificSite(website)
	}
	infStats := InfluencerStats{Username: username, Stats: stats}
	ps.InfluencerStats = append(ps.InfluencerStats, infStats)
}

func (ps *ParametersStatistics) AddEventForInterest(interests []string, eventType string, website string){
	for _, inter := range interests{
		for _, tgStat := range ps.TargetGroupsStats{
			if inter == tgStat.Interest{
				switch strings.ToLower(eventType){
				case "like":
					tgStat.Stats.Likes += 1
				case "dislike":
					tgStat.Stats.Dislikes += 1
				case "like_reset":
					tgStat.Stats.LikeResets += 1
				case "dislike_reset":
					tgStat.Stats.DislikeResets += 1
				case "comment":
					tgStat.Stats.Comments += 1
				case "visit":
					tgStat.Stats.TotalSiteVisits += 1
					tgStat.Stats.AddSpecificSite(website)
				}
				return
			}
		}
		stats := Stats{}
		switch strings.ToLower(eventType){
		case "like":
			stats.Likes += 1
		case "dislike":
			stats.Dislikes += 1
		case "like_reset":
			stats.LikeResets += 1
		case "dislike_reset":
			stats.DislikeResets += 1
		case "comment":
			stats.Comments += 1
		case "visit":
			stats.TotalSiteVisits += 1
			stats.AddSpecificSite(website)
		}
		tgStats := TargetGroupsStats{Interest: inter, Stats: stats}
		ps.TargetGroupsStats = append(ps.TargetGroupsStats, tgStats)
	}
}