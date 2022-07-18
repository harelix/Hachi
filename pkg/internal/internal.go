package internal

import (
	"context"
	"encoding/base64"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/agent"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/controller"
	"github.com/rills-ai/Hachi/pkg/cryptography"
	"github.com/rills-ai/Hachi/pkg/messages"
	"github.com/rills-ai/Hachi/pkg/messaging"
	"github.com/rills-ai/Hachi/pkg/webhooks"
	log "github.com/sirupsen/logrus"
	"strings"
)

type InternalResponse struct {
	Result    string
	Directive string
}

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

[A/C]
->API->Handler->Dispatch->NATS->Subscriber->Exec(Sinks)->Response->NATS->Subscriber->Confirmation


func DispatchCapsuleToMessageQueueSubscribers(ctx context.Context, capsule messages.Capsule) (messages.DefaultResponseMessage, error ) {

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


func internals(capsule messages.Capsule, directive string) InternalResponse {
	args := strings.Split(directive, "#")

	switch args[0] {
	case HachiContext.InternalsCryptoEncrypt:

		message := cryptography.Encryption(capsule.Message)

		return InternalResponse{
			Result:    base64.StdEncoding.EncodeToString([]byte(message)),
			Directive: directive,
		}
	case HachiContext.InternalsCryptoDecrypt:
		decodedMessage, err := base64.StdEncoding.DecodeString(capsule.Message)
		if err != nil {
			//todo: report this one
		}
		message := cryptography.Decryption(string(decodedMessage))
		return InternalResponse{
			Result:    message,
			Directive: directive,
		}
	default:
		return InternalResponse{
			Result: "NoActTak",
		}
	}
}
