package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/google/uuid"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/common-nighthawk/go-figure"
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/api"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/messaging"
	log "github.com/sirupsen/logrus"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}

func main() {
	log.SetReportCaller(true)

	go PrintHachiWelcome()

	var confile = flag.String("config", "", "default configuration file path")
	var valfile = flag.String("values", "", "default values file path")

	flag.Parse()

	cLog := log.WithFields(log.Fields{"app_name": "Hachi"})

	if *confile == "" {
		cLog.Error("Hachi err! missing cli args 'c' path to config file!")
		os.Exit(1)
	}

	//init config file
	err := config.New().LoadStanzaValues(*valfile)
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	err = config.New().ParseFile(*confile)
	if err != nil {
		fmt.Println("ERROR: " + err.Error())
		log.Printf(err.Error())
		os.Exit(1)
	}

	SelfProvisioning(config.New())
	//config main service context
	ctx, cancel := context.WithCancel(context.WithValue(context.Background(),
		HachiContext.ContextIAM, config.New().IAM))
	defer cancel()

	//bootstrap server
	go api.StartAPIServer(ctx)

	//init Hachi Neuron
	err = messaging.Get().Init(ctx)
	if err != nil {
		cLog.Errorf("failed to init Hachi: %v", err)
		os.Exit(1)
	}
	//go messaging.Get().SubscribeDefault()
	//messaging.Get().Subscribe(dendrite.GetSubscriptionSubjects(config.New()))
	//NATS connection close
	defer messaging.Get().Close()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	messaging.Get().Close()
	ctx.Done()
}

func SelfProvisioning(config *config.HachiConfig) {
	if config.Service.DNA.Controller.Enabled {
		//todo: maybe a constant identifier
	} else {
		config.Service.DNA.Agent.Identifiers =
			append(config.Service.DNA.Agent.Identifiers, fmt.Sprintf("%vRND%v", HachiContext.DefaultAgentIdentifierPrefix,
				strings.Replace(uuid.New().String(), "-", "", -1)))
	}
}

func PrintHachiWelcome() {
	time.Sleep(50 * time.Millisecond)
	myFigure := figure.NewFigure("8//Hachi", "doom", true)
	myFigure.Print()
	message := fmt.Sprintf("\nHachi@Relix, instance name '%s', version %d", config.New().Service.DNA.Name, config.New().Service.Version)
	fmt.Println(message)
}
