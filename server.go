package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	userQuotas map[UserID]Quota
}

func newServer(ip, port string) Server {
	return Server{userQuotas: bootstrapUserQuotas()}
}

func (s Server) run() {
	router := mux.NewRouter()
	router.Path("/kv").Queries("path", "{path}").HandlerFunc(s.PostData).Methods("POST")

	runningIpPort := "0.0.0.0:8000"
	srv := &http.Server{
		Addr:         runningIpPort,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		Handler:      router,
	}
	fmt.Println("Listening at :", runningIpPort)
	log.Fatal(srv.ListenAndServe())
	// Maybe gracefully shutdown here?
}

func (s *Server) PostData(w http.ResponseWriter, r *http.Request) {
	var input Input
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		panic(err)
	}
	// Do the work, check for limiter, then insert in store.
	quota, exists := s.userQuotas[input.UserID]
	if !exists {
		// panic
	}
	err = quota.checkMinute()
	if err != nil {
		// panic
	}
	err = quota.checkMonthly()
	if err != nil {
		// panic
	}

}
