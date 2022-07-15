package config

import (
	"fmt"
	"github.com/rills-ai/Hachi/pkg/interpolator"
	"io/ioutil"
	"sync"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsimple"
	"github.com/zclconf/go-cty/cty"
)

var once sync.Once
var instantiated *HachiConfig

type DNATypes int64

const (
	Agent      DNATypes = iota
	Controller          //todo: check for duplicates consider cluster
	Me
)

func (at DNATypes) String() string {
	switch at {
	case Controller:
		return "controller"
	case Agent:
		return "agent"
	case Me:
		return "relix"
	}
	return "unknown"
}

type HachiConfig struct {
	IAM                    IAgent
	Service                Service
	Internals              ServiceInternals
	Values                 ValuesConfig
	rawConfigValue         []byte
	internalTractsConfig   []byte
	rawValues              []byte
	interpolatedConfigFile []byte
	ValuesFile             []byte
}

type Service struct {
	Version int `hcl:"version"`
	Type    DNATypes
	DNA     DNAConfig `hcl:"dna,block"`
}

type ServiceInternals struct {
	DNA InternalDNA `hcl:"dna,block"`
}

type InternalDNA struct {
	Name   string       `hcl:"name,label"`
	Tracts TractsConfig `hcl:"tracts,block"`
	//Webhooks TractsConfig `hcl:"webhooks,block"`
}

func New() *HachiConfig {
	once.Do(func() {
		instantiated = &HachiConfig{}
	})
	return instantiated
}

type DNAConfig struct {
	Name           string            `hcl:"name,label"`
	API            APIConfig         `hcl:"api,block"`
	Controller     *controllerConfig `hcl:"controller,block"`
	Agent          *agentConfig      `hcl:"agent,block"`
	Storage        StorageConfig     `hcl:"storage,block"`
	Tracts         TractsConfig      `hcl:"tracts,block"`
	Stream         StreamConfig      `hcl:"stream,block"`
	HRL            HRLConfig         `hcl:"hrl,block"`
	Nats           NatsConfig        `hcl:"nats,block"`
	KV             KVConfig          `hcl:"kv_db,block"`
	Http           ServerConfig      `hcl:"http,block"`
	InternalTracts TractsConfig
}

type ValuesConfig struct {
	Values map[string]string `hcl:"values,attr"`
}

type ServerConfig struct {
	Addr string `hcl:"addr"`
	Port int    `hcl:"port"`
}
type APIConfig struct {
	Enabled   bool `hcl:"enabled"`
	AllowList bool `hcl:"allow_list"`
	Auth      Auth `hcl:"auth,block"`
	Version   int  `hcl:"version"`
}

type Auth struct {
	Enabled     bool   `hcl:"enabled"`
	TokenPrefix string `hcl:"token_prefix"`
	Provider    string `hcl:"provider"`
}

func GetAgent(srv Service) IAgent {
	if srv.Type == Controller {
		return srv.DNA.Controller
	} else {
		return srv.DNA.Agent
	}
}

type IAgent interface {
	GetType() DNATypes
	IsEnabled() bool
	GetInvocationTimeout() int
	GetIdentifiers() IdentifiersConfig
}

type controllerConfig struct {
	Enabled           bool               `hcl:"enabled"`
	InvocationTimeout int                `hcl:"invocation_timeout,optional"`
	Identifiers       *IdentifiersConfig `hcl:"identifiers,block"`
}

func (p controllerConfig) GetType() DNATypes {
	return Controller
}

func (p controllerConfig) IsEnabled() bool {
	return p.Enabled
}

func (p controllerConfig) GetInvocationTimeout() int {
	return p.InvocationTimeout
}

func (p controllerConfig) GetIdentifiers() IdentifiersConfig {
	return *p.Identifiers
}

type agentConfig struct {
	Enabled           bool               `hcl:"enabled"`
	InvocationTimeout int                `hcl:"invocation_timeout,optional"`
	Identifiers       *IdentifiersConfig `hcl:"identifiers,block"`
}

type IdentifiersConfig struct {
	Core        string   `hcl:"core"`
	Descriptors []string `hcl:"descriptors"`
}

func (p agentConfig) IsEnabled() bool {
	return p.Enabled
}

func (p agentConfig) GetInvocationTimeout() int {
	return p.InvocationTimeout
}

func (p agentConfig) GetIdentifiers() IdentifiersConfig {
	return *p.Identifiers
}

func (p agentConfig) GetType() DNATypes {
	return Agent
}

type StorageConfig struct {
	DataDir string `hcl:"data_dir"`
}

type TractsConfig struct {
	Streams []RouteConfig `hcl:"stream,block"`
}

type RemoteExecConfig struct {
	HTTP     *HTTPExecConfig `hcl:"http,block"`
	SSH      *SSHExecConfig  `hcl:"ssh,block"`
	Webhook  *WebhookConfig  `hcl:"webhook,block"`
	Internal *InternalConfig `hcl:"internal,block"`
}

