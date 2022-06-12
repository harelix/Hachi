package messages

import "github.com/gilbarco-ai/Hachi/pkg/config"

// TODO: since most libs will take a byte array for sending, should Message really be a string? - FYI strings are immutable
type Capsule struct {
	Message string              `json:"message"`
	Headers map[string][]string `json:"headers"`
	Subject []string            `json:"subject"`
	Route   *config.RouteConfig `json:"route"`
}
