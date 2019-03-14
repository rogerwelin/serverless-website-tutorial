package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/lambda"
)

const apiUrl = "https://www.metaweather.com/api/location/"

type Event struct {
	WoeID string `json:"woeid"`
}

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

type WeatherData struct {
	WeatherStateName string  `json:"weather_state_name"`
	MinTemp          float64 `json:"min_temp"`
	MaxTemp          float64 `json:"max_temp"`
	Title            string  `json:"title"`
	LattLong         string  `json:"latt_long"`
}

func returnHighestPredictability(data *rawWeatherData) (*WeatherData, error) {
	highestPredictability := 0
	var predictabilityIndex int

	for i, p := range data.ConsolidatedWeather {
		if p.Predictability > highestPredictability {
			highestPredictability = p.Predictability
			predictabilityIndex = i
		}
	}

	cleanedData := &WeatherData{
		WeatherStateName: data.ConsolidatedWeather[predictabilityIndex].WeatherStateName,
		MinTemp:          data.ConsolidatedWeather[predictabilityIndex].MinTemp,
		MaxTemp:          data.ConsolidatedWeather[predictabilityIndex].MaxTemp,
		Title:            data.Title,
		LattLong:         data.LattLong,
	}

	return cleanedData, nil
}

func HandleRequest(ctx context.Context, req Event) (WeatherData, error) {

	var weather *WeatherData

	resp, err := http.Get(apiUrl + req.WoeID)
	if err != nil {
		return *weather, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return *weather, err
	}
	defer resp.Body.Close()

	var data rawWeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		return *weather, err
	}

	weather, err = returnHighestPredictability(&data)
	if err != nil {
		return *weather, err
	}

	return *weather, nil
}

func main() {
	lambda.Start(HandleRequest)
}
