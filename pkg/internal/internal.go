package internal

import (
	"context"
	"encoding/base64"
	"fmt"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/agent"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/controller"
	"github.com/rills-ai/Hachi/pkg/cryptography"
	"github.com/rills-ai/Hachi/pkg/db"
	"github.com/rills-ai/Hachi/pkg/internal/selectors"
	"github.com/rills-ai/Hachi/pkg/messages"
	"github.com/rills-ai/Hachi/pkg/webhooks"
	log "github.com/sirupsen/logrus"
	"strings"
)

func ProcessCapsule(ctx context.Context, capsule messages.Capsule) (messages.DefaultResponseMessage, error) {

	/*
		-== Internal actions ==-
		These methods are common internal streams/api-endpoints that Hachi exposes
		for various behavioural requirements and restrictions (security, cryptography, etc. )
	*/
	if capsule.Route.Remote.Internal != nil {
		directive := capsule.Route.Remote.Internal.Type
		response := internals(capsule, directive)

		m := messages.DefaultResponseMessage{
			Error:     false,
			Data:      response,
			Selectors: capsule.Selectors,
		}
		return m, nil
	}

	/*
		-== Webhooks handlers ==-
	*/
	if capsule.Route.Remote.Webhook != nil {
		event := capsule.Route.Remote.Webhook.Event
		response := webhooks.Exec(capsule, event)
		m := messages.DefaultResponseMessage{
			Error:     false,
			Data:      response,
			Selectors: capsule.Selectors,
		}
		return m, nil
	}

	/*
		-==  Capsule dispatching engine  ==-
	*/

	//main message/action relay/execution
	DispatchCapsuleToMessageQueueSubscribers(ctx, capsule)

	responseMessage := messages.DefaultResponseMessage{
		Data:      HachiContext.PublishSuccessful,
		Selectors: capsule.Selectors,
	}

	return responseMessage, nil
}

func DispatchCapsuleToMessageQueueSubscribers(ctx context.Context, capsule messages.Capsule) (messages.DefaultResponseMessage, error) {

	switch config.New().IAM.GetType() {
	case config.Controller:
		//selectors
		return controller.DispatchCapsuleToMessageQueueSubscribers(ctx, capsule)

	case config.Agent:
		//only controller communication for now - unless we want mesh communication
		return agent.DispatchCapsuleToMessageQueueSubscribers(ctx, capsule)
	}

	log.Trace("this should not happen, service without valid IAM configuration.")
	return messages.DefaultResponseMessage{}, nil
}

func internals(capsule messages.Capsule, directive string) messages.InternalResponse {
	args := strings.Split(directive, "#")

	switch args[0] {
	case HachiContext.RegisterAgentInternalCommand:

		registerAgent(capsule)

		return messages.InternalResponse{
			Result:    base64.StdEncoding.EncodeToString([]byte("")),
			Directive: directive,
		}

	case HachiContext.InternalsCryptoEncrypt:

		message := cryptography.Encryption(capsule.Message)

		return messages.InternalResponse{
			Result:    base64.StdEncoding.EncodeToString([]byte(message)),
			Directive: directive,
		}
	case HachiContext.InternalsCryptoDecrypt:
		decodedMessage, err := base64.StdEncoding.DecodeString(capsule.Message)
		if err != nil {
			//todo: report this one
		}
		message := cryptography.Decryption(string(decodedMessage))
		return messages.InternalResponse{
			Result:    message,
			Directive: directive,
		}
	default:
		return messages.InternalResponse{
			Result: "NoActTak",
		}
	}
}

func registerAgent(capsule messages.Capsule) {

	agent, err := config.AgentConfigFromJSON(capsule.Message)
	if err != nil {
		log.Error("agent verification unmarshling failed")
	}
	agentId := agent.Identifiers.Core

	success, err := db.GetInstance().RegisterAgent(agentId,
		selectors.BuildAgentDedicatedChannelIdentifier(agentId))

	if err != nil {
		//todo:
	}

	fmt.Println(success)
	fmt.Println(agentId)
}
