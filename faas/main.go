package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
)

const (
	BaseURL = "https://ghibliapi.herokuapp.com"
)

var (
	wait time.Duration
)

type Movie struct {
	ID           string `json: "id"`
	Title        string `json: "title"`
	Description  string `json: "description"`
	Director     string `json: "director"`
	Release_Date string `json: "release_date"`
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "pong\n")
}

func movies(w http.ResponseWriter, r *http.Request) {
	resp, err := http.Get(BaseURL + "/films")
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	var mov []Movie
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &mov)
	if err != nil {
		log.Println(err)
		return
	}
	apa, _ := json.Marshal(mov)

	fmt.Fprintln(w, string(apa))
}

func movie(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	resp, err := http.Get(BaseURL + "/films/" + id)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var mov Movie
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &mov)

	apa, _ := json.Marshal(mov)

	fmt.Fprintln(w, string(apa))

}

func main() {
	wait = time.Second * 5

	router := mux.NewRouter()
	router.HandleFunc("/", rootHandler).Methods("GET")
	router.HandleFunc("/movies", movies).Methods("GET")
	router.HandleFunc("/movie/{id}", movie).Methods("GET")

	srv := &http.Server{
		Addr:         "0.0.0.0:4000",
		WriteTimeout: wait,
		ReadTimeout:  wait,
		Handler:      router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// graceful shutdown at SIGINT
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	srv.Shutdown(ctx)

	log.Println("Shutting down server...")
	os.Exit(0)
}
