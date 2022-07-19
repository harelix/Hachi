package haf

import (
	"fmt"
	"github.com/corazawaf/coraza/v2"
	"github.com/corazawaf/coraza/v2/seclang"
	"github.com/rills-ai/Hachi/pkg/messages"
)

/*
	HAF stands for Hachi Application Firewall
	==========================================
	we should include this on Hachi's API and NATS subscribers

	(https://github.com/corazawaf/coraza)

	messages coming on all ingress channels should be tested for validity
	- can a message pass-through to the next processing step
	- can a message act on the specific context (agent, controller)
		- controller
		- can a meesage act
*/

func HafParser(Capsule messages.Capsule) error {
	waf := coraza.NewWaf()
	parser, _ := seclang.NewParser(waf)
	// Now we parse our rules
	if err := parser.FromString(`SecRule REMOTE_ADDR "@rx .*" "id:1,phase:1,deny,status:403"`); err != nil {
		fmt.Println(err)
	}

	return nil
}
