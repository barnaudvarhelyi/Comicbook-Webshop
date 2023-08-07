package models

import "encoding/xml"

//All volumes models for API request
type VolumesResponse struct {
	XMLName        xml.Name       `xml:"response"`
	VolumesResults VolumesResults `xml:"results"`
}

type VolumesResults struct {
	XMLName    xml.Name             `xml:"results"`
	VolumesUrl []VolumeApiDetailUrl `xml:"volume"`
}

type VolumeApiDetailUrl struct {
	VolumeApiDetailUrl string `xml:"api_detail_url"`
}

//Volume models for API request
type VolumeResponse struct {
	XMLName       xml.Name `xml:"response"`
	VolumeResults []Volume `xml:"results"`
}

type Volume struct {
	Volume               xml.Name             `xml:"volume"`
	ID                   int                  `xml:"id"`
	Name                 string               `xml:"name"`
	Image                VolumeImage          `xml:"image"`
	Description          string               `xml:"description"`
	Publisher            Publisher            `xml:"publisher"`
	IssuesApiDetails     []IssueApiDetails    `xml:"issues>issue"`
	CharactersApiDetails []CharacterApiDetail `xml:"characters>character"`
}

type VolumeImage struct {
	OriginalURL string `xml:"original_url"`
}

type Publisher struct {
	Name string `xml:"name"`
}

type IssueApiDetails struct {
	Issue             xml.Name `xml:"issue"`
	IssueApiDetailUrl string   `xml:"api_detail_url"`
	ID                int      `xml:"id"`
}

type CharacterApiDetail struct {
	Character             xml.Name `xml:"character"`
	CharacterApiDetailUrl string   `xml:"api_detail_url"`
}

//Issue models for API request
type IssuesResponse struct {
	XMLName    xml.Name `xml:"response"`
	Issues     []Issue  `xml:"results"`
}

type Issue struct {
	ID          int        `xml:"id"`
	Name        string     `xml:"name"`
	Image       IssueImage `xml:"image"`
	CoverDate   string     `xml:"cover_date"`
	Date_Added  string     `xml:"date_added"`
	IssueNumber int        `xml:"issue_number"`
}

type IssueImage struct {
	OriginalURL string `xml:"original_url"`
}

//Character models for API request
type CharactersResponse struct {
	XMLName        xml.Name    `xml:"response"`
	Characters     []Character `xml:"results"`
}

type Character struct {
	ID           int            `xml:"id"`
	Name         string         `xml:"name"`
	Image        CharacterImage `xml:"image"`
	IssueCredits IssueCredits   `xml:"issue_credits"`
}

type CharacterImage struct {
	OriginalURL string `xml:"original_url"`
}

type IssueCredits struct {
	IssueCredit []IssueApiDetails `xml:"issue"`
}

//Error model for API request
type ApiError struct {
	XMLName xml.Name `xml:"response"`
	Error   string   `xml:"error"`
}
