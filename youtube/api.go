package youtube

import (
	"encoding/json"
	"github.com/mgrzeszczak/playlist-converter/oauth2"
	"github.com/franela/goreq"
	"fmt"
	"net/url"
)

const (
	api_url = "https://www.googleapis.com/youtube/v3/"
)

func CreatePlaylist(name string, auth *oauth2.AuthData) (string, error) {
	req := goreq.Request{
		Method: "POST",
		Uri: api_url+"playlists?part=snippet",
		Body:createPlaylistData{
			Snippet:snippet{Title:name},
		},
		ContentType:"application/json",
	}
	req.AddHeader("Authorization", fmt.Sprintf("Bearer %s", auth.AccessToken))

	resp, err := req.Do()
	if err != nil {
		return "", err
	}

	data := &createPlaylistResponse{}
	err = json.NewDecoder(resp.Body).Decode(data)
	if err != nil {
		return "", err
	}
	return data.Id, nil
}

func AddToPlaylist(playlistId, videoId string, auth *oauth2.AuthData) error {
	req := goreq.Request{
		Method: "POST",
		Uri:api_url+"playlistItems?part=snippet",
		ContentType:"application/json",
		Body:struct {
			Snippet snippet `json:"snippet"`
		}{
			snippet{
				PlaylistId:playlistId,
				Resource:resource{
					VideoId:videoId,
					Kind:"youtube#video",
				},
			},
		},
	}
	req.AddHeader("Authorization", fmt.Sprintf("Bearer %s", auth.AccessToken))
	_, err := req.Do()
	return err
}


func Search(query string, auth *oauth2.AuthData) ([]SearchItem,error){
	req := goreq.Request{
		Method: "GET",
		Uri: fmt.Sprintf(api_url+"search?part=snippet&q=%s",url.PathEscape(query)),
	}
	req.AddHeader("Authorization",fmt.Sprintf("Bearer %s",auth.AccessToken))
	resp,err := req.Do()
	if err!=nil {
		return nil,err
	}
	results := &SearchResults{}
	err = json.NewDecoder(resp.Body).Decode(results)
	if err!=nil{
		return nil,err
	}
	return results.Items,nil
}