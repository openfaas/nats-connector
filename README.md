# nats-connector

[![Docker Build Status](https://github.com/openfaas-incubator/nats-connector/workflows/docker/badge.svg)](https://github.com/openfaas-incubator/nats-connector/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/openfaas-incubator/nats-connector)](https://goreportcard.com/report/github.com/openfaas-incubator/nats-connector)
[![GoDoc](https://godoc.org/github.com/openfaas-incubator/nats-connector?status.svg)](https://godoc.org/github.com/openfaas-incubator/nats-connector)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![OpenFaaS](https://img.shields.io/badge/openfaas-serverless-blue.svg)](https://www.openfaas.com)

An OpenFaaS event-connector to trigger functions from NATS.

## Try it out

Test functions have been created to help you verify the installation of `nats-connector`.  See [`contrib/test-functions`](./contrib/test-functions).

### Deploy on Kubernetes

The following instructions show how to run and test `nats-connector` on Kubernetes.

1. Deploy the receiver functions, the receiver function must have the `topic` annotation:

   ```bash
   export OPENFAAS_URL="http://localhost:8080" # Set your gateway via env variable or the -g flag
   faas-cli deploy --name receive-message --image openfaas/nats-receive-message:latest --fprocess='./handler' --annotation topic="nats-test"
   ```

   Or deploy with the `stack.yml` provided in this repo:
   ```
   cd contrib/test-functions
   faas-cli template pull stack
   faas-cli deploy --filter receive-message
   ```

2. Deploy the connector with:

   ```bash
   kubectl apply -f ./yaml/kubernetes/connector-dep.yaml
   ```

3. Deploy the `publish-message` function

   ```bash
   faas-cli deploy --name publish-message --image openfaas/nats-publish-message:latest --fprocess='./handler' --env nats_url=nats://nats.openfaas:4222
   ```

   Or deploy via `stack.yml`

      ```
   cd contrib/test-functions
   faas-cli template pull stack
   faas-cli deploy --filter publish-message
   ```

4. Now publish a message on the `nats-test` topic. 

   Invoke the publisher
   ```bash
   faas-cli invoke publish-message <<< "test message"
   ```

4. Verify that the receiver was invoked by checking the logs

   ```bash
   faas-cli logs receive-message

   2019-12-29T19:06:50Z 2019/12/29 19:06:50 received "test message"
   ```

## Building

Build and release is done via CI, but you can also build your own version locally.

```bash
export TAG=0.2.1
make build push
```

### Configuration

Configuration is by environment variable, which can be set in the Kubernetes YAML file: [yaml/kubernetes/connector-dep.yaml](./yaml/kubernetes/connector-dep.yaml)

| Variable             | Description                   |  Default                                        |
| -------------------- | ------------------------------|--------------------------------------------------|
| `topics`             | Delimited list of topics    |  `nats-test,`                                   |
| `broker_host`        | The host, but not the port for NATS | `nats` |
| `async-invocation`   | Queue the invocation with the built-in OpenFaaS queue-worker and return immediately    |  `false` |
| `gateway_url`        | The URL for the OpenFaaS gateway | `http://gateway:8080` |
| `upstream_timeout`   | Timeout to wait for synchronous invocations | `60s` |
| `rebuild_interval`   | Interval at which to rebuild the map of topics <> functions | `5s`  |
| `topic_delimiter`    | Used to separate items in `topics` variable | `,` |

