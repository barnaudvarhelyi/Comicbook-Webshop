package models

import "encoding/xml"

type IssuesResponse struct {
	XMLName    xml.Name `xml:"response"`
	Issues     []Issue  `xml:"results"`
	IssueError string   `xml:"error"`
}

type Issue struct {
	ID                   int                  `xml:"id"`
	Name                 string               `xml:"name"`
	Image                IssueImage           `xml:"image"`
	CoverDate            string               `xml:"cover_date"`
	Date_Added           string               `xml:"date_added"`
	IssueNumber          int                  `xml:"issue_number"`
}

type IssueImage struct {
	OriginalURL string `xml:"original_url"`
}
