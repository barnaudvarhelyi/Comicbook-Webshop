package models

import "encoding/xml"

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

type VolumeResponse struct {
	XMLName       xml.Name `xml:"response"`
	VolumesError  string   `xml:"error"`
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
}

type CharacterApiDetail struct {
	Character             xml.Name `xml:"character"`
	CharacterApiDetailUrl string   `xml:"api_detail_url"`
}
