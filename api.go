package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "os"
	// "strconv"

	"github.com/gorilla/mux"
)


//2 handlers
type APIServer struct {
	listenAddr string
}

func NewAPIServer (listenAddr string) *APIServer {
	return &APIServer {
		listenAddr: listenAddr,
	}
}


func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/JSON")

	return json.NewEncoder(w).Encode(v)
}


type apiFunc func(http.ResponseWriter, *http.Request) error


type ApiError struct {
	Error string `json:"error"`
}

func makeHTTPHandleFunc (f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//handle error
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}


//start server up
func (s *APIServer) Run() {
	// from gorilla/mux
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))

	log.Println("JSON API server running on port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}


//account handlers
func (s *APIServer) handleAccount (w http.ResponseWriter, r *http.Request) error {
	
	if r.Method == "GET" {
		return s.handleGetAccount(w, r)
	}

	if r.Method == "POST" {
		return s.handleCreateAccount(w, r)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("Method not allowed %s", r.Method)
}

func (s *APIServer) handleGetAccount (w http.ResponseWriter, r *http.Request) error {
	account := NewAccount("Hakim", "Chulan")
	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleCreateAccount (w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteAccount (w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer (w http.ResponseWriter, r *http.Request) error {
	return nil
}

