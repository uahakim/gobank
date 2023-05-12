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

//start server up

type APIServer struct {
	listenAddr string
	store Storage
}

func NewAPIServer (listenAddr string, store Storage) *APIServer {
	return &APIServer {
		listenAddr: listenAddr,
		store: store,
	}
}



func (s *APIServer) Run() {
	// from gorilla/mux
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPHandleFunc(s.handleAccount))
	// router.HandleFunc("/account/{id}", makeHTTPHandleFunc(s.handleGetAccountByID))

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

// GET /accounts
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error { 
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}

// func (s *APIServer) handleGetAccountByID (w http.ResponseWriter, r *http.Request) error {
// 	// account := NewAccount("Hakim", "Chulan")
// 	// return WriteJSON(w, http.StatusOK, account)

// 	vars := mux.Vars(r)["id"]
// 	fmt.Println(id)
// 	return WriteJSON(w, http.StatusOK, &Account)
// }

// POST /account then send in body firstname and lastname
func (s *APIServer) handleCreateAccount (w http.ResponseWriter, r *http.Request) error {

	createAccountReq := new(CreateAccountRequest)
	if err := json.NewDecoder(r.Body).Decode(createAccountReq); err != nil {
		return err
	}

	// the same as above but using reference &
	// createAccountReq := CreateAccountRequest{}
	// if err := json.NewDecoder(r.Body).Decode(&createAccountReq); err != nil {
	// 	return err
	// }

	account := NewAccount(createAccountReq.FirstName, createAccountReq.LastName)
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, account)
}

func (s *APIServer) handleDeleteAccount (w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleTransfer (w http.ResponseWriter, r *http.Request) error {
	return nil
}


//best practice: put least important at the bottom

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	
	w.Header().Add("Content-Type", "application/JSON")
	w.WriteHeader(status)

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

