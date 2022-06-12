package messaging

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gilbarco-ai/Hachi/pkg/config"
	"github.com/nats-io/nats.go"
)

// Copyright 2022-2022 The HACHI Author
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// A HACHI agent  (https://hachi.io).
// Construct controller to agents Streaming tracts

func createDefaultStreams(hn *HachiNeuron) error {
	js := hn.JS
	for _, register := range hn.DefaultNATSRegisters {
		// Check if the KEYS_NATS_MAIN_STREAM stream already exists; if not, create it.
		_, err := js.StreamInfo(register.StreamName)

		if err != nil {
			if errors.Is(err, nats.ErrStreamNotFound) {
				log.Printf("creating stream %q and subject %q",
					register.StreamName,
					register.Subject)

				_, err = js.AddStream(&nats.StreamConfig{
					Name:     register.StreamName,
					Subjects: []string{register.Subject},
				})
				if err != nil {
					return fmt.Errorf("failed to create stream: %w", err)
				}
			} else {
				return fmt.Errorf("failed to check stream status: %w", err)
			}
		}
	}
	return nil
}
func createDefaultConsumers(hn *HachiNeuron) error {

	js := hn.JS
	for _, register := range hn.DefaultNATSRegisters {
		_, err := js.ConsumerInfo(
			register.StreamName,
			register.ListeningOn.ConsumerName)

		if err != nil {
			if errors.Is(err, nats.ErrConsumerNotFound) {
				_, err = hn.JS.AddConsumer(register.StreamName, &nats.ConsumerConfig{
					DeliverPolicy: nats.DeliverNewPolicy,
					Durable:       register.ListeningOn.ConsumerName,
					Description:   config.New().Service.Type.String() + " durable consumer",
					AckPolicy:     nats.AckExplicitPolicy,
					AckWait:       0,
					MaxDeliver:    3,
					MaxWaiting:    0,
					HeadersOnly:   false,
				})
				if err != nil {
					return fmt.Errorf("failed to create consumer: %w", err)
				}
			} else {
				return fmt.Errorf("failed to verify consumer exists: %w", err)
			}
		}
	}
	return nil

}

func bindDefaultPullSubscriber(hn *HachiNeuron, ch chan *PublishedMessage) {

	js := hn.JS

	for key, register := range hn.DefaultNATSRegisters {
		agentType := config.New().Service.Type.String()
		if key == strings.ToLower(agentType) {
			continue
		}

		subscriber, err := js.PullSubscribe(register.Subject,
			register.ListeningOn.ConsumerName,
			nats.PullMaxWaiting(512))

		if err != nil {
			ch <- &PublishedMessage{
				Message: nil,
				Error:   err,
			}
		}

		for {
			messages, _ := subscriber.Fetch(hn.GetBufferSize())
			for _, msg := range messages {
				if msg != nil {
					ch <- &PublishedMessage{
						Message: msg,
						Error:   nil,
					}
				}
			}
		}
	}
}

//NATSDefaultProvisioning Constructs controller and agents Streaming tracts
func NATSDefaultProvisioning(hn *HachiNeuron) error {
	err := createDefaultStreams(hn)
	if err != nil {
		return err
	}
	err = createDefaultConsumers(hn)
	if err != nil {
		return err
	}
	return nil
}