type InternalConfig struct {
	Type string `hcl:"directive"`
}

type WebhookConfig struct {
	Event string `hcl:"event"`
}

type HTTPExecConfig struct {
	URL string `hcl:"url"`
}
type SSHExecConfig struct {
}

type StreamConfig struct {
	CircuitBreaker CircuitBreakerConfig `hcl:"circuit_breaker,block"`
	Deduping       DedupingConfig       `hcl:"deduping,block"`
}

type HRLConfig struct {
	Crypto CryptoConfig `hcl:"crypto,block"`
}

type CryptoConfig struct {
	Provider        string `hcl:"provider"`
	EncryptEndpoint string `hcl:"encrypt_endpoint"`
	DecryptEndpoint string `hcl:"decrypt_endpoint"`
}

type NatsConfig struct {
	Address string `hcl:"addr"`
	Port    int    `hcl:"port"`
}

type CircuitBreakerConfig struct {
	Enabled     bool `hcl:"enabled"`
	MaxRequests int  `hcl:"max_requests"`
	Interval    int  `hcl:"interval"`
	Timeout     int  `hcl:"timeout"`
}

type KVConfig struct {
	//EncryptionMode string `hcl:""`
}

type DedupingConfig struct {
	Enabled  bool   `hcl:"enabled"`
	Strategy string `hcl:"strategy"`
}

func (config *HachiConfig) AppendInternalsToConfigFile() {
	config.Service.DNA.Tracts.Streams = append(config.Service.DNA.Tracts.Streams, config.Internals.DNA.Tracts.Streams...)
}

func (config *HachiConfig) ParseFile(filePath string) error {

	var err error
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"author": cty.StringVal("relix"),
		},
	}

	//init interpolation values from Stanza and Envars
	interpolator.InitInterpolationValues(config.Values.Values)
	//main configuration file parsing
	config.rawConfigValue, err = ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	baseConfigContent := string(config.rawConfigValue)
	bc, err := interpolator.InterpolateStrings(baseConfigContent)
	if err != nil {
		return fmt.Errorf("failed to interpolate config: %w", err)
	}

	//internal configuration file parsing
	config.internalTractsConfig, err = ioutil.ReadFile("conf.d/internals/internals.hcl")
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	internalConfigContent := string(config.internalTractsConfig)
	ic, err := interpolator.InterpolateStrings(internalConfigContent)
	if err != nil {
		return fmt.Errorf("failed to interpolate config: %w", err)
	}

	config.interpolatedConfigFile = []byte(bc)
	err = hclsimple.Decode(filePath, config.interpolatedConfigFile, ctx, &config.Service)
	if err != nil {
		return fmt.Errorf("failed to parse HCL file: %w", err)
	}

	tempString := []byte(ic)
	err = hclsimple.Decode(filePath, tempString, ctx, &config.Internals)
	if err != nil {
		return fmt.Errorf("failed to parse HCL file: %w", err)
	}

	config.AppendInternalsToConfigFile()
	//if config.Service.Agent.controller.Enabled() {
	if config.Service.DNA.Controller.Enabled {
		config.Service.Type = Controller
	}
	config.IAM = GetAgent(config.Service)

	return nil
}

func (config *HachiConfig) LoadStanzaValues(filePath string) error {
	var err error
	config.rawValues, err = ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	err = hclsimple.Decode(filePath, config.rawValues, nil, &config.Values)
	if err != nil {
		return fmt.Errorf("failed to parse HCL file: %w", err)
	}

	return nil
}

/*
var InterpolationRegex = regexp.MustCompile("{{\\.((local|remote|route|resolver)::(.*?))}}")

// InterpolateStrings  we currently support interpolation from envars and Hachi stanza vars
func (config *HachiConfig) InterpolateStrings(content string) (string, error) {
	//stanza vars override envars values
	stanza_vars := config.Values.Values
	lstanza_vars := helper.MapKeys[string, string](stanza_vars, strings.ToLower)

	//index envars
	envars := make(map[string]string)
	for _, e := range os.Environ() {
		before, after, ok := strings.Cut(e, "=")
		if !ok {
			continue
		}
		envars[strings.ToLower(before)] = after
	}

	matches := InterpolationRegex.FindAllString(content, -1)
	for _, v := range matches {
		interpolatedPlaceholder := v
		instructions := InterpolationRegex.FindStringSubmatch(v)
		instruct := instructions[2]
		key := instructions[3]
		//todo: add resolver implementation in the future
		if instruct == "local" {
			interpolatedValue := envars[key]
			if val, ok := lstanza_vars[key]; ok {
				content = strings.Replace(content, interpolatedPlaceholder, val, -1)
			} else {
				content = strings.Replace(content, interpolatedPlaceholder, interpolatedValue, -1)
				if interpolatedValue == "" {
					return "", errors.New("ERROR: key " + interpolatedPlaceholder + " value is missing, check your configuration file and machine ENVARS")
				}
			}
		}
	}
	return content, nil
}
*/
