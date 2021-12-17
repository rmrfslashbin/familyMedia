package main

import (
	"fetch/pkg/database"
	"fmt"
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

const (
	bucket = "is-family-media"
	region = "us-east-1"
)

type VideoRecord struct {
	Slug      string `yaml:"slug"`
	Title     string `yaml:"title"`
	HLSFile   string `yaml:"hlsFile"`
	Date      string `yaml:"date"`
	Thumbnail string `yaml:"thumbnail"`
}

func main() {

	ddb := database.NewDDB(database.SetDDBRegion(region), database.SetDDBLogger(log))
	data, err := ddb.GetVideos(bucket)
	if err != nil {
		log.Error(err)
	}
	records := make([]*VideoRecord, len(data))
	for i, v := range data {
		hlsFile := fmt.Sprintf("https://%s.s3.amazonaws.com/hls/%s/%s.m3u8", v.BucketName, v.Filename, v.Filename)
		thumbnail := fmt.Sprintf("https://%s.s3.amazonaws.com/%s/%s.%s", v.BucketName, v.ThumbnailPrefix, v.Filename, v.ThumbnailType)
		dt := fmt.Sprintf("%s-%s-%s", v.Year, v.Month, v.Day)
		records[i] = &VideoRecord{
			Slug:      v.Filename,
			Title:     v.Description,
			HLSFile:   hlsFile,
			Date:      dt,
			Thumbnail: thumbnail,
		}
	}

	for _, v := range records {
		file, err := yaml.Marshal(v)
		if err != nil {
			spew.Dump(v)
			log.Fatalf("error: %v", err)
		}
		file = append([]byte("---\n"), file...)
		file = append(file, []byte("---\n{{< video >}}\n")...)
		_ = ioutil.WriteFile(path.Join("exports", "videos", v.Slug+".md"), file, 0644)
	}

}

//https://is-family-media.s3.amazonaws.com/hls/0000.00.00-BolerjackFamilyGraves/0000.00.00-BolerjackFamilyGraves.m3u8
