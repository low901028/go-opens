package main

import (
	"log"
	"net/http"
	"github.com/low901028/go-opens/httpbatchlib"
)

func main(){
	mux := http.NewServeMux()

	mux.HandleFunc("/batch", func(w http.ResponseWriter, r *http.Request) {
		status := httpbatchlib.ValidRequest(r)
		if status != 0{
			w.WriteHeader(status)
		}
		requests, err := httpbatchlib.UnpackRequests(r)
		if err != nil{
			log.Printf("Error unpacking the request: %s", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		responses := httpbatchlib.BuildResponses(requests)
		if responses == nil || len(responses) == 0{
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if err = httpbatchlib.WriteResponses(w, responses); err != nil{
			log.Printf("Error packing the responses: %s", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	log.Fatal(http.ListenAndServe(":9898", mux))
}
