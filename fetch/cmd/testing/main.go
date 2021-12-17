package main

import (
	"fetch/pkg/database"
	"strings"

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
	AlbumName   string `json:"albumName"`
	Bucket      string `json:"bucket"`
	Description string `json:"description"`
}

const (
	bucket = "is-family-media"
	region = "us-east-1"
)

func main() {
	//toTitle("CherylWeaver1970PhotoAlbum")

	//ddb.AddAlbumDescr("btest", "atest", "this is a descr!")

	ddb := database.NewDDB(database.SetDDBRegion(region), database.SetDDBLogger(log))
	data, err := ddb.GetImageAlbums(bucket)
	//data, err := ddb.GetImageLabels("is-family-media::BobAndPatWithKathy")
	//data, err := ddb.GetVideos(bucket)
	//data, err := ddb.GetVideo(bucket, "0000.00.00-BolerjackFamilyGraves")
	//data, err := ddb.GetImagesWithFace("910df256-5f76-4171-8d4c-d4725563be4e")
	//data, err := ddb.GetFaceWithImage("bdb58dc6-24db-4886-96f5-4a29b1e6c062", "is-family-media::JapaneseGarden08")
	if err != nil {
		log.Error(err)
	}

	/*albumDescrs := make([]AlbumDescription, len(data))
	for i, v := range data {
		albumDescrs[i] = AlbumDescription{
			AlbumName:   v.AlbumName,
			Bucket:      v.BucketName,
			Description: toTitle(v.AlbumName),
		}
	}

	file, _ := json.MarshalIndent(albumDescrs, "", " ")
	_ = ioutil.WriteFile("descrs.json", file, 0644)
	*/
}

func toTitle(s string) string {
	opt := ""
	for _, v := range s {
		if v >= 'A' && v <= 'Z' {
			//fmt.Println(i, string(v))
			opt += " " + string(v)
		} else if v == '-' {
			opt += ":"
		} else {
			opt += string(v)
		}
	}
	opt = strings.TrimSpace(opt)
	return opt
}
