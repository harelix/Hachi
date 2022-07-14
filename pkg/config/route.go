package config

type RouteConfig struct {
	Async                      bool                `hcl:"async,optional"`
	Name                       string              `hcl:"name,label"`
	Selectors                  []string            `hcl:"selectors,optional"`
	Verb                       string              `hcl:"verb"`
	Local                      string              `hcl:"local"`
	Remote                     RemoteExecConfig    `hcl:"remote,block"`
	Headers                    map[string][]string `hcl:"headers,optional"`
	IndexedInterpolationValues map[string]string
	Payload                    string `hcl:"payload,optional"`
}
