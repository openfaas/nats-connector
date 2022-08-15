// Copyright (c) OpenFaaS Author(s) 2020. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/openfaas/connector-sdk/types"
	"github.com/openfaas/nats-connector/config"
	"github.com/openfaas/nats-connector/nats"
	"github.com/openfaas/nats-connector/version"
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
	}

	sha, ver := version.GetReleaseInfo()
	fmt.Printf(`==============================================================================
NATS Connector for OpenFaaS %s (%s)

Gateway URL: %s
NATS URL: nats://%s:4222
Topics: %s
Upstream timeout: %s
Topic-map rebuild interval: %s
Async invocation: %q
==============================================================================

`, sha,
		ver,
		config.GatewayURL,
		config.Broker,
		config.Topics,
		config.UpstreamTimeout,
		config.RebuildInterval,
		strconv.FormatBool(config.AsyncFunctionInvocation),
	)

	broker, err := nats.NewBroker(brokerConfig)
	if err != nil {
		log.Fatal(err)
	}

	if err = broker.Subscribe(controller, config.Topics); err != nil {
		log.Fatal(err)
	}
}
