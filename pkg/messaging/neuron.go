package messaging

import (
	"fmt"
	"github.com/rills-ai/Hachi/pkg/messaging/consumers"
	"strconv"
	"strings"
	"sync"
	"time"

	"context"
	"encoding/json"
	"github.com/nats-io/nats.go"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/messages"
	log "github.com/sirupsen/logrus"
)

var once sync.Once
var instance *HachiNeuron

type HachiNeuron struct {
	NC                   *nats.Conn
	JS                   nats.JetStreamContext
	ID                   string
	Agent                config.IAgent
	DefaultNATSRegisters map[string]DefaultNATSRegister
	bufferSize           int
}

func (hn *HachiNeuron) GetBufferSize() int {
	return hn.bufferSize
}

func Get() *HachiNeuron {
	once.Do(func() {
		instance = &HachiNeuron{}
		instance.bufferSize = 10
	})
	return instance
}

// RegisterMiddleware : for future use: register Hachi middleware
func (hn *HachiNeuron) RegisterMiddleware() error {
	return nil
}

func (hn *HachiNeuron) Init(ctx context.Context) error {
	return hn.Connect(ctx, 1)
}

func (hn *HachiNeuron) Connect(ctx context.Context, retry int) error {

	natsConf := config.New().Service.DNA.Nats
	hn.Agent = config.New().IAM
	hn.ID = HachiContext.ApplicationName
	// There is not a single listener for connection events in the NATS Go Client.
	// Instead, you can set individual event handlers using:

	connect, err := nats.Connect(natsConf.Address+":"+strconv.Itoa(natsConf.Port),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Printf("client disconnected: %v", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Printf("client reconnected")
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Printf("client closed")
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			//retry connection
		}))

	if err != nil {
		if retry <= HachiContext.MaxReconnectAttemptsToServer {
			retry += 1
			log.Printf("Error connecting to server; retry attempt #%d in 5 seconds: %v", retry, err)
			time.Sleep(time.Duration(HachiContext.ConnectingToServerIntervalInSeconds) * time.Second)
			return hn.Connect(ctx, retry)

		}
		return fmt.Errorf("NATS server is unavailable: Error connecting to server will not attempt further retries: %w", err)
	}

	hn.NC = connect
	hn.JS, err = hn.NC.JetStream()
	if err != nil {
		return fmt.Errorf("failed to initiate JetSteam connection: %w", err)
	}

	log.WithFields(log.Fields{
		"agent_type": config.New().IAM.GetType(),
		"NATS_addr":  config.New().Service.DNA.Nats.Address,
	}).Info("successfully connected to NATS server.")
	//Provision our Default JETStream Consumers and Handlers
	//======================================================
	hn.ProvisionNATSJetStream()
	//======================================================
	go hn.registerDedicatedConsumers()

	//if hn.NC == nil {
	//	return errors.New("could not connect to NATS on " + natsConf.Address + ":" + strconv.Itoa(natsConf.Port))
	//
	//}
	fmt.Println("Connected to NATS on " + natsConf.Address + ":" + strconv.Itoa(natsConf.Port))

	err = hn.NC.Flush()
	if err != nil {
		return fmt.Errorf("failed to flush connection to NATS: %w", err)
	}

	return nil
}

type DefaultNATSRegister struct {
	ListeningOn  *DefaultNATSRegister
	StreamName   string
	Subject      string
	ConsumerName string
	Active       bool
}

func (hn *HachiNeuron) ProvisionNATSJetStream() {
	//if config.New().Service.AgentType != config.controller {
	//	return
	//}
	controllerNATSConfig := DefaultNATSRegister{
		Active:       true,
		StreamName:   HachiContext.ControllerNATSDefaultStream,
		Subject:      HachiContext.ControllerNATSDefaultEfferentSubjects,
		ConsumerName: HachiContext.ControllerNATSDefaultStreamConsumerName}

	agentNATSConfig := DefaultNATSRegister{
		Active:       true,
		ListeningOn:  &controllerNATSConfig,
		StreamName:   HachiContext.AgentsNATSDefaultStream,
		Subject:      HachiContext.AgentsNATSDefaultAfferentSubjects,
		ConsumerName: HachiContext.AgentsNATSDefaultStreamConsumerName}

	controllerNATSConfig.ListeningOn = &agentNATSConfig

	hn.DefaultNATSRegisters = map[string]DefaultNATSRegister{
		strings.ToLower(config.Controller.String()): controllerNATSConfig,
		strings.ToLower(config.Agent.String()):      agentNATSConfig}

	err := NATSDefaultProvisioning(hn)
	if err != nil {
		log.Fatal("failed to provision JetStream: %w", err)
	}

}

var defaultConsumer = make(chan *PublishedMessage)

func (hn *HachiNeuron) registerDedicatedConsumers() {

	go bindDefaultPullSubscriber(hn, defaultConsumer)

	for {
		hn.handleIncomingMessage(<-defaultConsumer)
	}
}

type PublishedMessage struct {
	Message *nats.Msg
	Error   error
}

func (hn *HachiNeuron) handleIncomingMessage(pu *PublishedMessage) {
	message := pu.Message
	if message == nil {
		return
	}
	fmt.Println(message.Data)

	//Http Dispatch storix
	/*tracing.Trace(tracing.HachiContext{
		Method:      "",
		Path:        "",
		Resource:    "",
		ServiceName: "",
		ServiceType: "",
	})*/
	//invokeSink/Exec
	e := message.Ack()
	if e != nil {
		println("bind default pull subscriber: %w", e)
	}

}

/*
!!!!!!!!!!!! DO NOT DELETE THIS ONE !!!!!!!!!!!!!!!!!!
func (hn *HachiNeuron) SubscribeDefault(){
	subscriber, err := hn.JS.PullSubscribe(HachiContext.NATSDefaultSubjects, HachiContext.NATSDefaultStreamConsumerName,
		nats.PullMaxWaiting(512))
	noerr(err)

	for {
		cMessages, _ := subscriber.Fetch(10)
		for _, cMsg := range cMessages {
			err := cMsg.Ack()
			noerr(err)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(cMsg.Data) //exec router
			log.Println("execution-command recived")
		}
	}
}
*/

func (hn *HachiNeuron) Subscribe(subjects []string) error {
	for _, subject := range subjects {
		//todo: handle err
		hn.NC.Subscribe(subject, func(message *nats.Msg) {
			consumers.IncomingMessageHandler(message)
		})
	}
	return nil
}

func (hn *HachiNeuron) Publish(ctx context.Context, capsule messages.Capsule) error {
	msg, err := json.Marshal(capsule)
	if err != nil {
		return err
	}
	for _, subject := range capsule.Selectors {
		err = hn.NC.Publish(subject, []byte(msg))
		return err
	}
	return nil
}

func (hn *HachiNeuron) Close() {
	hn.NC.Close()
}
