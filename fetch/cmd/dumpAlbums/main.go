package main

import (
	"fetch/pkg/database"
	"io/ioutil"
	"path"

	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
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

type Album struct {
	URLRoot string   `yaml:"url_root"`
	Slug    string   `yaml:"slug"`
	Type    string   `yaml:"type"`
	Sizes   []string `yaml:"sizes"`
	Images  []string `yaml:"images"`
}

type AlbumRecord struct {
	Title string `yaml:"title"`
	Album *Album `yaml:"album"`
}

const (
	bucket = "is-family-media"
	region = "us-east-1"
)

func main() {
	//toTitle("CherylWeaver1970PhotoAlbum")

	ddb := database.NewDDB(database.SetDDBRegion(region), database.SetDDBLogger(log))
	data, err := ddb.GetImageAlbums(bucket)
	if err != nil {
		log.Error(err)
	}

	for _, v := range data {
		file, err := yaml.Marshal(&AlbumRecord{
			Title: v.Description,
			Album: &Album{
				URLRoot: "https://is-family-media.s3.amazonaws.com/images",
				Slug:    v.AlbumName,
				Type:    "png",
				Sizes:   []string{"source", "256", "768"},
				Images:  v.Images,
			},
		})
		if err != nil {
			spew.Dump(v)
			log.Fatalf("error: %v", err)
		}
		file = append([]byte("---\n"), file...)
		file = append(file, []byte("---\n{{< album >}}\n")...)
		_ = ioutil.WriteFile(path.Join("exports", "albums", v.AlbumName+".md"), file, 0644)
	}
}
