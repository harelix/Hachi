package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/messages"
	log "github.com/sirupsen/logrus"
	"time"
)

func IncomingMessageHandler(message *nats.Msg) {

	//todo: add error handling here!!!!
	//invocationTimeout := config.New().Service.Agent.agent.InvocationTimeout()
	invocationTimeout := config.New().IAM.GetInvocationTimeout()
	//in Milliseconds
	timeout := time.Duration(invocationTimeout)

	if invocationTimeout != -1 {
		invocationTimeout = HachiContext.ContextTimeoutMax
	}

	capsule := messages.Capsule{}
	err := json.Unmarshal(message.Data, &capsule)
	if err != nil {
		log.Error("Incoming message Unmarshal error " + string(message.Data))
	}

	//RBAC - Enforce / Check roles
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	context.WithValue(ctx, HachiContext.ContextCapsule, capsule)
	//context.WithValue(ctx,HachiContext.NATSDefaultAfferentSubjects, message.Subject)

	defer cancel()

	invocationCh := make(chan messages.Capsule)

	go ProcessIncomingCapsule(ctx, invocationCh, capsule)

	capsuleUpdates := <-invocationCh
	fmt.Println(capsuleUpdates)

	select {
	case a := <-invocationCh:
		fmt.Println("Finished long running task.", a)
	case <-ctx.Done():
		fmt.Println("Timed out context before task is finished.")
	}
}

func ProcessIncomingCapsule(ctx context.Context, ch chan messages.Capsule, capsule messages.Capsule) {

	printMsg(capsule)

	ch <- messages.Capsule{
		Message: "OK",
		Headers: nil,
		Subject: nil,
		Route:   nil,
	}

	ch <- messages.Capsule{
		Message: "AOK",
		Headers: nil,
		Subject: nil,
		Route:   nil,
	}
}

func printMsg(capsule messages.Capsule) {
	log.Printf("Received on [%s]: '%s', '%s'", capsule.Subject, string(capsule.Route.Name), capsule.Message)
}

// SubscribeToSubjects - register controller and agent to corresponding subjects
func GetSubscriptionSubjects(config *config.HachiConfig) []string {
	dna := config.Service.DNA
	if dna.Controller.Enabled {
		//heartbeat edge devices registration
		return dna.Controller.Identifiers
	} else {
		return dna.Agent.Identifiers
	}
}
