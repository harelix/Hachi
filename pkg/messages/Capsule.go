package messages

import (
	"encoding/json"
	"github.com/rills-ai/Hachi/pkg/config"
	log "github.com/sirupsen/logrus"
)

// TODO: since most libs will take a byte array for sending, should Message really be a string? - FYI strings are immutable
type Capsule struct {
	Message   string              `json:"message"`
	Headers   map[string][]string `json:"headers"`
	Selectors []string            `json:"selectors"`
	Route     *config.RouteConfig `json:"route"`
}

func CapsuleFromJSON(capsuleString string) (Capsule, error) {
	var interpolatedCapsule Capsule
	err := json.Unmarshal([]byte(capsuleString), &interpolatedCapsule)

	if err != nil {
		log.Warning("capsule message failed to Unmarshal, err: %w", err)
		return Capsule{}, err
	}
	return interpolatedCapsule, nil
}

func (c *Capsule) JSONFRomCapsule() (string, error) {
	marshaledCapsule, err := json.Marshal(c)
	if err != nil {
		log.Error("capsule message failed to interpolate, err: %w", err)
		return "", err
	}
	capsuleString := string(marshaledCapsule)
	return capsuleString, nil
}
