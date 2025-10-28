package main

import (
	"QuestionGame/repository/mysql"
	"QuestionGame/service/userservice"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const signKey = "secret_sign_key"

func main() {
	mysqlRepo := mysql.New()
	userservice.New(mysqlRepo, signKey)

	mux := http.NewServeMux()
	mux.HandleFunc("/users/register", registerHandler)
	mux.HandleFunc("/users/login", userLoginHandler)
	mux.HandleFunc("/users/profile", userProfileHandler)

	log.Println("Server is listening on port 8080 ...")

	http.ListenAndServe(":8080", mux)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"detail": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"detail": "Failed to read data from body"}`, http.StatusBadRequest)
		return
	}

	var req userservice.RegisterRequest
	err = json.Unmarshal(data, &req)
	if err != nil {
		fmt.Fprintf(w, `{"detail": "Failed to unmarshal data to json: %s"}`, err.Error())
		return
	}

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, signKey)

	_, err = userSvc.Register(req)
	if err != nil {
		fmt.Fprintf(w, `{"detail": "%s"}`, err.Error())
		return
	}

	w.Write([]byte(`{"detail": "User created successfully"}`))
}

func userLoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"detail": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, `{"detail": "Failed to read data from body"}`, http.StatusBadRequest)
		return
	}

	var lReq userservice.LoginRequest
	err = json.Unmarshal(data, &lReq)
	if err != nil {
		fmt.Fprintf(w, `{"detail": "Failed to unmarshal data to json: %s"}`, err.Error())
		return
	}

	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, signKey)

	response, err := userSvc.Login(lReq)
	if err != nil {
		fmt.Fprintf(w, `{"detail": "%s"}`, err.Error())
		return
	}

	data, mErr := json.Marshal(response)
	if mErr != nil {
		fmt.Fprintf(w, `{"detail": "%s"}`, mErr.Error())
		return
	}

	w.Write(data)
}

func userProfileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"detail": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	profileReq := userservice.ProfileRequest{UserID: 0}
	mysqlRepo := mysql.New()
	userSvc := userservice.New(mysqlRepo, signKey)

	resp, err := userSvc.Profile(profileReq)
	if err != nil {
		fmt.Fprintf(w, `{"detail": "%s"}`, err.Error())
		return
	}

	data, err := json.Marshal(resp)
	if err != nil {
		fmt.Fprintf(w, `{"detail": "%s"}`, err.Error())
		return
	}

	w.Write(data)
}
