package models

import "encoding/xml"

type ApiError struct {
	XMLName xml.Name `xml:"response"`
	Error   string   `xml:"error"`
}
