package core

import (
	"context"
	"fmt"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/internal"
	"github.com/rills-ai/Hachi/pkg/messages"
)

type ExecutionCommand interface {
	Exec(cap *messages.Capsule) (*messages.ExecutionResponse, error)
}

func ProcessIncomingCapsule(ctx context.Context, capsule messages.Capsule) (messages.ExecutionResponse, error) {

	execQualifier := capsule.Route.Remote.GetExecIdentifier()

	switch execQualifier {
	case config.Internal:
		_, err := internal.ProcessCapsule(ctx, capsule)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println(execQualifier)

	return messages.ExecutionResponse{}, nil
}
