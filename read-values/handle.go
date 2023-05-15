package function

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"strconv"

	dapr "github.com/dapr/go-sdk/client"
)

var STATE_STORE_NAME = "statestore"
var daprClient dapr.Client
var DAPR_HOST = "my-ambient-dapr-ambient.default.svc.cluster.local"
var DAPR_PORT = "50001"

type MyValues struct {
	Values []string
}

// Handle an HTTP Request.
func Handle(ctx context.Context, res http.ResponseWriter, req *http.Request) {

	daprClient, err := dapr.NewClientWithAddress(net.JoinHostPort(DAPR_HOST, DAPR_PORT))
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
	avg = float64(total / count)
	respondWithJSON(res, http.StatusOK, avg)

}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
