package internal

import (
	HachiContext "github.com/rills-ai/Hachi/pkg"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/integrity"
	"github.com/rills-ai/Hachi/pkg/messages"
	"reflect"
	"testing"
)

func Test_registerAgent(t *testing.T) {

	agentCfg := config.AgentConfig{
		Enabled:           false,
		InvocationTimeout: 0,
		Identifiers: &config.IdentifiersConfig{
			Core:        integrity.GenerateAgentID(),
			Descriptors: []string{"agents.generation.alpha", "inception.author.relix", "context.yearly.seasonal.summer"},
		},
		VerifiedOn: "",
		Verified:   false,
	}

	agentJSONString, _ := agentCfg.ToJSON()

	type args struct {
		capsule messages.Capsule
	}
	tests := []struct {
		name    string
		args    args
		want    messages.ExecutionResponse
		wantErr bool
	}{
		{
			name: "",
			args: args{
				capsule: messages.Capsule{
					Message:   agentJSONString,
					Headers:   nil,
					Selectors: nil,
					Route: &config.RouteConfig{
						Async:     false,
						Name:      "",
						Selectors: nil,
						Verb:      "",
						Local:     "",
						Remote: config.RemoteExecConfig{
							HTTP:    nil,
							SSH:     nil,
							Webhook: nil,
							Internal: &config.InternalConfig{
								Type: HachiContext.RegisterAgentInternalCommand,
							},
						},
						Headers:                    nil,
						IndexedInterpolationValues: nil,
						Payload:                    "",
					},
				},
			},
			want:    messages.ExecutionResponse{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := registerAgent(tt.args.capsule)
			if (err != nil) != tt.wantErr {
				t.Errorf("registerAgent() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("registerAgent() got = %v, want %v", got, tt.want)
			}
		})
	}
}
