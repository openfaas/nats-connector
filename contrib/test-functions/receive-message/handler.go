package function

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	handler "github.com/openfaas-incubator/go-function-sdk"
)

// Handle a serverless request
func Handle(req handler.Request) (handler.Response, error) {
	log.Printf("Received: %q", string(req.Body))

	if val, ok := os.LookupEnv("wait"); ok && len(val) > 0 {
		parsedVal, _ := time.ParseDuration(val)
		log.Printf("Waiting for %s before returning", parsedVal.String())
		time.Sleep(parsedVal)
	}

	return handler.Response{
		Body:       []byte(fmt.Sprintf("Received: %q", string(req.Body))),
		StatusCode: http.StatusOK,
	}, nil
}
