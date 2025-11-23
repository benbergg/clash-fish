package config

// GetDefaultConfig 返回默认配置
func GetDefaultConfig() *Config {
	return &Config{
		Port:               7890,
		SocksPort:          7891,
		AllowLan:           false,
		Mode:               "rule",
		LogLevel:           "info",
		ExternalController: "127.0.0.1:9090",
		TUN: TUNConfig{
			Enable:              true,
			Stack:               "system",
			DNSHijack:           []string{"any:53"},
			AutoRoute:           true,
			AutoDetectInterface: true,
		},
		DNS: DNSConfig{
			Enable:       true,
			Listen:       "198.18.0.2:53",
			EnhancedMode: "fake-ip",
			FakeIPRange:  "198.18.0.1/16",
			Nameserver: []string{
				"223.5.5.5",
				"114.114.114.114",
			},
			Fallback: []string{
				"tls://1.1.1.1:853",
				"tls://8.8.8.8:853",
			},
		},
		Proxies: []Proxy{
			{
				Name:     "example-proxy",
				Type:     "ss",
				Server:   "example.com",
				Port:     8388,
				Cipher:   "aes-256-gcm",
				Password: "password",
				UDP:      true,
			},
		},
		ProxyGroups: []ProxyGroup{
			{
				Name: "PROXY",
				Type: "select",
				Proxies: []string{
					"example-proxy",
					"DIRECT",
				},
			},
		},
		Rules: []string{
			"GEOIP,PRIVATE,DIRECT",
			"GEOIP,CN,DIRECT",
			"MATCH,PROXY",
		},
	}
}

// GetExampleConfig 返回示例配置（带注释说明）
func GetExampleConfig() string {
	return `# Clash-Fish Configuration File
# 更多配置选项请参考: https://wiki.metacubex.one/

# HTTP 代理端口
port: 7890

# SOCKS5 代理端口
socks-port: 7891

# 允许局域网连接
allow-lan: false

# 代理模式: rule(规则) / global(全局) / direct(直连)
mode: rule

# 日志级别: info / warning / error / debug / silent
log-level: info

# RESTful API 控制端口
external-controller: 127.0.0.1:9090

# TUN 模式配置（透明代理核心）
tun:
  enable: true
  stack: system                # system / gvisor
  dns-hijack:
    - any:53                   # 劫持所有 DNS 请求
  auto-route: true             # 自动配置路由表
  auto-detect-interface: true  # 自动检测网卡

# DNS 配置
dns:
  enable: true
  listen: 198.18.0.2:53
  enhanced-mode: fake-ip       # fake-ip / redir-host
  fake-ip-range: 198.18.0.1/16
  nameserver:
    - 223.5.5.5                # 阿里 DNS
    - 114.114.114.114          # 114 DNS
  fallback:
    - tls://1.1.1.1:853        # Cloudflare DNS over TLS
    - tls://8.8.8.8:853        # Google DNS over TLS

# 代理服务器配置
proxies:
  - name: "example-proxy"
    type: ss
    server: example.com
    port: 8388
    cipher: aes-256-gcm
    password: password
    udp: true

# 代理组配置
proxy-groups:
  - name: "PROXY"
    type: select
    proxies:
      - example-proxy
      - DIRECT

# 规则配置
rules:
  - GEOIP,PRIVATE,DIRECT       # 私网地址直连（VPN 内网会走这里）
  - GEOIP,CN,DIRECT            # 国内地址直连
  - MATCH,PROXY                # 其他流量走代理
`
}
