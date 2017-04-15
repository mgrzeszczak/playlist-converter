package spotify

type TracksData struct {
	Total int `json:"total"`
	Next string `json:"next"`
	Previous string `json:"previous"`
	Href string `json:"href"`
	Limit int `json:"limit"`
	Offset int `json:"offset"`
	Items []TrackInfo `json:"items"`
}

type TrackInfo struct {
	AddedAt string `json:"added_at"`
	Track Track `json:"track"`
}

type Track struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Album Album `json:"album"`
	Artists []Artist `json:"artists"`
	Duration int `json:"duration_ms"`
}

type Album struct {
	Id string `json:"id"`
	Name string `json:"name"`
	Artists []Artist `json:"artists"`
}

type Artist struct {
	Name string `json:"name"`
	Id string `json:"id"`
}
