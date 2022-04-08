FROM --platform=${BUILDPLATFORM:-linux/amd64} golang:1.17 as build

ARG TARGETPLATFORM
ARG BUILDPLATFORM
ARG TARGETOS
ARG TARGETARCH

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/openfaas/nats-connector

COPY go.mod	.
COPY go.sum	.
COPY config	    config
COPY nats	    nats
COPY main.go    .

# Run a gofmt and exclude all vendored code.
RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))"

RUN go test -v ./...
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build -a -ldflags "-s -w" -installsuffix cgo -o /usr/bin/connector

FROM --platform=${TARGETPLATFORM:-linux/amd64} alpine:3.15 as ship
RUN apk add --no-cache ca-certificates

COPY --from=build /usr/bin/connector /usr/bin/connector
WORKDIR /root/

CMD ["/usr/bin/connector"]
