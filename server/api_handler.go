package main

// handles input format and arg types validation before passing to controller

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

func StartEndUserAPI() {
	log.Println("Starting EndUser API...")
	var apiHandler ApiHandler
	router := mux.NewRouter().StrictSlash(false)
	router.HandleFunc("/AddJob/", apiHandler.AddJob)
	router.HandleFunc("/GetJobInfo/", apiHandler.GetJobInfo)
	// TODO addr, timeouts should be in config
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

type ApiHandler struct {
}

func (m *ApiHandler) AddJob(w http.ResponseWriter, r *http.Request) {

	decoder := json.NewDecoder(r.Body)
	var newJob Job
	err := decoder.Decode(&newJob)
	if err != nil {
		log.Println("Decoding Error: " + err.Error())
		w.Write([]byte("Wrong or missing parameters!"))
	}
	defer r.Body.Close()

	var c Controller
	err = c.AddJob(newJob)
	if err != nil {
		log.Println(err.Error())
		w.Write([]byte(err.Error()))
	}
}

func (m *ApiHandler) GetJobInfo(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var requestedJob Job
	err := decoder.Decode(&requestedJob)
	if err != nil {
		log.Println("Decoding Error: " + err.Error())
		w.Write([]byte("Wrong or missing parameters!"))
	}
	defer r.Body.Close()

	var c Controller
	job := c.GetJob(requestedJob.JobID)
	log.Printf("Job details: %+v\n", job)

	// TODO func
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(jobs)
	w.Write(b.Bytes())
}
