FROM golang:1.9.2
RUN mkdir -p /go/src/github.com/openfaas-incubator/nats-connector
WORKDIR /go/src/github.com/openfaas-incubator/nats-connector

COPY vendor     vendor
COPY config     config
COPY nats       nats
COPY main.go    .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))"

RUN go test -v ./...

# Stripping via -ldflags "-s -w" 
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o ./connector

CMD ["./connector"]
