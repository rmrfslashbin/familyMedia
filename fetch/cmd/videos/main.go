package main

import (
	"fetch/pkg/database"

	"github.com/davecgh/go-spew/spew"
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

func main() {
	ddb := database.NewDDB(database.SetDDBRegion("us-east-1"), database.SetDDBLogger(log))
	data, err := ddb.GetImageAlbums("is-family-media")
	//data, err := ddb.GetImageLabels("is-family-media::BobAndPatWithKathy")
	//data, err := ddb.GetVideos("is-family-media")
	//data, err := ddb.GetVideo("is-family-media", "0000.00.00-BolerjackFamilyGraves")
	//data, err := ddb.GetImagesWithFace("910df256-5f76-4171-8d4c-d4725563be4e")
	//data, err := ddb.GetFaceWithImage("bdb58dc6-24db-4886-96f5-4a29b1e6c062", "is-family-media::JapaneseGarden08")
	if err != nil {
		log.Error(err)
	}
	spew.Dump(data)
}
