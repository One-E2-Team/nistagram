package service

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/johnfercher/maroto/pkg/color"
	"github.com/johnfercher/maroto/pkg/consts"
	"github.com/johnfercher/maroto/pkg/pdf"
	"github.com/johnfercher/maroto/pkg/props"
	"io/ioutil"
	"net/http"
	"nistagram/agent/dto"
	"nistagram/agent/model"
	"nistagram/agent/util"
	"strings"
)

type CampaignService struct {
}

func (service *CampaignService) SaveCampaignReport(campaignId uint) error {
	resp, err := util.NistagramRequest(http.MethodGet, "/agent-api/statistics/"+util.Uint2String(campaignId),
		nil, map[string]string{})

	if err != nil {
		return err
	}

	var stat dto.StatisticsDTO

	err = json.NewDecoder(resp.Body).Decode(&stat)
	if err != nil {
		return err
	}

	var report model.CampaignReport
	var basicInfo model.BasicInformation
	var overallStat model.OverallStatistics
	var oStats model.Stats
	var paramStat []model.ParametersStatistics

	basicInfo.CampaignId = stat.CampaignId
	basicInfo.PostID = stat.Campaign.PostID
	basicInfo.AgentID = stat.Campaign.AgentID
	basicInfo.CampaignType = stat.Campaign.CampaignType
	basicInfo.Start = stat.Campaign.Start

	//overall stats
	for _, event := range stat.Events {
		switch strings.ToLower(event.EventType) {
		case "like":
			oStats.Likes += 1
		case "dislike":
			oStats.Dislikes += 1
		case "like_reset":
			oStats.LikeResets += 1
		case "dislike_reset":
			oStats.DislikeResets += 1
		case "comment":
			oStats.Comments += 1
		case "visit":
			oStats.TotalSiteVisits += 1
			oStats.AddSpecificSite(event.WebSite)
		}
	}

	//params stats
	for _, params := range stat.Campaign.CampaignParameters {
		var ps model.ParametersStatistics
		ps.Start = params.Start
		ps.End = params.End
		ps.Timestamps = params.Timestamps
		basicInfo.End = params.End

		var infNotAccepted []string
		for _, req := range params.CampaignRequests {
			if strings.ToLower(req.RequestStatus) == "declined" {
				infNotAccepted = append(infNotAccepted, req.InfluencerUsername)
			}
		}

		for _, event := range stat.Events {
			if event.Timestamp.After(params.Start) && event.Timestamp.Before(params.End) {
				if event.InfluencerId != 0 {
					ps.AddEventForInf(event.InfluencerUsername, event.EventType, event.WebSite)
				} else if len(event.Interests) != 0 {
					ps.AddEventForInterest(event.Interests, event.EventType, event.WebSite)
				} else {
					//TODO: direct event
				}
			}
		}

		ps.InfluencerWhoDidNotAccept = infNotAccepted
		paramStat = append(paramStat, ps)
	}

	overallStat.Stats = oStats
	report.BasicInformation = basicInfo
	report.OverallStatistics = overallStat
	report.ParametersStatistics = paramStat

	output, err := xml.MarshalIndent(&report, "  ", "    ")
	if err != nil {
		return err
	}

	campIdString := util.Uint2String(report.BasicInformation.CampaignId)
	resp, err = util.ExistDBRequest(http.MethodPut, "/exist/rest/collection/report" + campIdString + ".xml", output, map[string]string{})
	if err != nil {
		return err
	}
	fmt.Println(resp.StatusCode)
	fmt.Println("XML document successfully written!")

	return nil
}

func (service *CampaignService) GetMyCampaigns() (*http.Response, error) {
	return util.NistagramRequest(http.MethodGet, "/agent-api/campaign/my-campaigns",
		nil, map[string]string{})
}

func (service *CampaignService) CreateCampaign(requestBody []byte) (*http.Response, error) {
	return util.NistagramRequest(http.MethodPost, "/agent-api/campaign/create",
		requestBody, map[string]string{"Content-Type": "application/json"})
}

func (service *CampaignService) GetInterests() (*http.Response, error) {
	return util.NistagramRequest(http.MethodGet, "/agent-api/campaign/interests",
		nil, map[string]string{})
}

