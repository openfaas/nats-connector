// Copyright (c) OpenFaaS Project 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.
package main

import (
	"github.com/openfaas-incubator/connector-sdk/types"
	"github.com/openfaas-incubator/nats-connector/config"
	"github.com/openfaas-incubator/nats-connector/nats"
)

func main() {
	creds := types.GetCredentials()
	config := config.Get()

	controllerConfig := &types.ControllerConfig{
		RebuildInterval: config.RebuildInterval,
		GatewayURL:      config.GatewayURL,
		PrintResponse:   config.PrintResponse,
	}

	controller := types.NewController(creds, controllerConfig)
	controller.BeginMapBuilder()

	brokerConfig := nats.BrokerConfig{
		Host:        config.Broker,
		ConnTimeout: config.UpstreamTimeout,
	}

	broker := nats.NewBroker(brokerConfig)
	broker.Subscribe(controller, config.Topics)
}
