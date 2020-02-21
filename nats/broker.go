/// Copyright (c) OpenFaaS Author(s) 2020. All rights reserved.
/// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package nats

import (
	"log"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/openfaas-incubator/connector-sdk/types"
)

const queueGroup = "openfaas_nats_worker_group"
const clientName = "openfaas_connector"

// BrokerConfig high level config for the broker
type BrokerConfig struct {
	Host        string
	ConnTimeout time.Duration
}

// Broker used to subscribe to NATS subjects
type Broker interface {
	Subscribe(types.Controller, []string)
}

type broker struct {
	client *nats.Conn
}

// NewBroker loops until we are able to connect to the NATS server
func NewBroker(config BrokerConfig) Broker {
	broker := &broker{}

	brokerURL := "nats://" + config.Host + ":4222"
	for {
		client, err := nats.Connect(brokerURL, nats.Timeout(config.ConnTimeout), nats.Name(clientName))
		if client != nil && err == nil {
			broker.client = client
			break
		}

		if client != nil {
			client.Close()
		}
		log.Println("Wait for brokers to come up.. ", brokerURL)
		time.Sleep(1 * time.Second)
		// TODO Add healthcheck
	}
	return broker
}

// Subscribe to a list of NATS subjects and block until interrupted
func (b *broker) Subscribe(controller types.Controller, topics []string) {
	log.Printf("Configured topics: %v", topics)

	wg := sync.WaitGroup{}
	wg.Add(1)

	for _, topic := range topics {
		log.Printf("Binding to topic: %v", topic)
		// check client not nil
		b.client.QueueSubscribe(topic, queueGroup, func(m *nats.Msg) {
			log.Printf("Received topic: %s, message: %s", m.Subject, string(m.Data))
			controller.Invoke(m.Subject, &m.Data)
		})
	}

	// interrupt handling
	wg.Wait()
	b.client.Close()
}
