// Copyright (c) OpenFaaS Author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/openfaas-incubator/connector-sdk/types"
	"github.com/openfaas-incubator/nats-connector/config"
	"github.com/openfaas-incubator/nats-connector/nats"
)

func main() {
	creds := types.GetCredentials()
	config := config.Get()

	controllerConfig := &types.ControllerConfig{
		UpstreamTimeout:          config.UpstreamTimeout,
		GatewayURL:               config.GatewayURL,
		RebuildInterval:          config.RebuildInterval,
		PrintResponse:            config.PrintResponse,
		PrintResponseBody:        config.PrintResponseBody,
		TopicAnnotationDelimiter: config.TopicAnnotationDelimiter,
		AsyncFunctionInvocation:  config.AsyncFunctionInvocation,
		PrintSync:                config.PrintSync,
	}

	controller := types.NewController(creds, controllerConfig)
	controller.BeginMapBuilder()

	brokerConfig := nats.BrokerConfig{
		Host:        config.Broker,
		ConnTimeout: config.UpstreamTimeout, // ConnTimeout isn't the same as UpstreamTimeout, it's just the delay to connect to NATS.
		Credentials: config.Credentials,
	}

	fmt.Printf(`==============================================================================
NATS Connector for OpenFaaS

Gateway URL: %s
NATS URL: nats://%s:4222
Credentials: %s
Topics: %s
Upstream timeout: %s
Topic-map rebuild interval: %s
Async invocation: %q
==============================================================================

`, config.GatewayURL,
		config.Broker,
		config.Credentials,
		config.UpstreamTimeout,
		config.Topics,
		config.RebuildInterval,
		strconv.FormatBool(config.AsyncFunctionInvocation),
	)

	broker, err := nats.NewBroker(brokerConfig)

	if err != nil {
		log.Fatal(err)
	}

	err = broker.Subscribe(controller, config.Topics)
	if err != nil {
		log.Fatal(err)
	}
}
