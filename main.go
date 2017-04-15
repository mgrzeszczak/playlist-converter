package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mgrzeszczak/playlist-converter/oauth2"
	"github.com/mgrzeszczak/playlist-converter/spotify"
	"github.com/mgrzeszczak/playlist-converter/utils"
	"github.com/mgrzeszczak/playlist-converter/youtube"
	"log"
)

const (
	user_library_read   = "user-library-read"
	user_library_modify = "user-library-modify"
	user_read_private   = "user-read-private"
	user_read_birthdate = "user-read-birthdate"
	user_read_email     = "user-read-email"
	user_top_read       = "user-top-read"

	spotify_code_url  = "https://accounts.spotify.com/authorize"
	spotify_token_url = "https://accounts.spotify.com/api/token"

	youtube_code_url  = "https://accounts.google.com/o/oauth2/v2/auth"
	youtube_token_url = "https://www.googleapis.com/oauth2/v4/token"
)

func main() {
	conf := utils.LoadConfig()

	spotAuth, err := loginSpotify(conf.Spotify)
	if err != nil {
		log.Fatalf("Login to spotify failed %v\n", err)
	}
	ytAuth, err := loginYoutube(conf.Youtube)
	if err != nil {
		log.Fatalf("Login to youtube failed %v\n", err)
	}

	data, err := spotify.GetLibrary(spotAuth)
	if err != nil {
		log.Fatalf("Failed to get spotify library: %v\n", err)
	}
	log.Printf("Found %d songs\n",len(data))

	playlistName, err := spotifyToYoutube(data,ytAuth)
	if err!=nil {
		log.Fatalf("Failed to export playlist %v\n",err)
	}

	log.Printf("Export successful\nPlaylist name: %v\n",playlistName)

	/*log.Printf("Library song count: %d\n", len(data))
	search, err := youtube.Search("star wars", ytAuth)
	if err != nil {
		log.Fatalf("Failed to get search results %v\n", err)
	}

	//res := utils.BestResult("star wars",search)

	log.Println(search[0])*/


	//log.Println(csv.WriteFile("output.csv",data))
}

func loginSpotify(credentials oauth2.Credentials) (*oauth2.AuthData, error) {
	spotifyAuthArgs := oauth2.AuthArgs{
		ClientSecret: credentials.ClientSecret,
		ClientId:     credentials.ClientId,
		AuthTokenUrl: spotify_token_url,
		AuthCodeUrl:  spotify_code_url,
		Scopes:       []string{user_library_read, user_read_private, user_read_email},
	}
	spotifyAuth, err := oauth2.Authorize(spotifyAuthArgs)
	if err != nil {
		return nil, err
	}
	return spotifyAuth, nil
}
func loginYoutube(credentials oauth2.Credentials) (*oauth2.AuthData, error) {
	ytAuthArgs := oauth2.AuthArgs{
		ClientSecret: credentials.ClientSecret,
		ClientId:     credentials.ClientId,
		Scopes:       []string{"https://www.googleapis.com/auth/youtube"},
		AuthCodeUrl:  youtube_code_url,
		AuthTokenUrl: youtube_token_url,
	}
	ytAuth, err := oauth2.Authorize(ytAuthArgs)
	if err != nil {
		return nil, err
	}
	return ytAuth, nil
}

func spotifyToYoutube(data []spotify.Track, ytAuth *oauth2.AuthData) (string, error) {
	playlistName := uuid.New().String()
	log.Printf("Creating playlist: %s\n",playlistName)
	playlistId, err := youtube.CreatePlaylist(playlistName, ytAuth)
	if err != nil {
		return "", err
	}
	log.Println("Playlist export started")
	for i, track := range data {

		log.Printf("[%d/%d] %s - %s",i+1,len(data), track.Name,track.Artists[0].Name)

		results, err := youtube.Search(fmt.Sprintf("%s %s", track.Name, track.Artists[0].Name), ytAuth)
		if err != nil {
			return "", err
		}

		videoId := results[0].Id.VideoId

		err = youtube.AddToPlaylist(playlistId, videoId, ytAuth)
		if err != nil {
			return "", err
		}
	}

	return playlistName, err
}
