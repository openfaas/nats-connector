# nats-connector

[![Go Report Card](https://goreportcard.com/badge/github.com/openfaas-incubator/nats-connector)](https://goreportcard.com/report/github.com/openfaas-incubator/nats-connector) [![Build
Status](https://travis-ci.org/openfaas-incubator/nats-connector.svg?branch=master)](https://travis-ci.org/openfaas-incubator/nats-connector) [![GoDoc](https://godoc.org/github.com/openfaas-incubator/nats-connector?status.svg)](https://godoc.org/github.com/openfaas-incubator/nats-connector) [![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![OpenFaaS](https://img.shields.io/badge/openfaas-serverless-blue.svg)](https://www.openfaas.com)

The NATS connector connects OpenFaaS functions to NATS topics.

## Building

```
export TAG=0.2.0
make build push
```

## Try it out

### Deploy on Kubernetes

The following instructions show how to run `kafka-connector` on Kubernetes.

Deploy a function with a `topic` annotation:

```bash
$ faas store deploy figlet --annotation topic="faas-request" --gateway <faas-netes-gateway-url>
```

Deploy the connector with:

```bash
kubectl apply -f ./yaml/kubernetes/connector-dep.yml
```

Now publish a message on the faas-request topic.
