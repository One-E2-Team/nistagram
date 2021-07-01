package model

import (
	"strings"
	"time"
)

type CampaignReport struct{
	BasicInformation 		BasicInformation 			`json:"basicInformation"`
	OverallStatistics		OverallStatistics			`json:"overallStatistics"`
	ParametersStatistics	[]ParametersStatistics		`json:"parametersStatistics"`
}

type BasicInformation struct {
	CampaignId		   uint				        `json:"campaignId"`
	PostID             string               	`json:"postId"`
	AgentID            uint                 	`json:"agentId"`
	CampaignType       string               	`json:"campaignType"`
	Start              time.Time            	`json:"start"`
	End                time.Time            	`json:"end"`
}

type OverallStatistics struct {
	Stats 				Stats 			`json:"stats"`
}

type ParametersStatistics struct {
	Start            			 time.Time            `json:"start"`
	End              			 time.Time            `json:"end"`
	Timestamps  	 			 []time.Time	      `json:"timestamps"`
	InfluencerStats	 			 []InfluencerStats	  `json:"influencerStats"`
	InfluencerWhoDidNotAccept	 []string	  		  `json:"influencerWhoDidNotAccept"`
	TargetGroupsStats	 		 []TargetGroupsStats  `json:"targetGroupsStats"`
}

type InfluencerStats struct {
	Username			string		    `json:"username"`
	Stats 				Stats 			`json:"stats"`
}

type TargetGroupsStats struct {
	Interest			string		    `json:"interest"`
	Stats 				Stats 			`json:"stats"`
}

type Stats struct{
	Likes				uint 			`json:"likes"`
	Dislikes 			uint			`json:"dislikes"`
	LikeResets			uint			`json:"likeResets"`
	DislikeResets		uint			`json:"dislikeResets"`
	Comments			uint			`json:"comments"`
	TotalSiteVisits		uint			`json:"totalSiteVisits"`
	SpecificSiteVisits	[]SiteVisit		`json:"specificSiteVisits"`
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