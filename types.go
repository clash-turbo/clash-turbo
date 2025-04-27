package main

type appConfig struct {
	ProxyType    int         `yaml:"proxy_type"`
	GitHubProxy  string      `yaml:"github_proxy"`
	HTTPProxy    string      `yaml:"http_proxy"`
	Core         core        `yaml:"core"`
	GuiPort      int         `yaml:"gui_port"`
	Profiles     []profile   `yaml:"profiles"`
	PrependRules []string    `yaml:"prepend_rules"`
	CoverConfig  coverConfig `yaml:"cover_config"`
}

type core struct {
	Port               int    `yaml:"port"`
	MixedPort          int    `yaml:"mixed_port"`
	ExternalController string `yaml:"external_controller"`
	Secret             string `yaml:"secret"`
}

type profile struct {
	Name    string   `yaml:"name"`
	URL     string   `yaml:"url"`
	Filters []filter `yaml:"filters"`
}

type filter struct {
	Name   string `yaml:"name"`
	FType  string `yaml:"type"`
	Filter string `yaml:"filter"`
}

type coverConfig struct {
	Dns          map[string]any `yaml:"dns"`
	Tun          map[string]any `yaml:"tun"`
	Sniffer      map[string]any `yaml:"sniffer"`
	Rules        []any          `yaml:"rules"`
	Ipv6         bool           `yaml:"ipv6"`
	AllowLan     bool           `yaml:"allow-lan"`
	UnifiedDelay bool           `yaml:"unified-delay"`
	//ExternalUiUrl string                 `yaml:"external_ui_url"`
}
