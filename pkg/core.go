package HachiContext

import (
	"fmt"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/internal"
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

func ProcessIncomingCapsule(ctx context.Context, capsule messages.Capsule) (ExecutionResponse, error) {

	execQualifier := capsule.Route.Remote.GetExecIdentifier()
	fmt.Println(execQualifier)
	fmt.Println(config.Internal)
	fmt.Println(internal.InternalResponse{})
	/*
		switch execQualifier {
		case config.Internal:
			_, err := internal.ProcessCapsule(ctx, capsule)
			if err != nil {
				fmt.Println(err)
			}
		}
		fmt.Println(execQualifier)
	*/
	return ExecutionResponse{
		response: Ok,
	}, nil
}
