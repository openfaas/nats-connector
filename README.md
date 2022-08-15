# nats-connector

[![build](https://github.com/openfaas/nats-connector/actions/workflows/build.yml/badge.svg)](https://github.com/openfaas/nats-connector/actions/workflows/build.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/openfaas-incubator/nats-connector)](https://goreportcard.com/report/github.com/openfaas-incubator/nats-connector)
[![GoDoc](https://godoc.org/github.com/openfaas-incubator/nats-connector?status.svg)](https://godoc.org/github.com/openfaas-incubator/nats-connector)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![OpenFaaS](https://img.shields.io/badge/openfaas-serverless-blue.svg)](https://www.openfaas.com)

An OpenFaaS event-connector to trigger functions from NATS pub/sub topics.

This integration only supports pub/sub topics, so if you need durability, and the ability to scale out, check out [JetStream for OpenFaaS instead](https://www.openfaas.com/blog/jetstream-for-openfaas/).

## Quick start

### Deploy the connector to faasd

See the [eBook Serverless For Everyone Else](https://openfaas.gumroad.com/l/serverless-for-everyone-else) for instructions and sample YAML to add to your faasd host.

### Or deploy the connector to Kubernetes

1. Deploy the connector using arkade

   ```bash
   arkade install nats-connector
   ```

   Alternatively, see [the Helm chart](https://github.com/openfaas/faas-netes/tree/master/chart/nats-connector)

2. Deploy the two test functions

   ```bash
   git clone https://github.com/openfaas/nats-connector --depth=1
   cd nats-connector/contrib/test-functions
   ```

   Deploy the functions using `stack.yml`, see how the `receive-message` function has the `topic=nats-test` annotation? That's how a function binds itself to a particular *NATS Subject*.

   ```bash
   faas-cli template pull stack
   faas-cli deploy
   ```

3. Now publish a message on the `nats-test` topic. 

   Invoke the publisher
   ```bash
   faas-cli invoke publish-message <<< "test message"
   ```

4. Verify that the receiver was invoked by checking the logs

   ```bash
   faas-cli logs receive-message

   2019-12-29T19:06:50Z 2019/12/29 19:06:50 received "test message"
   ```

### Configuration

Configuration of the binary is by environment variable. The names vary for the values.yaml file in [the Helm chart](https://github.com/openfaas/faas-netes/tree/master/chart/nats-connector).

| Variable             | Description                   |  Default                                        |
| -------------------- | ------------------------------|--------------------------------------------------|
| `topics`             | Delimited list of topics    |  `nats-test,`                                   |
| `broker_host`        | The host, but not the port for NATS | `nats` |
| `async-invocation`   | Queue the invocation with the built-in OpenFaaS queue-worker and return immediately    |  `false` |
| `gateway_url`        | The URL for the OpenFaaS gateway | `http://gateway:8080` |
| `upstream_timeout`   | Timeout to wait for synchronous invocations | `60s` |
| `rebuild_interval`   | Interval at which to rebuild the map of topics <> functions | `5s`  |
| `topic_delimiter`    | Used to separate items in `topics` variable | `,` |
