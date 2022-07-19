package integrity

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rills-ai/Hachi/pkg/config"
	"github.com/rills-ai/Hachi/pkg/cryptography"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

const DAT_FILE = "data/__id.dat"

//todo: hide our token :)
const AUTHENTICITY_TOKEN string = "Hachi::"

func GenerateAgentID() string {
	return AUTHENTICITY_TOKEN + "rlx" + strings.Replace(uuid.New().String(), "-", "", -1) + ":" + strconv.FormatInt(time.Now().Unix(), 10)
}
func ProvisionAgent(agentConf config.AgentConfig) string {

	id := ""
	cfg := config.New()
	agentConf.Enabled = cfg.Service.DNA.Agent.Enabled
	if agentConf.Identifiers == nil {
		agentConf.Identifiers = &config.IdentifiersConfig{
			Core:        "",
			Descriptors: cfg.Service.DNA.Agent.Identifiers.Descriptors,
		}
	}

	if agentConf.Identifiers.Core == "" {
		id = GenerateAgentID()
	}

	agentConf.Identifiers.Core = id

	agentJSON, _ := agentConf.ToJSON()
	encryptedJSON := cryptography.Encryption(agentJSON)
	base64JSON := base64.StdEncoding.EncodeToString([]byte(encryptedJSON))
	err := ioutil.WriteFile(DAT_FILE, []byte(base64JSON), 0644)
	if err != nil {
		log.Fatal("Can not create DAT FILE", err)
		return ""
	}
	return agentJSON
}

func ValidateAgentID() (string, error) {
	_, err := os.Stat(DAT_FILE)
	if errors.Is(err, os.ErrNotExist) {
		ProvisionAgent(config.AgentConfig{})
	}
	data, err := ioutil.ReadFile(DAT_FILE)
	if err != nil {
		log.Error("Can not open DAT FILE from storage ", err)
		ProvisionAgent(config.AgentConfig{})
	}

	decodedMessage, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		log.Error("DAT FILE decoding failed, generating new ID  for agent. ", err)
		ProvisionAgent(config.AgentConfig{})
	}
	id := cryptography.Decryption(string(decodedMessage))
	if strings.Contains(id, AUTHENTICITY_TOKEN) {
		log.Info("Valid agent identifier received.")
	} else {
		return "", fmt.Errorf("agent integrity compromised")
	}
	return strings.Replace(id, AUTHENTICITY_TOKEN, "", -1), nil
}
