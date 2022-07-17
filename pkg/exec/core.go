package exec

import (
	"fmt"
	"github.com/rills-ai/Hachi/pkg/messages"
)

type ExecutionResponseCode int

const (
	Ok ExecutionResponseCode = iota
	Failure
	Error
	Unknown
)

func (er ExecutionResponseCode) String() string {
	switch er {
	case Ok:
		return "ok"
	case Failure:
		return "failure"
	case Error:
		return "error"
	case Unknown:
		return "unknown"
	}
	return "unknown"
}

type ExecutionCommand interface {
	Exec(cap *messages.Capsule) (*ExecutionResponse, error)
}

type ExecutionResponse struct {
	response ExecutionResponseCode
}

func ProcessIncomingCapsule(capsule messages.Capsule) (ExecutionResponse, error) {
	fmt.Println(capsule)
	//capsule.Route.Remote.GetExecIdentifier()
	return ExecutionResponse{
		response: Ok,
	}, nil
}
