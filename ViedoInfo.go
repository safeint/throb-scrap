package main

import "encoding/xml"

type VideoInfo struct {
	Actor         []Actor   `xml:"actor"`
	Aired         string    `xml:"aired"`
	Art           Art       `xml:"art"`
	AspectRatio   string    `xml:"aspectratio"`
	Biography     string    `xml:"biography"`
	Continent     string    `xml:"continent"`
	Cover         string    `xml:"cover"`
	Country       string    `xml:"country"`
	CountryCode   int       `xml:"countrycode"`
	CriticRating  float32   `xml:"criticrating"`
	Credits       []string  `xml:"credits"`
	CustomRating  float32   `xml:"customrating"`
	DateAdded     string    `xml:"dateadded"`
	DisplayOrder  int       `xml:"displayorder"`
	Director      string    `xml:"director"`
	EndDate       string    `xml:"enddate"`
	FileInfo      string    `xml:"fileinfo"`
	FanArt        []string  `xml:"fanart"`
	Formed        string    `xml:"formed"`
	Genre         []string  `xml:"genre"`
	Language      string    `xml:"language"`
	LastPlayed    string    `xml:"lastplayed"`
	LockData      bool      `xml:"lockdata"`
	LockedFields  string    `xml:"lockedfields"`
	LocalTitle    string    `xml:"localtitle"`
	Maker         string    `xml:"maker"`
	Mpaa          string    `xml:"mpaa"`
	OriginalTitle string    `xml:"originaltitle"`
	Plot          string    `xml:"plot"`
	PlayCount     int       `xml:"playcount"`
	PlotOutline   string    `xml:"plotoutline"`
	Poster        string    `xml:"poster"`
	Premiered     string    `xml:"premiered"`
	Rating        float32   `xml:"rating"`
	Ratings       []Ratings `xml:"ratings"`
	ReleaseDate   string    `xml:"releasedate"`
	Runtime       int       `xml:"runtime"`
	SortName      string    `xml:"sortname"`
	SortTitle     string    `xml:"sorttitle"`
	Style         []Inner   `xml:"style"`
	Studio        string    `xml:"studio"`
	Tag           []string  `xml:"tag"`
	TagLine       string    `xml:"tagline"`
	Thumb         string    `xml:"thumb"`
	Title         string    `xml:"title"`
	Trailer       []string  `xml:"trailer"`
	UniqueId      string    `xml:"uniqueid"`
	Watched       bool      `xml:"watched"`
	Writer        []string  `xml:"writer"`
	XMLName       xml.Name  `xml:"movie"`
	Year          int       `xml:"year"`
}

type Inner struct {
	Inner string `xml:"innerxml"`
}

type Actor struct {
	Name      string `xml:"name"`
	Role      string `xml:"role"`
	Type      string `xml:"type"`
	SortOrder int    `xml:"sortorder"`
	Thumb     string `xml:"thumb"`
}

type Ratings struct {
	Name   string  `xml:"name"`
	Rating float32 `xml:"rating"`
}

type Art struct {
	Poster string   `xml:"poster"`
	FanArt []string `xml:"fanart"`
}
