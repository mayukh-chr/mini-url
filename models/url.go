package models

type URL struct {
	ID         int    `json:"id"`
	URL        string `json:"url"`
	ShortCode  string `json:"short_code"`
	AccessCount int   `json:"access_count"`
}
