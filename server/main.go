package main

import (
	_ "strconv"
	"log"
	"github.com/boltdb/bolt"
	_ "encoding/json"
)

type Job struct {
	JobID           string
	Prioritiy       int
	Cmd             string
	RetryCount      int
	RetryIntervalMs int
	Status          string // status waiting/failed/finished/in_progress
	OutputsSTDOUT   string
	OutputsSTDERR   string
}

func main() {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	
	err = db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte("jobs"))
		if err != nil {
			log.Println("could not create jobs bucket")
			return err
		}

		return err
	})

	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("jobs"))
		v := b.Get([]byte("69"))
		log.Printf("The answer is: " + string( v[:] ))
		return nil
	})
	db.Close()

	go StartEndUserAPI()
	StartAgentAPI()
}
