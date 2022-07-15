package integrity

import (
	"encoding/base64"
	"errors"
	"github.com/google/uuid"
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

func ProvisionAgentID() string {

	id := AUTHENTICITY_TOKEN + "rlx" + strings.Replace(uuid.New().String(), "-", "", -1) + ":" + strconv.FormatInt(time.Now().Unix(), 10)
	encryptedID := cryptography.Encryption(id)
	base64Id := base64.StdEncoding.EncodeToString([]byte(encryptedID))
	err := ioutil.WriteFile(DAT_FILE, []byte(base64Id), 0644)
	if err != nil {
		log.Fatal("Can not create DAT FILE", err)
		return ""
	}
	return id
}

func ValidateAgentID() string {
	_, err := os.Stat(DAT_FILE)
	//errors.Is(nil)
	if errors.Is(err, os.ErrNotExist) {
		return ProvisionAgentID()
	}
	data, err := ioutil.ReadFile(DAT_FILE)
	if err != nil {
		log.Error("Can not open DAT FILE from storage ", err)
		return ProvisionAgentID()
	}

	decodedMessage, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		log.Error("DAT FILE decoding failed, generating new ID  for agent. ", err)
		return ProvisionAgentID()
	}
	id := cryptography.Decryption(string(decodedMessage))
	if strings.Contains(id, AUTHENTICITY_TOKEN) {
		log.Info("Valid agent identifier received.")
	} else {
		log.Error("Invalid Identifier for Hachi agent, invalid value is: %v. Initializing new ID and sending notification.", id)
		return ProvisionAgentID()
	}

	return strings.Replace(id, AUTHENTICITY_TOKEN, "", -1)
}
