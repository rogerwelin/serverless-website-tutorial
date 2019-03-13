package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const apiUrl = "https://www.metaweather.com/api/location/"

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

var woeid = map[string]string{
	"london":    "44418",
	"stockholm": "906057",
}

func returnHighestPredictability(data *rawWeatherData) ([]uint8, error) {
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

	weather, err := json.Marshal(cleanedData)
	if err != nil {
		return nil, err
	}
	return weather, nil
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(apiUrl + woeid["stockholm"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data rawWeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	weather, err := returnHighestPredictability(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(weather)
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/hello", helloHandler)

	srv := &http.Server{
		Addr:         "0.0.0.0:5000",
		WriteTimeout: time.Second * 5,
		ReadTimeout:  time.Second * 5,
		Handler:      router,
	}
	log.Fatal(srv.ListenAndServe())
}