func (service *CampaignService) GetActiveParams(campaignID string) (*http.Response, error) {
	return util.NistagramRequest(http.MethodGet, "/agent-api/campaign/"+campaignID+"/params/active",
		nil, map[string]string{})
}

func (service *CampaignService) EditCampaign(postID string, requestBody []byte) (*http.Response, error) {
	return util.NistagramRequest(http.MethodPut, "/agent-api/campaign/update/"+postID,
		requestBody, map[string]string{"Content-Type": "application/json"})
}

func (service *CampaignService) GeneratePdfForCampaign(campaignId uint) error {
	resp, err := util.ExistDBRequest(http.MethodGet, "/exist/rest/collection/reports/report" + util.Uint2String(campaignId) + ".xml", []byte(""), map[string]string{})
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var report model.CampaignReport
	err = xml.Unmarshal(body, &report)
	if err != nil{
		return err
	}

	m := pdf.NewMaroto(consts.Portrait, consts.Letter)

	m.Row(20, func() {
		m.Col(4, func() {
			m.Text("CAMPAIGN REPORT - Overall statistics", props.Text{
				Top:         12,
				Size:        20,
				Extrapolate: true,
			})
		})
		m.ColSpace(4)
	})

	m.Line(5)

	m.Row(10, func() {
		m.Col(4, func() {
			m.Text("Basic campaign information:", props.Text{
				Size:        14,
				Top:         22,
				Extrapolate: true,
				Style: consts.BoldItalic,
			})
		})
		m.ColSpace(4)
	})

	m.Row(6, func() {
		m.Col(4, func() {
			m.Text(" - Campaign ID: " + util.Uint2String(report.BasicInformation.CampaignId), props.Text{
				Size: 12,
				Top:  22,
			})
		})
		m.ColSpace(4)
	})
	m.Row(6, func() {
		m.Col(4, func() {
			m.Text(" - Campaign type: " + report.BasicInformation.CampaignType, props.Text{
				Size: 12,
				Top:  22,

			})
		})
		m.ColSpace(4)
	})
	m.Row(6, func() {
		m.Col(4, func() {
			m.Text(" - Post ID: " + report.BasicInformation.PostID, props.Text{
				Size: 12,
				Top:  22,

			})
		})
		m.ColSpace(4)
	})
	m.Row(6, func() {
		m.Col(4, func() {
			m.Text(" - Agent ID: " + util.Uint2String(report.BasicInformation.AgentID), props.Text{
				Size: 12,
				Top:  22,

			})
		})
		m.ColSpace(4)
	})
	m.Row(6, func() {
		m.Col(4, func() {
			m.Text(" - Start date: " + report.BasicInformation.Start.Format("2006-01-02 15:04:05"), props.Text{
				Size: 12,
				Top:  22,

			})
		})
		m.ColSpace(4)
	})
	m.Row(6, func() {
		m.Col(4, func() {
			m.Text(" - End date: " + report.BasicInformation.End.Format("2006-01-02 15:04:05"), props.Text{
				Size: 12,
				Top:  22,

			})
		})
		m.ColSpace(4)
	})

	m.Row(10, func() {
		m.Col(4, func() {
			m.Text("Overall campaign statistics:", props.Text{
				Size:        14,
				Top:         22,
				Extrapolate: true,
				Style: consts.BoldItalic,
			})
		})
		m.ColSpace(4)
	})

	header := []string{"Likes", "Dislikes", "Like resets", "Dislike resets", "Comments", "Site visits"}
	contents := [][]string{
		{util.Uint2String(report.OverallStatistics.Stats.Likes),
			util.Uint2String(report.OverallStatistics.Stats.Dislikes),
			util.Uint2String(report.OverallStatistics.Stats.LikeResets),
			util.Uint2String(report.OverallStatistics.Stats.DislikeResets),
			util.Uint2String(report.OverallStatistics.Stats.Comments),
			util.Uint2String(report.OverallStatistics.Stats.TotalSiteVisits),
		},
	}

	m.Row(10, func(){})

	m.Row(30, func() {
		m.TableList(header, contents, props.TableList{
			HeaderProp: props.TableListContent{
				Size:      13,
			},
			ContentProp: props.TableListContent{
				Size:      13,
			},
			Align:                consts.Center,
			AlternatedBackground: &color.Color{Red: 150, Green: 150, Blue: 150},
			HeaderContentSpace:   1,
			Line:                 false,
		})
	})

	m.Row(10, func(){})

	header = []string{"Website", "Visits"}
	contents = [][]string{
	}

	for _, item := range report.OverallStatistics.Stats.SpecificSiteVisits{
		var cont []string
		cont = append(cont, item.WebSite)
		cont = append(cont, util.Uint2String(item.Visits))
		contents = append(contents, cont)
	}

	m.Row(10, func(){})

	m.Row(50, func() {
		m.TableList(header, contents, props.TableList{
			HeaderProp: props.TableListContent{
				Size:      13,
			},
			ContentProp: props.TableListContent{
				Size:      13,
			},
			Align:                consts.Center,
			AlternatedBackground: &color.Color{Red: 150, Green: 150, Blue: 150},
			HeaderContentSpace:   1,
			Line:                 false,
		})
	})

	m.Row(15, func() {
		m.Col(4, func() {
			m.Text("Parameters statistics:", props.Text{
				Size:        14,
				Top:         22,
				Extrapolate: true,
				Style: consts.BoldItalic,
			})
		})
		m.ColSpace(4)
	})

	for _, params := range report.ParametersStatistics {

		m.Row(6, func() {
			m.Col(4, func() {
				m.Text(" - Period: " + params.Start.Format("2006-01-02 15:04:05") +
				" - " + params.End.Format("2006-01-02 15:04:05"), props.Text{
					Size: 12,
					Top:  22,
				})
			})
			m.ColSpace(4)
		})

		timestamps := ""
		for _, ts := range params.Timestamps{
			timestamps += " " + ts.Format("8:00") + ","
		}

		m.Row(6, func() {
			m.Col(12, func() {
				m.Text(" - Timestamps when the campaign has been launched: " + timestamps, props.Text{
					Size: 12,
					Top:  22,
				})
			})
			m.ColSpace(4)
		})

		m.Row(15, func() {
			m.Col(4, func() {
				m.Text("- Influencers:", props.Text{
					Size:        14,
					Top:         22,
					Extrapolate: true,
					Style:       consts.Bold,
				})
			})
			m.ColSpace(4)
		})

		m.Row(10, func() {})

		header = []string{"Username", "Likes", "Dislikes", "Like resets", "Dislike resets", "Comments", "Site visits"}
		contents = [][]string{
		}

		for _, item := range params.InfluencerStats{
			var cont []string
			cont = append(cont, item.Username)
			cont = append(cont, util.Uint2String(item.Stats.Likes))
			cont = append(cont, util.Uint2String(item.Stats.Dislikes))
			cont = append(cont, util.Uint2String(item.Stats.LikeResets))
			cont = append(cont, util.Uint2String(item.Stats.DislikeResets))
			cont = append(cont, util.Uint2String(item.Stats.Comments))
			cont = append(cont, util.Uint2String(item.Stats.TotalSiteVisits))
			contents = append(contents, cont)
		}

		m.Row(50, func() {
			m.TableList(header, contents, props.TableList{
				HeaderProp: props.TableListContent{
					Size: 12,
				},
				ContentProp: props.TableListContent{
					Size: 12,
				},
				Align:                consts.Center,
				AlternatedBackground: &color.Color{Red: 150, Green: 150, Blue: 150},
				HeaderContentSpace:   1,
				Line:                 false,
			})
		})

		infNotAccepted := ""
		for _, inf := range params.InfluencerWhoDidNotAccept{
			infNotAccepted += " " + inf + ","
		}

		m.Row(15, func() {
			m.Col(4, func() {
				m.Text("- Influencers who did not accept a campaign request: " + infNotAccepted, props.Text{
					Size:        14,
					Top:         22,
					Extrapolate: true,
					Style:       consts.Bold,
				})
			})
			m.ColSpace(4)
		})

		/*
			m.Row(10, func(){})

			m.Row(15, func() {
				m.Col(4, func() {
					m.Text("- Direct interation with agent profile: ", props.Text{
						Size:        14,
						Top:         22,
						Extrapolate: true,
						Style: consts.Bold,
					})
				})
				m.ColSpace(4)
			})*/

		m.Row(10, func() {})

		m.Row(15, func() {
			m.Col(4, func() {
				m.Text("- Target groups:", props.Text{
					Size:        14,
					Top:         22,
					Extrapolate: true,
					Style:       consts.Bold,
				})
			})
			m.ColSpace(4)
		})

		m.Row(10, func() {})

		header = []string{"Interest", "Likes", "Dislikes", "Like resets", "Dislike resets", "Comments", "Site visits"}
		contents = [][]string{
		}

		for _, item := range params.TargetGroupsStats{
			var cont []string
			cont = append(cont, item.Interest)
			cont = append(cont, util.Uint2String(item.Stats.Likes))
			cont = append(cont, util.Uint2String(item.Stats.Dislikes))
			cont = append(cont, util.Uint2String(item.Stats.LikeResets))
			cont = append(cont, util.Uint2String(item.Stats.DislikeResets))
			cont = append(cont, util.Uint2String(item.Stats.Comments))
			cont = append(cont, util.Uint2String(item.Stats.TotalSiteVisits))
			contents = append(contents, cont)
		}

		m.Row(50, func() {
			m.TableList(header, contents, props.TableList{
				HeaderProp: props.TableListContent{
					Size: 12,
				},
				ContentProp: props.TableListContent{
					Size: 12,
				},
				Align:                consts.Center,
				AlternatedBackground: &color.Color{Red: 150, Green: 150, Blue: 150},
				HeaderContentSpace:   1,
				Line:                 false,
			})
		})
	}

	err = m.OutputFileAndClose("../../nistagramstaticdata/data/report" +
		util.Uint2String(campaignId) + ".pdf")
	if err != nil {
		return err
	}

	return nil
}


