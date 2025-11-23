package config

// Config 主配置结构
type Config struct {
	Port               int          `yaml:"port"`
	SocksPort          int          `yaml:"socks-port"`
	AllowLan           bool         `yaml:"allow-lan"`
	Mode               string       `yaml:"mode"`
	LogLevel           string       `yaml:"log-level"`
	ExternalController string       `yaml:"external-controller"`
	TUN                TUNConfig    `yaml:"tun"`
	DNS                DNSConfig    `yaml:"dns"`
	Proxies            []Proxy      `yaml:"proxies"`
	ProxyGroups        []ProxyGroup `yaml:"proxy-groups"`
	Rules              []string     `yaml:"rules"`
}

// TUNConfig TUN 模式配置
type TUNConfig struct {
	Enable              bool     `yaml:"enable"`
	Stack               string   `yaml:"stack"`
	DNSHijack           []string `yaml:"dns-hijack"`
	AutoRoute           bool     `yaml:"auto-route"`
	AutoDetectInterface bool     `yaml:"auto-detect-interface"`
}

// DNSConfig DNS 配置
type DNSConfig struct {
	Enable       bool     `yaml:"enable"`
	Listen       string   `yaml:"listen"`
	EnhancedMode string   `yaml:"enhanced-mode"`
	FakeIPRange  string   `yaml:"fake-ip-range"`
	Nameserver   []string `yaml:"nameserver"`
	Fallback     []string `yaml:"fallback"`
}

// Proxy 代理配置
type Proxy struct {
	Name     string                 `yaml:"name"`
	Type     string                 `yaml:"type"`
	Server   string                 `yaml:"server"`
	Port     int                    `yaml:"port"`
	Cipher   string                 `yaml:"cipher,omitempty"`
	Password string                 `yaml:"password,omitempty"`
	UDP      bool                   `yaml:"udp,omitempty"`
	Plugin   string                 `yaml:"plugin,omitempty"`
	PluginOpts map[string]interface{} `yaml:"plugin-opts,omitempty"`
}

// ProxyGroup 代理组配置
type ProxyGroup struct {
	Name    string   `yaml:"name"`
	Type    string   `yaml:"type"`
	Proxies []string `yaml:"proxies"`
	URL     string   `yaml:"url,omitempty"`
	Interval int     `yaml:"interval,omitempty"`
}
