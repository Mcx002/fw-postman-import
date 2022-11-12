package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"os"
)

type ResponseBody struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"body"`
}

type Version struct {
	Version string `json:"version"`
	AppName string `json:"appName"`
}

func Res(w http.ResponseWriter, data any) {
	out, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	io.WriteString(w, string(out))
}

func getStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Print("got / request\n")

	data := &ResponseBody{
		Status:  200,
		Message: "OK",
		Data: &Version{
			Version: "1.0.0",
			AppName: "Forward Postman Collection",
		},
	}

	Res(w, data)
}

func getCollection(w http.ResponseWriter, r *http.Request) {
	var err error
	var client = &http.Client{}
	var data any

	vars := mux.Vars(r)

	request, err := http.NewRequest("GET", "https://api.getpostman.com/collections/"+vars["id"], nil)
	request.Header.Set("X-API-Key", "PMAK-636f228f76ebe938af3f6a73-04dcbaa3566360d18fb3d1b58f8acc8b3c")
	if err != nil {
		panic(err)
	}

	response, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	defer response.Body.Close()

	if response.StatusCode == 401 {
		Res(w, &ResponseBody{
			Status:  401,
			Message: "Unauthorized",
			Data:    nil,
		})
	}

	err = json.NewDecoder(response.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	Res(w, data)
}

func main() {
	r := mux.NewRouter()
	//http.HandleFunc("/", getStatus)
	r.HandleFunc("/collection/{id}", getCollection)

	err := http.ListenAndServe(":3333", r)
	if errors.Is(err, http.ErrServerClosed) {
		fmt.Printf("server closed\n")
	} else if err != nil {
		fmt.Printf("error starting server: %s\n", err)
		os.Exit(1)
	}
}
