package youtube

type createPlaylistData struct {
	Snippet snippet `json:"snippet"`
}

type snippet struct {
	Title      string `json:"title"`
	PlaylistId string `json:"playlistId,omitempty"`
	Resource   resource `json:"resourceId,omitempty"`
}

type resource struct {
	VideoId string `json:"videoId,omitempty"`
	Kind    string `json:"kind"`
}

type createPlaylistResponse struct {
	Id string `json:"id"`
}

type SearchResults struct {
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	Id Id `json:"id"`
	Snippet snippet `json:"snippet"`
}

type Id struct {
	Kind string `json:"kind"`
	VideoId string `json:"videoId"`
}
