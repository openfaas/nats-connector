/// Copyright (c) OpenFaaS Author(s) 2020. All rights reserved.
/// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package nats

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	nats "github.com/nats-io/nats.go"
	"github.com/openfaas/connector-sdk/types"
)

const queueGroup = "openfaas_nats_worker_group"

const clientName = "openfaas_connector"

// BrokerConfig high level config for the broker
type BrokerConfig struct {

	// Host is the NATS address, the port is hard-coded to 4222
	Host string

	// ConnTimeout is the timeout for Dial on a connection.
	ConnTimeout time.Duration
}

// Broker used to subscribe to NATS subjects
type Broker interface {
	Subscribe(types.Controller, []string) error
}

type broker struct {
	client *nats.Conn
}

// NATSPort hard-coded port for NATS
const NATSPort = "4222"

// NewBroker loops until we are able to connect to the NATS server
func NewBroker(config BrokerConfig) (Broker, error) {
	broker := &broker{}
	brokerURL := fmt.Sprintf("nats://%s:%s", config.Host, NATSPort)

	for {
		client, err := nats.Connect(brokerURL,
			nats.Timeout(config.ConnTimeout),
			nats.Name(clientName))

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

	return broker, nil
}

// Subscribe to a list of NATS subjects and block until interrupted
func (b *broker) Subscribe(controller types.Controller, topics []string) error {
	log.Printf("Configured topics: %v", topics)

	if b.client == nil {
		return fmt.Errorf("client was nil, try to reconnect")
	}

	wg := sync.WaitGroup{}
	wg.Add(1)

	subs := []*nats.Subscription{}
	for _, topic := range topics {
		log.Printf("Binding to topic: %q", topic)

		sub, err := b.client.QueueSubscribe(topic, queueGroup, func(m *nats.Msg) {
			log.Printf("Topic: %s, message: %q", m.Subject, string(m.Data))
			controller.Invoke(m.Subject, &m.Data, http.Header{
				"X-Topic": []string{m.Subject},
			})
		})
		subs = append(subs, sub)

		if err != nil {
			log.Printf("Unable to bind to topic: %s", topic)
		}
	}

	for _, sub := range subs {
		log.Printf("Subscription: %s ready", sub.Subject)
	}

	// interrupt handling
	wg.Wait()

	b.client.Close()

	return nil
}
