package main

import (
	"errors"
	"log"
)

type Controller struct {
}

// TODO replace job status magic variables
func (c *Controller) AddJob(newJob Job) error {
	// Checks
	if newJob.RetryCount < 0 {
		return errors.New("Negative Retry Count")
	}
	if newJob.JobID < 0 {
		return errors.New("Negative ID")
	}

	// Init sane defaults
	newJob.Status = "waiting"
	if newJob.RetryIntervalMs == 0 {
		newJob.RetryIntervalMs = 1000
	}

	// TODO make a storage layer, eventually
	jobs = append(jobs, newJob)
	log.Printf("Jobs in queue: %+v\n", jobs)
	return nil
}

func (c *Controller) GetJob(JobID int) *Job {
	for _, job := range jobs {
		if job.JobID == JobID {
			return &job
		}
	}
	return nil
}

// TODO remove magicvariable
func (c *Controller) GetWaitingJob() *Job {
	for _, job := range jobs {
		if job.Status == "waiting" {
			job.Status = "in_progress"
			return &job
		}
	}
	return nil
}

func (c *Controller) UpdateJob(updatedJob Job) error {
	for idx, job := range jobs {
		if job.JobID == updatedJob.JobID {
			jobs[idx] = updatedJob
		}
	}
	return nil
}
