package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
)

// TODO Config
// TODO abstract everything with interfaces

type Job struct {
	JobID           int
	Prioritiy       int
	OwnerHash       string
	Cmd             string
	RetryCount      int
	RetryIntervalMs int
	Status          string // status inactive/failed/finished/in_progress
	OutputsSTDOUT   string
	OutputsSTDERR   string
}

func main() {
	for {
		log.Println("Starting agent...")

		job, err := GetJob()
		if err != nil {
			log.Println("Error getting job: " + err.Error())
			time.Sleep(3 * time.Second)
			continue
		}

		if job.JobID == 0 {
			log.Println("There is no current JOB!")
			time.Sleep(3 * time.Second)
			continue
		}

		// Execute Job
		stdout, stderr, err := ExecCommand(job.Cmd)
		if err != nil {
			log.Println("Error executing job: " + err.Error())

			if job.RetryCount < 1 {
				job.RetryCount -= 1
				continue
			}
			job.Status = "failed"
		} else {
			job.Status = "finished"
		}

		log.Println("STDOUT: " + stdout + "\nSTDERR: " + stderr)

		// TODO Update Job
		SetJobStatus(job)
	}
}

func GetJob() (job Job, err error) {
	resp, err := http.Get("http://localhost:8081/GetJob/")
	if err != nil {
		return job, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return job, err
	}

	json.Unmarshal(bodyBytes, &job)
	return job, nil
}

func ExecCommand(cmd_str string) (cmd_stdout string, cmd_stderr string, err error) {
	cmd := exec.Command("sh", "-c", cmd_str)
	stderr, err := cmd.StderrPipe()
	stdout, err := cmd.StdoutPipe()

	if err != nil {
		return cmd_stdout, cmd_stderr, err
	}

	if err := cmd.Start(); err != nil {
		return cmd_stdout, cmd_stderr, err
	}

	cmd_stdout_bytes, _ := ioutil.ReadAll(stdout)
	cmd_stdout = string(cmd_stdout_bytes)

	cmd_stderr_bytes, _ := ioutil.ReadAll(stderr)
	cmd_stderr = string(cmd_stderr_bytes)
	if err = cmd.Wait(); err != nil {
		return cmd_stdout, cmd_stderr, err
	}
	return cmd_stdout, cmd_stderr, nil
}

//TODO return only error
func SetJobStatus(job Job) error {
	buf := new(bytes.Buffer)
	json.NewEncoder(buf).Encode(job)
	resp, err := http.Post("http://localhost:8081/SetJobStatus/", "text/json", buf)

	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}
