package config

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/gilbarco-ai/Hachi/pkg/helper"
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
	Values                 ValuesConfig
	rawConfigValue         []byte
	rawValues              []byte
	interpolatedConfigFile []byte
	ValuesFile             []byte
}

type Service struct {
	Version int `hcl:"version"`
	Type    DNATypes
	DNA     DNAConfig `hcl:"dna,block"`
}

func New() *HachiConfig {
	once.Do(func() {
		instantiated = &HachiConfig{}
	})
	return instantiated
}

type DNAConfig struct {
	Name       string           `hcl:"name,label"`
	API        APIConfig        `hcl:"api,block"`
	Controller controllerConfig `hcl:"controller,block"`
	Agent      agentConfig      `hcl:"agent,block"`
	Storage    StorageConfig    `hcl:"storage,block"`
	Tracts     TractsConfig     `hcl:"tracts,block"`
	Stream     StreamConfig     `hcl:"stream,block"`
	HRL        HRLConfig        `hcl:"hrl,block"`
	Nats       NatsConfig       `hcl:"nats,block"`
	KV         KVConfig         `hcl:"kv_db,block"`
	Http       ServerConfig     `hcl:"http,block"`
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
	GetIdentifiers() []string
}

type controllerConfig struct {
	Enabled           bool     `hcl:"enabled"`
	InvocationTimeout int      `hcl:"invocation_timeout,optional"`
	Identifiers       []string `hcl:"identifiers,optional"`
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

func (p controllerConfig) GetIdentifiers() []string {
	return p.Identifiers
}

type agentConfig struct {
	Enabled           bool     `hcl:"enabled"`
	InvocationTimeout int      `hcl:"invocation_timeout,optional"`
	Identifiers       []string `hcl:"identifiers,optional"`
}

func (p agentConfig) IsEnabled() bool {
	return p.Enabled
}

func (p agentConfig) GetInvocationTimeout() int {
	return p.InvocationTimeout
}

func (p agentConfig) GetIdentifiers() []string {
	return p.Identifiers
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

type RouteConfig struct {
	Async                      bool                `hcl:"async,optional"`
	Name                       string              `hcl:"name,label"`
	Subject                    []string            `hcl:"subject,optional"`
	Verb                       string              `hcl:"verb"`
	Local                      string              `hcl:"local"`
	Remote                     string              `hcl:"remote,optional"`
	Headers                    map[string][]string `hcl:"headers,optional"`
	IndexedInterpolationValues map[string]string
	Payload                    string `hcl:"payload,optional"`
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

func (config *HachiConfig) ParseFile(filePath string) error {
	var err error
	ctx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"author": cty.StringVal("relix"),
		},
	}

	config.rawConfigValue, err = ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	c, err := config.InterpolateStrings()
	if err != nil {
		return fmt.Errorf("failed to interpolate config: %w", err)
	}

	config.interpolatedConfigFile = []byte(c)
	err = hclsimple.Decode(filePath, config.interpolatedConfigFile, ctx, &config.Service)
	if err != nil {
		return fmt.Errorf("failed to parse HCL file: %w", err)
	}

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

var InterpolationRegex = regexp.MustCompile("{{\\.((local|remote|route|resolver)::(.*?))}}")

// InterpolateStrings  we currently support interpolation from envars and Hachi stanza vars
func (config *HachiConfig) InterpolateStrings() (string, error) {
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

	content := string(config.rawConfigValue)
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
