package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"os"

	"github.com/go-chi/chi"
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
	port := GetenvOrDefault("APP_PORT", "8080")

	r := chi.NewRouter()

	// Dapr subscription routes orders topic to this route
	r.Post("/", printRoot)

	// Dapr subscription routes orders topic to this route
	r.Post("/notifications", notifications)

	// Add handlers for readiness and liveness endpoints
	r.Get("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	log.Printf("Starting Subscriber in Port: %s", port)
	// Start the server; this is a blocking call
	err := http.ListenAndServe(":"+port, r)
	if err != http.ErrServerClosed {
		log.Panic(err)
	}
}

func GetenvOrDefault(envName, defaultValue string) string {
	v := os.Getenv(envName)
	if v != "" {
		return v
	}
	return defaultValue
}
