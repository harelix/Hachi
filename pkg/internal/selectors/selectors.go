package selectors

import "fmt"

func GetAgentsForSelectors() {

}

func BuildAgentDedicatedChannelIdentifier(agentID string) string {
	return fmt.Sprintf("agents.dedicated.%v", agentID)
}
