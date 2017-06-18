package main

import (
         "net/http"
         "github.com/gorilla/mux"
         "encoding/json"
         "log"
         "time"
       )

// TODO Configuration
func StartAgentAPI () {
   log.Println("Starting Agent API...")
   router := mux.NewRouter().StrictSlash(false)
   var agentHandler AgentHandler
   router.HandleFunc("/GetJob/", agentHandler.GetJob )
   router.HandleFunc("/SetJobStatus/", agentHandler.SetJobStatus )
   srv := &http.Server{
      Handler:      router,
      Addr:         ":8081",
      WriteTimeout: 15 * time.Second,
      ReadTimeout:  15 * time.Second,
    }

    log.Fatal(srv.ListenAndServe())
}

type AgentHandler struct {
}

// TODO send agent ID, so the server knows what type of Requests the server can take
// TODO what to do with empty response
func (m *AgentHandler) GetJob(w http.ResponseWriter, r *http.Request) {
   var c Controller
   job := c.GetWaitingJob()
   b, _ := json.Marshal(job)
   w.Write([]byte(b))
}

func (m *AgentHandler) SetJobStatus(w http.ResponseWriter, r *http.Request) {

   decoder := json.NewDecoder(r.Body)
   var job Job   
   err := decoder.Decode(&job)
   if err != nil {
      log.Println( "Decoding Error: " + err.Error() ) 
      w.Write([]byte( "Wrong or missing parameters!" ) )
   }
 
   var c Controller
   c.UpdateJob( job )
}

// TODO get job from buff utility function
