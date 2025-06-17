package models


type Track struct {
	Title  string `json:"title"`
	ReleaseDate string `json:"release_date"`
	Duration int    `json:"duration"`
	Artist struct {
		Name string `json:"name"`
	} `json:"artist"`
	Album struct {
		Title string `json:"title"`
		Cover string `json:"cover_big"`
	} `json:"album"`
}