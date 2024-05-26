package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	userQuotas map[UserID]Quota
	stopC      chan bool
}

func newServer(stopC chan bool) Server {
	return Server{userQuotas: bootstrapUserQuotas(), stopC: stopC}
}

// Run to bootstrap http server.
func (s *Server) Run() {
	router := mux.NewRouter()
	router.Path("/kv").HandlerFunc(s.PostData).Methods("POST")

	runningIpPort := "127.0.0.1:8000"
	srv := &http.Server{
		Addr:         runningIpPort,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
		Handler:      router,
	}
	fmt.Println("Listening at :", runningIpPort)
	log.Fatal(srv.ListenAndServe())
	s.stopC <- true
}

func (s *Server) PostData(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var input Input
	err := s.decode(r, &input)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(fmt.Sprintf(`{"error": "server error", "text": "%v"}`, err.Error())))
		return
	}
	err = s.validateInput(input)
	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(fmt.Sprintf(`{"error": "validation error", "text": "%v"}`, err.Error())))
		return
	}
	// Do the work, check for limiter, then insert in store.
	quota, exists := s.userQuotas[input.UserID]
	if !exists {
		w.WriteHeader(403)
		w.Write([]byte(`{"error": "validation error", "text": "user does not have any quota"}`))
		return
	}
	err = input.checkRate(quota)
	if err != nil {
		w.WriteHeader(403)
		w.Write([]byte(fmt.Sprintf(`{"error": "quota exceeded", "text": "%v"}`, err.Error())))
		return
	}
	err = input.Accept()
	if err != nil {
		w.WriteHeader(403)
		w.Write([]byte(fmt.Sprintf(`{"error": "storage error", "text": "%v"}`, err.Error())))
		return
	}
	w.WriteHeader(201)
	w.Write([]byte(`{"ok": true, "text": "data saved successfully"}`))
}

func (s Server) decode(r *http.Request, output any) error {
	decoder := json.NewDecoder(r.Body)
	return decoder.Decode(output)
}

func (s Server) validateInput(input Input) error {
	// We assume that 0 is nil here.
	if input.UniqueID == 0 || input.UserID == 0 {
		return errors.New("unique_id and user_id must be present and valid")
	}
	return nil
}
