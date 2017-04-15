package csv

import (
	"github.com/mgrzeszczak/playlist-converter/spotify"
	"encoding/csv"
	"os"
)

func WriteFile(filename string, tracks []spotify.Track) error {
	file,err := os.Create(filename)
	if err!=nil {
		return err
	}
	writer := csv.NewWriter(file)

	data := make([][]string,0)

	for _,track := range(tracks){
		row := make([]string,3)
		row[0] = track.Name
		row[1] = track.Album.Name
		row[2] = track.Artists[0].Name
		data = append(data,row)
	}

	err = writer.WriteAll(data)
	if err != nil {
		return err
	}
	writer.Flush()
	file.Close()
	return nil
}

