package function

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"

	nats "github.com/nats-io/nats.go"
	handler "github.com/openfaas-incubator/go-function-sdk"
)

var (
	subject        = "nats-test"
	defaultMessage = "Hello World"
)

// Handle a serverless request
func Handle(req handler.Request) (handler.Response, error) {
	msg := defaultMessage
	if len(req.Body) > 0 {
		msg = string(bytes.TrimSpace(req.Body))
	}

	natsURL := nats.DefaultURL
	val, ok := os.LookupEnv("nats_url")
	if ok {
		natsURL = val
	}

	nc, err := nats.Connect(natsURL)
	if err != nil {
		errMsg := fmt.Sprintf("can not connect to nats: %s", err)
		log.Printf(errMsg)
		r := handler.Response{
			Body:       []byte(errMsg),
			StatusCode: http.StatusInternalServerError,
		}
		return r, err
	}
	defer nc.Close()

	log.Printf("Publishing %d bytes to: %q\n", len(msg), subject)

	err = nc.Publish(subject, []byte(msg))
	if err != nil {
		log.Println(err)

		r := handler.Response{
			Body:       []byte(fmt.Sprintf("can not publish to nats: %s", err)),
			StatusCode: http.StatusInternalServerError,
		}
		return r, err
	}

	return handler.Response{
		Body:       []byte(fmt.Sprintf("Published %d bytes to: %q", len(msg), subject)),
		StatusCode: http.StatusOK,
	}, nil
}
