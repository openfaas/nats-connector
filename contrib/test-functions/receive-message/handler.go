package function

import (
	"fmt"
	"log"
	"net/http"

	handler "github.com/openfaas-incubator/go-function-sdk"
)

// Handle a serverless request
func Handle(req handler.Request) (handler.Response, error) {
	log.Printf("received %q", string(req.Body))

	return handler.Response{
		Body:       []byte(fmt.Sprintf("Received msg: %q", string(req.Body))),
		StatusCode: http.StatusOK,
	}, nil
}
