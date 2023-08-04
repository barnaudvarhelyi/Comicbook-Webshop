package models

import "encoding/xml"

type CharactersResponse struct {
	XMLName        xml.Name    `xml:"response"`
	Characters     []Character `xml:"results"`
	CharacterError string      `xml:"error"`
}

type Character struct {
	ID    int            `xml:"id"`
	Name  string         `xml:"name"`
	Image CharacterImage `xml:"image"`
}

type CharacterImage struct {
	OriginalURL string `xml:"original_url"`
}
