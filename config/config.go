// Copyright (c) OpenFaaS Project 2018. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package config

import (
	"log"
	"os"
	"strings"
	"time"
)

// Config for the NATS Connector
type Config struct {
	Broker string
	Topics []string

	GatewayURL               string
	UpstreamTimeout          time.Duration
	RebuildInterval          time.Duration
	PrintResponse            bool
	PrintResponseBody        bool
	PrintSync                bool
	AsyncFunctionInvocation  bool
	TopicAnnotationDelimiter string
}

// Get will load the NATS Connector config from the environment variables
func Get() Config {
	broker := "nats"
	if val, exists := os.LookupEnv("broker_host"); exists {
		broker = val
	}

	topics := []string{}
	if val, exists := os.LookupEnv("topics"); exists {
		for _, topic := range strings.Split(val, ",") {
			if len(topic) > 0 {
				topics = append(topics, topic)
			}
		}
	}
	if len(topics) == 0 {
		log.Fatal(`Provide a list of topics i.e. topics="payment_published,slack_joined"`)
	}

	gatewayURL := "http://gateway:8080"
	if val, exists := os.LookupEnv("gateway_url"); exists {
		gatewayURL = val
	}

	upstreamTimeout := time.Second * 30
	rebuildInterval := time.Second * 3

	if val, exists := os.LookupEnv("upstream_timeout"); exists {
		parsedVal, err := time.ParseDuration(val)
		if err == nil {
			upstreamTimeout = parsedVal
		}
	}

	if val, exists := os.LookupEnv("rebuild_interval"); exists {
		parsedVal, err := time.ParseDuration(val)
		if err == nil {
			rebuildInterval = parsedVal
		}
	}

	printResponse := false
	if val, exists := os.LookupEnv("print_response"); exists {
		printResponse = (val == "1" || val == "true")
	}

	printResponseBody := false
	if val, exists := os.LookupEnv("print_response_body"); exists {
		printResponseBody = (val == "1" || val == "true")
	}

	printSync := false
	if val, exists := os.LookupEnv("print_sync"); exists {
		printSync = (val == "1" || val == "true")
	}

	asyncFunctionInvocation := true
	if val, exists := os.LookupEnv("asynchronous_invocation"); exists {
		asyncFunctionInvocation = (val == "1" || val == "true")
	}

	delimiter := ","
	if val, exists := os.LookupEnv("topic_delimiter"); exists {
		if len(val) > 0 {
			delimiter = val
		}
	}

	return Config{
		Broker:                   broker,
		Topics:                   topics,
		GatewayURL:               gatewayURL,
		UpstreamTimeout:          upstreamTimeout,
		RebuildInterval:          rebuildInterval,
		PrintResponse:            printResponse,
		PrintResponseBody:        printResponseBody,
		PrintSync:                printSync,
		AsyncFunctionInvocation:  asyncFunctionInvocation,
		TopicAnnotationDelimiter: delimiter,
	}
}
