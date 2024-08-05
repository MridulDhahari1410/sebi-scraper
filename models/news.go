package models

// S3NewsMoreInfoDetails is the more info details.
type S3NewsMoreInfoDetails struct {
	Text       string `json:"text"`
	Link       string `json:"link"`
	SourceName string `json:"source_name"`
	SourceLogo string `json:"source_logo"`
}

// S3News is used for s3 news structure.
type S3News struct {
	Title          string                `json:"title"`
	NewsID         string                `json:"news_id"`
	Slug           string                `json:"slug"`
	BannerLink     string                `json:"banner_link"`
	ImageCourtesy  string                `json:"image_courtesy"`
	ShareAssetLink string                `json:"share_asset_link"`
	URL            string                `json:"url"`
	Content        string                `json:"content"`
	YoutubeLink    string                `json:"youtube_link"`
	VideoLink      string                `json:"video_link"`
	SourceLink     string                `json:"source_link"`
	Companies      []string              `json:"companies"`
	Commodities    []string              `json:"commodities"`
	Indices        []string              `json:"indices"`
	Forex          []string              `json:"forex"`
	Tags           []string              `json:"tags"`
	Notification   bool                  `json:"notification"`
	Lang           string                `json:"lang"`
	MoreInfo       S3NewsMoreInfoDetails `json:"more_info"`
	CreatedAt      string                `json:"created_at"`
	UpdatedAt      string                `json:"updated_at"`
	DeletedAt      string                `json:"deleted_at"`
}
