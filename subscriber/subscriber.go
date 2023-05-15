package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/gorilla/mux"
)

type Result struct {
	Data string `json:"data"`
}

func notifications(w http.ResponseWriter, r *http.Request) {

	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(string(requestDump))

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var result Result
	err = json.Unmarshal(data, &result)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Subscriber received on /notifications:", string(result.Data))

	obj, err := json.Marshal(data)
	if err != nil {
		log.Fatal(err.Error())
	}
	_, err = w.Write(obj)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func printRoot(w http.ResponseWriter, r *http.Request) {
	requestDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		fmt.Println(err)
	}
	log.Println(string(requestDump))
}

func main() {
	appPort := os.Getenv("APP_PORT")
	if appPort == "" {
		appPort = "8080"
	}

	r := mux.NewRouter()

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/", printRoot).Methods("POST")

	// Dapr subscription routes orders topic to this route
	r.HandleFunc("/notifications", notifications).Methods("POST")

	// Add handlers for readiness and liveness endpoints
	r.HandleFunc("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	log.Printf("Starting Subscriber in Port: %s", appPort)
	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+appPort, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}
