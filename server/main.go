package main

import (
	_ "strconv"
)

// TODO abstract everything with interfaces

var jobs []Job

type Job struct {
	JobID           int
	Prioritiy       int
	OwnerHash       string //TODO config defaults by ownerHash
	Cmd             string
	RetryCount      int
	RetryIntervalMs int
	Status          string // status waiting/failed/finished/in_progress
	OutputsSTDOUT   string
	OutputsSTDERR   string
}

func main() {
	jobs = make([]Job, 0)
	
	go StartEndUserAPI()
	StartAgentAPI()
}
