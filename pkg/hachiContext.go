package HachiContext

import "errors"

const (
	//auto-heal
	ConnectingToServerIntervalInSeconds int64  = 5
	MaxReconnectAttemptsToServer        int    = 5
	ApplicationName                     string = "8HACHI"
	ContextIAM                          string = "IAM"
	ContextCapsule                      string = "capsule"
	ContextTimeoutMax                   int    = 1000 * 60

	//NATS - controller
	ControllerNATSDefaultStream string = "controllerToagents-NeuroStream"
	//from controller to agents
	ControllerNATSDefaultEfferentSubjects   string = "neurostream.controller.to.agents"
	ControllerNATSDefaultStreamConsumerName string = "controllerDefaultNeuroStreamConsumer"
	ControllerNATSDefaultStreamUPMessage    string = "neuro-stream link is connected from controller"

	//NATS - agents
	AgentsNATSDefaultStream string = "agentsTocontroller-NeuroStream"
	//from agents to controller
	AgentsNATSDefaultAfferentSubjects   string = "neurostream.agent.to.controller"
	AgentsNATSDefaultStreamConsumerName string = "agentDefaultNeuroStreamConsumer"
	AgentsNATSDefaultStreamUPMessage    string = "neuro-stream link is connected from agent"

	//Server responses
	PublishSuccessful string = "message hyper-jumped through wormhole, i.e. everything is groovy, we're golden!"
)

var (
	ErrBodyEmpty = errors.New("empty body content is not allowed in POST handler")
)
