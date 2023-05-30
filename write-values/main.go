package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi"
)

var STATE_STORE_NAME = "statestore"

var PUB_SUB_NAME = "notifications-pubsub"
var PUB_SUB_TOPIC = "notifications"

var daprClient dapr.Client
var DAPR_HOST = "my-ambient-dapr-ambient.default.svc.cluster.local"
var DAPR_PORT = "50001"

type MyValues struct {
	Values []string
}

func main() {
	port := GetenvOrDefault("APP_PORT", "8080")
	r := chi.NewRouter()
	r.Post("/", Handle)
	log.Printf("Starting Write Values App in Port: %s", port)
	http.ListenAndServe(":"+port, r)
}

func Handle(res http.ResponseWriter, req *http.Request) {
	ctx := context.Background()
	daprClient, err := dapr.NewClientWithAddress(fmt.Sprintf("%s:%s", DAPR_HOST, DAPR_PORT))
	if err != nil {
		panic(err)
	}

	value := req.URL.Query().Get("value")

	result, _ := daprClient.GetState(ctx, STATE_STORE_NAME, "values", nil)
	myValues := MyValues{}
	if result.Value != nil {
		json.Unmarshal(result.Value, &myValues)
	}

	if myValues.Values == nil || len(myValues.Values) == 0 {
		myValues.Values = []string{value}
	} else {
		myValues.Values = append(myValues.Values, value)
	}

	jsonData, err := json.Marshal(myValues)

	err = daprClient.SaveState(ctx, STATE_STORE_NAME, "values", jsonData, nil)
	if err != nil {
		panic(err)
	}

	daprClient.PublishEvent(context.Background(), PUB_SUB_NAME, PUB_SUB_TOPIC, []byte(value))

	fmt.Println("Published data:", value)

	respondWithJSON(res, http.StatusOK, myValues)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func GetenvOrDefault(envName, defaultValue string) string {
	v := os.Getenv(envName)
	if v != "" {
		return v
	}
	return defaultValue
}
