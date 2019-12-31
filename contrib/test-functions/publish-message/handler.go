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
		r := handler.Response{
			Body:       []byte(fmt.Sprintf("can not connect to nats: %s", err)),
			StatusCode: http.StatusInternalServerError,
		}
		return r, err
	}
	defer nc.Close()

	log.Printf("Sending %q to %q\n", msg, subject)
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
		Body:       []byte(fmt.Sprintf("Sent %q", msg)),
		StatusCode: http.StatusOK,
	}, nil
}
