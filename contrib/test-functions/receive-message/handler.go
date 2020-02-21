package function

import (
	"fmt"
	handler "github.com/openfaas-incubator/go-function-sdk"
	"log"
	"net/http"
	"os"
	"time"
)

// Handle a serverless request
func Handle(req handler.Request) (handler.Response, error) {
	log.Printf("received %q", string(req.Body))

	if val, ok := os.LookupEnv("wait"); ok && len(val) > 0 {
		parsedVal, _ := time.ParseDuration(val)
		time.Sleep(parsedVal)
	}

	return handler.Response{
		Body:       []byte(fmt.Sprintf("Received msg: %q", string(req.Body))),
		StatusCode: http.StatusOK,
	}, nil
}
