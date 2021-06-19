package dto

type ShowReportDTO struct {
	ReportID 		   string	 `json:"reportId"`
	PostID 		       string 	 `json:"postId"`
	Reason 		       string 	 `json:"reason"`
	PublisherId        uint      `json:"publisherId"`
	PublisherUsername  string    `json:"publisherUsername"`
	Medias             []Media   `json:"medias"`
	Description        string    `json:"description"`
}

type Media struct {
	FilePath string `json:"filePath"`
}
