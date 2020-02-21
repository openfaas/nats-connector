FROM golang:1.13 as build

ENV GO111MODULE=off
ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/openfaas-incubator/nats-connector

COPY vendor     vendor
COPY config	    config
COPY nats	    nats
COPY main.go    .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" \
 && go test -v ./... \
 && CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w" -installsuffix cgo -o /usr/bin/connector

FROM alpine:3.11 as ship
RUN apk add --no-cache ca-certificates

COPY --from=build /usr/bin/connector /usr/bin/connector
WORKDIR /root/

CMD ["/usr/bin/connector"]
