package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

const apiUrl = "https://www.metaweather.com/api/location/"

var (
	woeIDS = []string{"906057",
		"2487956",
		"44418",
		"638242",
		"2295420",
		"2151330",
		"615702",
		"721943",
		"727232",
		"766273"}

	S3Bucket = os.Getenv("S3BUCKET")
	S3Region = os.Getenv("REGION")
)

type rawWeatherData struct {
	ConsolidatedWeather []struct {
		WeatherStateName string  `json:"weather_state_name"`
		MinTemp          float64 `json:"min_temp"`
		MaxTemp          float64 `json:"max_temp"`
		Humidity         int     `json:"humidity"`
		Predictability   int     `json:"predictability"`
	} `json:"consolidated_weather"`
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	Woeid        int    `json:"woeid"`
	LattLong     string `json:"latt_long"`
}

type AggregatedWeather struct {
	WeatherItems []WeatherData `json:"weather_items"`
}

type WeatherData struct {
	Woeid            int     `json:"woeid"`
	WeatherStateName string  `json:"weather_state_name"`
	MinTemp          float64 `json:"min_temp"`
	MaxTemp          float64 `json:"max_temp"`
	Title            string  `json:"title"`
	LattLong         string  `json:"latt_long"`
}

func returnHighestPredictability(data *rawWeatherData) (WeatherData, error) {
	highestPredictability := 0
	var predictabilityIndex int

	for i, p := range data.ConsolidatedWeather {
		if p.Predictability > highestPredictability {
			highestPredictability = p.Predictability
			predictabilityIndex = i
		}
	}

	cleanedData := WeatherData{
		Woeid:            data.Woeid,
		WeatherStateName: data.ConsolidatedWeather[predictabilityIndex].WeatherStateName,
		MinTemp:          data.ConsolidatedWeather[predictabilityIndex].MinTemp,
		MaxTemp:          data.ConsolidatedWeather[predictabilityIndex].MaxTemp,
		Title:            data.Title,
		LattLong:         data.LattLong,
	}

	return cleanedData, nil
}

func (ag *AggregatedWeather) fetchApiData(woeID string, wg *sync.WaitGroup, mu *sync.Mutex) {
	defer wg.Done()

	var weather WeatherData
	var data rawWeatherData

	response, err := http.Get(apiUrl + woeID)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}
	defer response.Body.Close()

	if err := json.Unmarshal(body, &data); err != nil {
		return
	}

	weather, err = returnHighestPredictability(&data)
	if err != nil {
		return
	}

	mu.Lock()
	ag.WeatherItems = append(ag.WeatherItems, weather)
	mu.Unlock()
}

func uploadToS3(fileName string) error {
	sess := session.Must(session.NewSession(aws.NewConfig().WithRegion(S3Region)))
	uploader := s3manager.NewUploader(sess)

	file, err := os.Open(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	// upload to S3
	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(S3Bucket),
		Key:    aws.String("data.json"),
		Body:   file,
		ACL:    aws.String(s3.ObjectCannedACLPublicRead),
	})

	if err != nil {
		return err
	}

	fmt.Printf("File uploaded to: %s\n", aws.StringValue(&result.Location))
	return nil
}

func handleRequest(ctx context.Context, cloudWatchEvent events.CloudWatchEvent) error {
	var wg sync.WaitGroup
	var mu sync.Mutex
	ag := AggregatedWeather{}

	for _, item := range woeIDS {
		wg.Add(1)
		go ag.fetchApiData(item, &wg, &mu)
	}

	wg.Wait()
	re, err := json.Marshal(ag)
	if err != nil {
		return err
	}
	// fmt.Println(string(re))

	err = ioutil.WriteFile("/tmp/data.json", re, 0644)
	if err != nil {
		return err
	}
	err = uploadToS3("/tmp/data.json")
	if err != nil {
		return err
	}
	return nil

}

func main() {
	lambda.Start(handleRequest)
}
