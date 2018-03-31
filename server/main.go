package main

import (
	_ "strconv"
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
	initiateDB()
	go StartEndUserAPI()
	StartAgentAPI()
}

