package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	jwt "github.com/golang-jwt/jwt/v4"
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
	router.HandleFunc("/account/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetAccountByID)))
	router.HandleFunc("/transfer", makeHTTPHandleFunc(s.handleTransfer))

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

	//moved into handleGetAccountByID
	// if r.Method == "DELETE" {
	// 	return s.handleDeleteAccount(w, r)
	// }

	return fmt.Errorf("method not allowed %s", r.Method)
}

// GET /account
func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error { 
	accounts, err := s.store.GetAccounts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, accounts)
}


// GET: /account/id
func (s *APIServer) handleGetAccountByID (w http.ResponseWriter, r *http.Request) error {
// 	// account := NewAccount("Hakim", "Chulan")
// 	// return WriteJSON(w, http.StatusOK, account)

	if r.Method == "GET" {
		id, err := getID(r);

		if err != nil {
			return err
		}

		account, err := s.store.GetAccountByID(id)

		if err != nil {
			return err
		}

		return WriteJSON(w, http.StatusOK, account)
	}

	if r.Method == "DELETE" {
		return s.handleDeleteAccount(w, r)
	}

	return fmt.Errorf("method not allowed %s", r.Method)
}

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

	id, err := getID(r);

	if err != nil {
		return err
	}

	if err := s.store.DeleteAccount(id); err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, map[string]int{"deleted": id})
}

func (s *APIServer) handleTransfer (w http.ResponseWriter, r *http.Request) error {

	transferReq := new(TransferRequest)

	if err := json.NewDecoder(r.Body).Decode(transferReq); err != nil {
		return err
	}
	defer r.Body.Close()


	return WriteJSON(w, http.StatusOK, transferReq)
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

func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("calling AWT auth middleware")

		tokenString := r.Header.Get("x-jwt-token")

		_, err := validateJWT(tokenString)

		if err != nil {
			WriteJSON(w, http.StatusForbidden, ApiError{Error: "invalid token"})
			return 
		}

		handlerFunc(w, r)
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {

	secret := os.Getenv("JWT_SECRET")

	// from https://pkg.go.dev/github.com/golang-jwt/jwt/v4#example-Parse-Hmac
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
	
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	
}

func getID(r *http.Request) (int, error) {
	idStr := mux.Vars(r)["id"]

	id, err := strconv.Atoi(idStr)

	if err != nil {
		return id, fmt.Errorf("invalid id give %s", idStr)
	}

	return id, nil
}