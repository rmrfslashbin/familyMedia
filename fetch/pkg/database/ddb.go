package database

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/rekognition"
	"github.com/sirupsen/logrus"
)

type DDBOption func(config *DDBDriver)

const (
	t_family_media_faces        = "family-media-faces"
	t_family_media_image_albums = "family-media-image-albums"
	t_family_media_image_labels = "family-media-image-labels"
	t_family_media_videos       = "family-media-videos"
)

type DDBDriver struct {
	log        *logrus.Logger
	driverName string
	region     string
	db         *dynamodb.DynamoDB
}

type MediaFace struct {
	FaceId     string  `json:"faceId"`
	ObjectURI  string  `json:"objectUri"`
	Confidence float64 `json:"confidence"`
	ImageId    string  `json:"imageId"`
}

type MediaImageAlbum struct {
	BucketName  string   `json:"bucketName"`
	AlbumName   string   `json:"albumName"`
	Images      []string `json:"images"`
	Description string   `json:"description"`
}

type MediaVideo struct {
	BucketName      string `json:"bucketName"`
	Filename        string `json:"filename"`
	ThumbnailPrefix string `json:"thumbnailPrefix"`
	VideoPrefix     string `json:"videoPrefix"`
	Month           string `json:"month"`
	Year            string `json:"year"`
	VideoType       string `json:"videoType"`
	Day             string `json:"day"`
	Description     string `json:"description"`
	ThumbnailType   string `json:"thumbnailType"`
	Title           string `json:"title"`
}

type MediaImageLabel struct {
	ObjectURI  string                        `json:"objectUri"`
	Moderation []rekognition.ModerationLabel `json:"moderation"`
	Text       []rekognition.TextDetection   `json:"text"`
	Faces      []map[string]interface{}      `json:"faces"`
	Labels     []rekognition.Label           `json:"labels"`
}

func NewDDB(opts ...func(*DDBDriver)) *DDBDriver {
	config := &DDBDriver{}
	config.driverName = "ddb"

	// apply the list of options to Config
	for _, opt := range opts {
		opt(config)
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	// Create DynamoDB client
	svc := dynamodb.New(sess)
	config.db = svc

	return config
}

func SetDDBRegion(region string) func(*DDBDriver) {
	return func(config *DDBDriver) {
		config.region = region
	}
}

func SetDDBLogger(logger *logrus.Logger) func(*DDBDriver) {
	return func(config *DDBDriver) {
		config.log = logger
	}
}

func (config *DDBDriver) GetDriverName() string {
	return config.driverName
}

func (config *DDBDriver) GetImageAlbums(bucketName string) ([]*MediaImageAlbum, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":bucketName": {
				S: aws.String(bucketName),
			},
		},
		KeyConditionExpression: aws.String("bucketName = :bucketName"),
		TableName:              aws.String(t_family_media_image_albums),
	}

	result, err := config.db.Query(input)
	if err != nil {
		return nil, err
	}

	results := []*MediaImageAlbum{}
	for _, i := range result.Items {
		item := &MediaImageAlbum{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil

}

func (config *DDBDriver) GetImageLabels(objectURI string) ([]*MediaImageLabel, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":objectURI": {
				S: aws.String(objectURI),
			},
		},
		KeyConditionExpression: aws.String("objectURI = :objectURI"),
		TableName:              aws.String(t_family_media_image_labels),
	}

	result, err := config.db.Query(input)
	if err != nil {
		return nil, err
	}

	results := []*MediaImageLabel{}
	for _, i := range result.Items {
		item := &MediaImageLabel{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (config *DDBDriver) GetVideos(bucketName string) ([]*MediaVideo, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":bucketName": {
				S: aws.String(bucketName),
			},
		},
		KeyConditionExpression: aws.String("bucketName = :bucketName"),
		TableName:              aws.String(t_family_media_videos),
	}

	result, err := config.db.Query(input)
	if err != nil {
		return nil, err
	}

	results := []*MediaVideo{}
	for _, i := range result.Items {
		item := &MediaVideo{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (config *DDBDriver) GetVideo(bucketName string, filename string) ([]*MediaVideo, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":bucketName": {
				S: aws.String(bucketName),
			},
			":filename": {
				S: aws.String(filename),
			},
		},
		KeyConditionExpression: aws.String("bucketName = :bucketName and filename = :filename"),
		TableName:              aws.String(t_family_media_videos),
	}

	result, err := config.db.Query(input)
	if err != nil {
		return nil, err
	}

	results := []*MediaVideo{}
	for _, i := range result.Items {
		item := &MediaVideo{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (config *DDBDriver) GetImagesWithFace(faceId string) ([]*MediaFace, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":faceId": {
				S: aws.String(faceId),
			},
		},
		KeyConditionExpression: aws.String("faceId = :faceId"),
		TableName:              aws.String(t_family_media_faces),
	}

	result, err := config.db.Query(input)
	if err != nil {
		return nil, err
	}

	results := []*MediaFace{}
	for _, i := range result.Items {
		item := &MediaFace{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}
	return results, nil
}

func (config *DDBDriver) GetFaceWithImage(faceId string, objectURI string) ([]*MediaFace, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":faceId": {
				S: aws.String(faceId),
			},
			":objectURI": {
				S: aws.String(objectURI),
			},
		},
		KeyConditionExpression: aws.String("faceId = :faceId and objectURI = :objectURI"),
		TableName:              aws.String(t_family_media_faces),
	}

	result, err := config.db.Query(input)
	if err != nil {
		return nil, err
	}

	results := []*MediaFace{}
	for _, i := range result.Items {
		item := &MediaFace{}
		err = dynamodbattribute.UnmarshalMap(i, &item)
		if err != nil {
			return nil, err
		}
		results = append(results, item)
	}

	return results, nil
}

func (config *DDBDriver) AddAlbumDescr(bucketName string, albumName string, descr string) error {
	input := &dynamodb.UpdateItemInput{

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":descr": {
				S: aws.String(descr),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"bucketName": {
				S: aws.String(bucketName),
			},
			"albumName": {
				S: aws.String(albumName),
			},
		},
		UpdateExpression: aws.String("set description = :descr"),
		TableName:        aws.String(t_family_media_image_albums),
	}

	_, err := config.db.UpdateItem(input)
	return err
}

func (config *DDBDriver) AddVideoDescr(bucketName string, filename string, description string) error {
	input := &dynamodb.UpdateItemInput{

		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":description": {
				S: aws.String(description),
			},
		},
		Key: map[string]*dynamodb.AttributeValue{
			"bucketName": {
				S: aws.String(bucketName),
			},
			"filename": {
				S: aws.String(filename),
			},
		},
		UpdateExpression: aws.String("set description = :description"),
		TableName:        aws.String(t_family_media_videos),
	}

	_, err := config.db.UpdateItem(input)
	return err
}
