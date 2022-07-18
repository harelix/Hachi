package agent

import (
	"context"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/integrity"
	"github.com/rills-ai/Hachi/pkg/messages"
	"github.com/rills-ai/Hachi/pkg/messaging"
	log "github.com/sirupsen/logrus"
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

func SelfProvision() error {
	agentCfg, err := integrity.ValidateAgentID()
	if err != nil {
		return err
	}
	agentFromCFG, _ := config.AgentConfigFromJSON(agentCfg)
	config.New().Service.DNA.Agent = &agentFromCFG
	log.Info("agent loaded from JSON, v%", config.New().Service.DNA.Agent.Identifiers.Core)
	return nil
}

func Verify() error {
	//todo: TOM-HA & RL
	//1. Dispatch message to controller - here I AM! here's my ID and labels/selectors
	//2. Controller -> save agent selectors and params (Postgres)
	//3. Controller to Agent confirmation - you've been registered
	//4. save verification / date to agent file
}
