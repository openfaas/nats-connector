FROM golang:1.13 as build

RUN mkdir -p /go/src/github.com/openfaas-incubator/nats-connector
WORKDIR /go/src/github.com/openfaas-incubator/nats-connector

COPY vendor     vendor
COPY config	    config
COPY nats	    nats
COPY main.go    .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))"

RUN go test -v ./...

# Stripping via -ldflags "-s -w" 
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o /usr/bin/producer

FROM alpine:3.9 as ship
RUN apk add --no-cache ca-certificates

COPY --from=build /usr/bin/producer /usr/bin/producer
WORKDIR /root/

CMD ["/usr/bin/producer"]
