package dtos

type Volume struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Image       string `json:"image"`
	Description string `json:"description"`
	Publisher   string `json:"publisher"`
}

type VolumeById struct {
	Volume     `json:"volume"`
	IssuesCount int     `json:"issues_count"`
	Issues     []Issue `json:"issues"`
}

type Issue struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IssueNumebr int    `json:"issue_number"`
	Image       string `json:"image"`
	CoverDate   string `json:"cover_date"`
	DateAdded   string `json:"date_added"`
}
