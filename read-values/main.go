package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	dapr "github.com/dapr/go-sdk/client"
	"github.com/go-chi/chi"
)

var (
	daprClient       dapr.Client
	STATE_STORE_NAME = GetenvOrDefault("STATE_STORE_NAME", "statestore")
	DAPR_HOST        = GetenvOrDefault("DAPR_HOST", "my-ambient.default.svc.cluster.local")
	DAPR_PORT        = GetenvOrDefault("DAPR_PORT", "50001")
)

type MyValues struct {
	Values []string
}

func main() {
	r := chi.NewRouter()
	r.Get("/", Handle)
	r.Get("/health/{endpoint:readiness|liveness}", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})
	port := GetenvOrDefault("APP_PORT", "8080")
	log.Printf("Starting Read Values App in Port: %s", port)
	http.ListenAndServe(":"+port, r)
}

// Handle an HTTP Request.
func Handle(res http.ResponseWriter, req *http.Request) {

	ctx := context.Background()

	daprClient, err := dapr.NewClientWithAddress(fmt.Sprintf("%s:%s", DAPR_HOST, DAPR_PORT))
	if err != nil {
		panic(err)
	}

	result, err := daprClient.GetState(ctx, STATE_STORE_NAME, "values", nil)
	if err != nil {
		panic(err)
	}
	myValues := MyValues{}
	json.Unmarshal(result.Value, &myValues)

	var total int
	var count int
	for _, v := range myValues.Values {
		intVar, _ := strconv.Atoi(v)
		total += intVar
		count++
	}

	var avg float64
	if count == 0 {
		avg = float64(0)
	} else {
		avg = float64(total / count)
	}
	respondWithJSON(res, http.StatusOK, avg)
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
