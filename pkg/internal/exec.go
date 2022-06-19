package internal

import (
	"encoding/base64"
	"github.com/rills-ai/Hachi/pkg/cryptography"
	"github.com/rills-ai/Hachi/pkg/messages"
	"strings"
)

type InternalResponse struct {
	Result    string
	Directive string
}

func Exec(capsule messages.Capsule, directive string) InternalResponse {
	args := strings.Split(directive, "#")

	switch args[0] {
	case "__internal__.crypto.encrypt":

		message := cryptography.Encryption(capsule.Message)

		return InternalResponse{
			Result:    base64.StdEncoding.EncodeToString([]byte(message)),
			Directive: directive,
		}
	case "__internal__.crypto.decrypt":
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
