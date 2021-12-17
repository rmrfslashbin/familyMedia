package main

import (
	"encoding/json"
	"fetch/pkg/database"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

var log *logrus.Logger

func init() {
	// Set the log level
	log = logrus.New()
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

type AlbumDescription struct {
	BucketName  string `json:"bucketname"`
	Filename    string `json:"filename"`
	Description string `json:"description"`
}

const (
	bucket = "is-family-media"
	region = "us-east-1"
)

func main() {
	ddb := database.NewDDB(database.SetDDBRegion(region), database.SetDDBLogger(log))

	/*
		data, err := ddb.GetVideos(bucket)
		if err != nil {
			panic(err)
		}
		opts := make([]*AlbumDescription, len(data))
		for i, v := range data {
			opts[i] = &AlbumDescription{
				BucketName:  v.BucketName,
				Filename:    v.Filename,
				Description: v.Description,
			}
		}
		file, _ := json.MarshalIndent(opts, "", " ")
		_ = ioutil.WriteFile("exports/videos.json", file, 0644)
	*/

	jsonFile, err := os.Open("exports/videos.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var albumDescr []AlbumDescription
	json.Unmarshal(byteValue, &albumDescr)
	for _, v := range albumDescr {
		log.Infof("Updating description for video %s :: %s :: %s", v.BucketName, v.Filename, v.Description)
		err := ddb.AddVideoDescr(v.BucketName, v.Filename, v.Description)
		if err != nil {
			panic(err)
		}
	}

}
