package agent

import (
	"context"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/messages"
	"github.com/rills-ai/Hachi/pkg/messaging"
)

func DispatchCapsuleToMessageQueueSubscribers(ctx context.Context, capsule messages.Capsule) (messages.DefaultResponseMessage, error) {

	err := messaging.Get().Publish(ctx, capsule, []string{HachiContext.AgentsNATSDefaultAfferentSubjects})
	if err != nil {
		return messages.DefaultResponseMessage{
			Error: true,
		}, err
	}
	return messages.DefaultResponseMessage{}, nil
}
