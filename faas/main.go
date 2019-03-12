package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const apiUrl = "https://www.metaweather.com/api/location/"

type WeatherData struct {
	ConsolidatedWeather []struct {
		ID               int64   `json:"id"`
		WeatherStateName string  `json:"weather_state_name"`
		MinTemp          float64 `json:"min_temp"`
		MaxTemp          float64 `json:"max_temp"`
		WindDirection    float64 `json:"wind_direction"`
		Humidity         int     `json:"humidity"`
		Predictability   int     `json:"predictability"`
	} `json:"consolidated_weather"`
	Title        string `json:"title"`
	LocationType string `json:"location_type"`
	Woeid        int    `json:"woeid"`
	LattLong     string `json:"latt_long"`
}

var woeid = map[string]string{
	"london":    "44418",
	"stockholm": "906057",
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(apiUrl + woeid["london"])
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

	var data WeatherData
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, data)
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