func (service *CampaignService) GeneratePdfForSortedCampaigns() error {
	//var query = map[string]string{"_query": "collection(\"/\")/campaign_report"}
	var query = map[string]string{"_query": "//basic_information | //overall_statistics"}
	resp, err := util.ExistDBRequest(http.MethodGet, "/exist/rest/collection"+util.GenerateFuckingExistDBHTTPRequestParametersQuery(query), []byte(""), map[string]string{})
	if err != nil {
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var reports []model.CampaignReport
	err = xml.Unmarshal(body, &reports)
	if err != nil{
		return err
	}
	/*
	expr, err := xpath.Compile("//basic_information | //overall_statistics")
	if err != nil {
		fmt.Println(err)
	}
	node, err := xmlquery.Parse(bytes.NewReader(body))
	if err != nil{
		return err
	}

	nodeIter := expr.Evaluate(xmlquery.CreateXPathNavigator(node)).(*xpath.NodeIterator)

	var reports []model.CampaignReport

	for nodeIter.MoveNext() == true{

		if nodeIter.Current().LocalName() == "campaign_id"{
			var report model.CampaignReport
			report.BasicInformation.CampaignId = util.String2Uint(nodeIter.Current().Value())
			reports = append(reports, report)
		}

	}*/

	reports = sortReports(reports)

	m := pdf.NewMaroto(consts.Portrait, consts.Letter)

	m.Row(20, func() {
		m.Col(4, func() {
			m.Text("CAMPAIGN REPORT - Campaign comparison", props.Text{
				Top:         12,
				Size:        20,
				Extrapolate: true,
			})
		})
		m.ColSpace(4)
	})

	m.Line(5)

	for _, report := range reports{
		m.Row(10, func() {
			m.Col(4, func() {
				m.Text("Basic campaign information:", props.Text{
					Size:        14,
					Top:         22,
					Extrapolate: true,
					Style: consts.BoldItalic,
				})
			})
			m.ColSpace(4)
		})

		m.Row(6, func() {
			m.Col(4, func() {
				m.Text(" - Campaign ID: " + util.Uint2String(report.BasicInformation.CampaignId), props.Text{
					Size: 12,
					Top:  22,
				})
			})
			m.ColSpace(4)
		})
		m.Row(6, func() {
			m.Col(4, func() {
				m.Text(" - Campaign type: " + report.BasicInformation.CampaignType, props.Text{
					Size: 12,
					Top:  22,

				})
			})
			m.ColSpace(4)
		})
		m.Row(6, func() {
			m.Col(4, func() {
				m.Text(" - Post ID: " + report.BasicInformation.PostID, props.Text{
					Size: 12,
					Top:  22,

				})
			})
			m.ColSpace(4)
		})
		m.Row(6, func() {
			m.Col(4, func() {
				m.Text(" - Agent ID: " + util.Uint2String(report.BasicInformation.AgentID), props.Text{
					Size: 12,
					Top:  22,

				})
			})
			m.ColSpace(4)
		})
		m.Row(6, func() {
			m.Col(4, func() {
				m.Text(" - Start date: " + report.BasicInformation.Start.Format("2006-01-02 15:04:05"), props.Text{
					Size: 12,
					Top:  22,

				})
			})
			m.ColSpace(4)
		})
		m.Row(6, func() {
			m.Col(4, func() {
				m.Text(" - End date: " + report.BasicInformation.End.Format("2006-01-02 15:04:05"), props.Text{
					Size: 12,
					Top:  22,

				})
			})
			m.ColSpace(4)
		})

		m.Row(10, func() {
			m.Col(4, func() {
				m.Text("Overall campaign statistics:", props.Text{
					Size:        14,
					Top:         22,
					Extrapolate: true,
					Style: consts.BoldItalic,
				})
			})
			m.ColSpace(4)
		})

		header := []string{"Likes", "Dislikes", "Like resets", "Dislike resets", "Comments", "Site visits"}
		contents := [][]string{
			{util.Uint2String(report.OverallStatistics.Stats.Likes),
				util.Uint2String(report.OverallStatistics.Stats.Dislikes),
				util.Uint2String(report.OverallStatistics.Stats.LikeResets),
				util.Uint2String(report.OverallStatistics.Stats.DislikeResets),
				util.Uint2String(report.OverallStatistics.Stats.Comments),
				util.Uint2String(report.OverallStatistics.Stats.TotalSiteVisits),
			},
		}

		m.Row(10, func(){})

		m.Row(30, func() {
			m.TableList(header, contents, props.TableList{
				HeaderProp: props.TableListContent{
					Size:      13,
				},
				ContentProp: props.TableListContent{
					Size:      13,
				},
				Align:                consts.Center,
				AlternatedBackground: &color.Color{Red: 150, Green: 150, Blue: 150},
				HeaderContentSpace:   1,
				Line:                 false,
			})
		})

		m.Row(10, func(){})

		header = []string{"Website", "Visits"}
		contents = [][]string{
		}

		for _, item := range report.OverallStatistics.Stats.SpecificSiteVisits{
			var cont []string
			cont = append(cont, item.WebSite)
			cont = append(cont, util.Uint2String(item.Visits))
			contents = append(contents, cont)
		}

		m.Row(10, func(){})

		m.Row(50, func() {
			m.TableList(header, contents, props.TableList{
				HeaderProp: props.TableListContent{
					Size:      13,
				},
				ContentProp: props.TableListContent{
					Size:      13,
				},
				Align:                consts.Center,
				AlternatedBackground: &color.Color{Red: 150, Green: 150, Blue: 150},
				HeaderContentSpace:   1,
				Line:                 false,
			})
		})

		m.Row(20, func(){})
	}

	err = m.OutputFileAndClose("../../nistagramstaticdata/data/reportComparison.pdf")
	if err != nil {
		return err
	}

	return nil
}

func sortReports(reports []model.CampaignReport) []model.CampaignReport{
	var ret []model.CampaignReport

	for i := 0; i < len(reports); i++{
		ret[i] = reports[i]
		istart := ret[i].BasicInformation.Start
		iend := ret[i].BasicInformation.Start
		idays := iend.Sub(istart).Hours()/24
		istats := ret[i].OverallStatistics.Stats
		ipoints := float64(istats.Likes + istats.DislikeResets + istats.Comments + istats.TotalSiteVisits -
			istats.Dislikes - istats.LikeResets) / idays
		for j := i + 1; j < len(reports); j++{
			jstart := reports[j].BasicInformation.Start
			jend := reports[j].BasicInformation.End
			jdays := jend.Sub(jstart).Hours()/24
			jstats := reports[j].OverallStatistics.Stats
			jpoints := float64(jstats.Likes + jstats.DislikeResets + jstats.Comments + jstats.TotalSiteVisits -
				jstats.Dislikes - jstats.LikeResets) / jdays
			if jpoints > ipoints{
				t := ret[i]
				ret[i] = reports[j]
				reports[j] = t
			}
		}
	}

	return ret
}