package browser

type Landing struct {
	Type          string  `json:"type"`
	ID            int64   `json:"id"`
	Nmae          string  `json:"name"`
	Logo          *string `json:"logo"`
	Link          string  `json:"link"`
	Subscribers   *int64  `json:"subscribers"`
	IsInactive    bool    `json:"IsInactive"`
	IsDeactivated bool    `json:"isDeactivated"`
}

type Video struct {
	ID                 int64   `json:"id"`
	Duration           int     `json:"duration"`
	Created            int64   `json:"created"`
	Title              string  `json:"title"`
	ThumbID            int64   `json:"thumbId"`
	VideoType          string  `json:"videoType"`
	PageURL            string  `json:"pageURL"`
	ThumbURL           string  `json:"thumbURL"`
	ImageURL           string  `json:"imageURL"`
	PreviewThumbURL    string  `json:"previewThumbURL"`
	SpriteURL          string  `json:"spriteURL"`
	TrailerURL         string  `json:"trailerURL"`
	Views              int64   `json:"views"`
	Landing            Landing `json:"landing"`
	IsUHD              bool    `json:"isUHD"`
	IsWatched          bool    `json:"isWatched"`
	IsThumbCustom      bool    `json:"isThumbCustom"`
	IsAdminCustomThumb bool    `json:"isAdminCustomThumb"`
}

type VideoList struct {
	Videos []Video `json:"videos"`
}
