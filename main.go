package main

import (
	"encoding/json"
	"net/http"
	"github.com/lithammer/shortuuid"
	"github.com/redis/go-redis/v9"
	"flag"
	"log"
	"time"
	"context"
	"fmt"
	"github.com/gorilla/mux"
)

var redisUrl = flag.String("u", "redis://redis:6379", "Redis url")
var serverPort = flag.String("p", "3000", "The server running port")

type HttpController struct {
	client *redis.Client
}

type RequestBody struct {
	Url string `json:"url"`
	Ttl uint   `json:"ttl"`
}

func (ctrl* HttpController) redirectHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	log.Printf("searching for url uuid %s shortened", vars["uuid"])

	url, err := ctrl.client.Get(context.Background(), vars["uuid"]).Result()

	// Get 
	if  err != nil {
		log.Println(err.Error())
		http.Error(w, "Cannot redirect. URL not found.", http.StatusNotFound)
		return
	}

	// Perform a temporary redirect (HTTP 302) to another URL
	http.Redirect(w, r, url, http.StatusFound)
}

func (ctrl* HttpController) shortenUrlHandler(w http.ResponseWriter, r *http.Request) {
	var body RequestBody
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&body); err != nil {
		log.Println(err.Error())
		http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
	}

	uuid := shortuuid.New()

	log.Printf("url: %s\n ttl: %d", body.Url, body.Ttl)
	if err := ctrl.client.Set(context.Background(), uuid, body.Url, time.Duration(body.Ttl) * time.Second).Err(); err != nil {
		log.Println(err.Error())
		http.Error(w, "InternalServerHerror", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintln(w, uuid);
}

func main() {
	flag.Parse()
	
	opt, err := redis.ParseURL(*redisUrl)
	if err != nil {
		panic(err)
	}
	client := redis.NewClient(opt)
	defer client.Close()

	controller := &HttpController{client}
	router := mux.NewRouter()

	router.HandleFunc("/{uuid}", controller.redirectHandler).Methods("GET")
	router.HandleFunc("/generate", controller.shortenUrlHandler).Methods("POST")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", *serverPort), router))
}
