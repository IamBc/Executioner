package main

import (
	"errors"
	"log"
	"encoding/json"
	"github.com/boltdb/bolt"
)

type Controller struct {
}

func (c *Controller) AddJob(newJob Job) error {
	// Checks
	if newJob.RetryCount < 0 {
		return errors.New("Negative Retry Count")
	}
	if newJob.JobID == "" {
		return errors.New("Negative ID")
	}

	// Init sane defaults
	newJob.Status = "waiting"
	if newJob.RetryIntervalMs == 0 {
		newJob.RetryIntervalMs = 1000
	}

	buf, err := json.Marshal( newJob )
	if err != nil {
		log.Println( "Could not marshal newJob" )
		return err
	}
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("jobs"))
		err := b.Put([]byte(newJob.JobID), buf )
		if err != nil {
			log.Println( "Could not insert new Job with id: " + newJob.JobID )
			return err
		}
		return nil
	})
	db.Close()
	log.Println( "Succesfully added Job with ID: " + newJob.JobID )
	return nil
}

func (c *Controller) GetJob(JobID string) *Job {
	var existingJob Job

	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("jobs"))
		v := b.Get([]byte( JobID ))
		log.Printf("The existing job is: " + string( v[:] ))
		json.Unmarshal(v, &existingJob)
		return nil
	})
	db.Close()

	return &existingJob
}

func (c *Controller) GetWaitingJob() *Job {
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	var waitingJob Job
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte("jobs"))

		b.ForEach(func(k, v []byte) error {
			log.Printf("key=%s, value=%s\n", k, v)
			var job Job
			err := json.Unmarshal(v, &job)
			if err != nil {
				log.Println( "Error unmarshalling JSON" )
				return err
			}

			if job.Status == "waiting" {
				err = json.Unmarshal(v, &waitingJob)
				if err != nil {
					log.Println( "Error unmarshalling JSON" )
					return err
				}
			}
			return nil
		})
		return nil
	})
	db.Close()

	return &waitingJob
}

func (c *Controller) UpdateJob(updatedJob Job) error {
	if updatedJob.JobID == "" {
		return errors.New("Negative ID")
	}
 
	buf, err := json.Marshal( updatedJob )
	if err != nil {
		log.Println( "Could not marshal newJob" )
		return err
	}
	db, err := bolt.Open("my.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("jobs"))
		err := b.Put([]byte(updatedJob.JobID), buf )
		if err != nil {
			log.Println( "Could not update Job with id: " + updatedJob.JobID )
			return err
		}
		return nil
	})
	db.Close()
	log.Println( "Succesfully updated Job with ID: " + updatedJob.JobID )
	return nil
}
