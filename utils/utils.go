package utils

import (
	"github.com/mgrzeszczak/playlist-converter/youtube"
	"github.com/texttheater/golang-levenshtein/levenshtein"
	"github.com/mgrzeszczak/playlist-converter/oauth2"
	"os"
	"log"
	"encoding/json"
)

const (
	config_name ="config.json"
)

type holder struct {
	Item youtube.SearchItem
	Distance int
}

func BestResult(query string,res []youtube.SearchItem) youtube.SearchItem {

	dist := make(map[string]holder)
	for _,v := range(res){
		d := levenshtein.DistanceForStrings([]rune(query),[]rune(v.Snippet.Title),levenshtein.DefaultOptions)
		dist[v.Snippet.Title] = holder{
			Item : v,
			Distance : d,
		}
	}

	min := dist[res[0].Snippet.Title]
	for _,v := range(dist){
		if v.Distance<min.Distance{
			min = v
		}
	}

	return min.Item
}

func LoadConfig() *oauth2.Config {
	f,err := os.Open(config_name)
	if err!=nil{
		log.Fatalf("Cannot open config file: %v\n",err)
	}
	conf := &oauth2.Config{}
	err = json.NewDecoder(f).Decode(conf)
	if err!=nil {
		log.Fatalf("Failed to parse config file %v\n",err)
	}
	return conf
}