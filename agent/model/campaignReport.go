package model

import (
	"encoding/xml"
	"strings"
	"time"
)

type MonitoringResult struct {
	Results []CampaignReport      				`json:"result" xml:"result"`
}

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
	for i, _ := range stats.SpecificSiteVisits{
		if stats.SpecificSiteVisits[i].WebSite == site{
			stats.SpecificSiteVisits[i].Visits += 1
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
	for i, _ := range ps.InfluencerStats{
		if username == ps.InfluencerStats[i].Username{
			switch strings.ToLower(eventType){
			case "like":
				ps.InfluencerStats[i].Stats.Likes += 1
			case "dislike":
				ps.InfluencerStats[i].Stats.Dislikes += 1
			case "like_reset":
				ps.InfluencerStats[i].Stats.LikeResets += 1
			case "dislike_reset":
				ps.InfluencerStats[i].Stats.DislikeResets += 1
			case "comment":
				ps.InfluencerStats[i].Stats.Comments += 1
			case "visit":
				ps.InfluencerStats[i].Stats.TotalSiteVisits += 1
				ps.InfluencerStats[i].Stats.AddSpecificSite(website)
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
		for i, _ := range ps.TargetGroupsStats{
			if inter == ps.TargetGroupsStats[i].Interest{
				switch strings.ToLower(eventType){
				case "like":
					ps.TargetGroupsStats[i].Stats.Likes += 1
				case "dislike":
					ps.TargetGroupsStats[i].Stats.Dislikes += 1
				case "like_reset":
					ps.TargetGroupsStats[i].Stats.LikeResets += 1
				case "dislike_reset":
					ps.TargetGroupsStats[i].Stats.DislikeResets += 1
				case "comment":
					ps.TargetGroupsStats[i].Stats.Comments += 1
				case "visit":
					ps.TargetGroupsStats[i].Stats.TotalSiteVisits += 1
					ps.TargetGroupsStats[i].Stats.AddSpecificSite(website)
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