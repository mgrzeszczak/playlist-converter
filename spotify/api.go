package spotify

import (
	"github.com/franela/goreq"
	"fmt"
	"log"
	"encoding/json"
	"github.com/mgrzeszczak/playlist-converter/oauth2"
	"net/http"
)

const (
	api_url = "https://api.spotify.com"
)

func GetLibrary(auth *oauth2.AuthData) ([]Track,error) {
	tracks := make([]Track,0)

	decodeBody := func(b *goreq.Body) (*TracksData,error){
		data := &TracksData{}
		err := json.NewDecoder(b).Decode(data)
		if err!=nil {
			return nil,err
		}
		return data,nil
	}

	appendTracks := func (data *TracksData){
		for _,track := range(data.Items){
			tracks = append(tracks,track.Track)
		}
	}

	body,err := executeRequest(makeApiRequest("/v1/me/tracks?limit=50",http.MethodGet,auth.AccessToken))
	if err!=nil {
		return nil,err
	}

	data,err := decodeBody(body)
	if err != nil {
		return nil,err
	}

	appendTracks(data)

	for data.Next!="" {
		request := goreq.Request {
			Uri: data.Next,
			Method: http.MethodGet,
		}
		request.AddHeader("Authorization",fmt.Sprintf("Bearer %s",auth.AccessToken))
		body,err = executeRequest(request)
		if err!=nil {
			return nil,err
		}
		data,err = decodeBody(body)
		if err!=nil {
			return nil,err
		}

		appendTracks(data)
	}

	return tracks,nil
}

func executeRequest(r goreq.Request) (*goreq.Body,error) {
	resp,err := r.Do()
	if err!=nil {
		return nil,err
	}

	return resp.Body,nil
}

func makeApiRequest(path,method,token string) goreq.Request {
	request := goreq.Request {
		Uri: fmt.Sprintf("%s%s",api_url,path),
		Method: method,
	}
	request.AddHeader("Authorization",fmt.Sprintf("Bearer %s",token))
	log.Println(request)
	return request
}
