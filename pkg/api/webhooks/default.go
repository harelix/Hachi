package webhooks

import (
	"bytes"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/rills-ai/Hachi/pkg/api/helpers"
	"net/http"
	"sync"
	"time"
)

var once sync.Once

var instance *Dispatcher

type Dispatcher struct {
	client           *http.Client
	registers        map[string]WebhookRegister
	base             *echo.Group
	mu               *sync.Mutex
	notificationChan chan string
	stopChan         chan bool
}

type WebhookRegister struct {
	// Header property of the request, matching Storix Server property name
	Name        string `json:"name" validate:"required"`
	Destination string `json:"destination" validate:"required"`
	//registers represent the events that this webhook is subscribing to
	Registers []string `json:"registers" validate:"required"`
	Token     string   `json:"token" validate:"required"`
}

func Construct() *Dispatcher {

	once.Do(func() {

		instance = &Dispatcher{
			client:           &http.Client{},
			registers:        make(map[string]WebhookRegister),
			mu:               &sync.Mutex{},
			notificationChan: make(chan string),
			stopChan:         make(chan bool),
		}
	})

	return instance
}

func (d *Dispatcher) BindWebhooks(base *echo.Group) {
	d.base = base
	d.bindDefaultWebhook()
	go d.hookDispatchChannel()
}

func (d *Dispatcher) Notify(message string) {
	d.notificationChan <- message
}

func (d *Dispatcher) bindDefaultWebhook() {

	d.base.POST("/webhooks/", func(c echo.Context) error {

		webhookRegister, err := helpers.BindAndValidate[WebhookRegister](c)
		if err != nil {
			responseError := &helpers.HachiResponseMessage{
				Error:   true,
				Message: "Bind and Validate request failed",
			}
			return c.JSON(http.StatusBadRequest, responseError)
		}

		d.subscribe(webhookRegister.Name, *webhookRegister)
		registerResponse := &helpers.HachiResponseMessage{
			Error:   false,
			Message: "Webhook registered",
		}
		return c.JSON(http.StatusOK, registerResponse)
	})
}

func (d *Dispatcher) hookDispatchChannel() {

	for {
		select {
		case message := <-d.notificationChan:
			fmt.Println("received", message)
			d.dispatch(message)
		case stopMessage := <-d.stopChan:
			fmt.Println("a stop dispatching webhooks command", stopMessage)
		}
	}
}

func (d *Dispatcher) subscribe(name string, webhook WebhookRegister) {
	d.mu.Lock()
	d.registers[name] = webhook
	d.mu.Unlock()
}

func (d *Dispatcher) dispatch(message string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	for name, webhookRegister := range d.registers {
		fmt.Println(message)
		go func(name string, webhookRegister WebhookRegister) {
			req, err := http.NewRequest("POST", webhookRegister.Destination, bytes.NewBufferString(fmt.Sprintf("Hello %s, current time is %s", name, time.Now().String())))
			if err != nil {
				// probably don't allow creating invalid destinations
				return
			}

			resp, err := d.client.Do(req)
			if err != nil {
				// should probably check response status code and retry if it's timeout or 500
				return
			}

			fmt.Printf("Webhook to '%s' dispatched, response code: %d \n", webhookRegister, resp.StatusCode)

		}(name, webhookRegister)
	}
}
