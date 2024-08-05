package models

// IdentityData is used to store the user information.
type IdentityData struct {
	CountryCode string `json:"country_code"`
	Mobile      string `json:"mob_no"`
	UserID      string `json:"user_id"`
	Source      string `json:"source"`
	AppID       string `json:"app_id"`
}
