package persistence

import ()

type LanguageTags struct {
	Tags []Tag `xml:"tag"`
}

type Tag struct {
	Id      string `xml:"id,attr"`
	Message string `xml:",chardata"`
}
