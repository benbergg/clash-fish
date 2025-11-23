package system

import (
	"net"
	"strings"
)

// VPNInfo VPN 连接信息
type VPNInfo struct {
	Active    bool
	Interface string
	IP        string
	Network   string
}

// DetectVPN 检测 VPN 连接
func DetectVPN() (*VPNInfo, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	for _, iface := range interfaces {
		// 检查 utun 接口（macOS VPN 特征）
		if strings.HasPrefix(iface.Name, "utun") {
			addrs, err := iface.Addrs()
			if err != nil {
				continue
			}

			for _, addr := range addrs {
				ipNet, ok := addr.(*net.IPNet)
				if !ok {
					continue
				}

				ip := ipNet.IP
				// 检查是否是 IPv4 私网地址（VPN 特征）
				if ip.To4() != nil && isPrivateIP(ip) {
					return &VPNInfo{
						Active:    true,
						Interface: iface.Name,
						IP:        ip.String(),
						Network:   ipNet.String(),
					}, nil
				}
			}
		}
	}

	return &VPNInfo{Active: false}, nil
}

// isPrivateIP 判断是否是私网地址
func isPrivateIP(ip net.IP) bool {
	privateRanges := []string{
		"10.0.0.0/8",
		"172.16.0.0/12",
		"192.168.0.0/16",
	}

	for _, cidr := range privateRanges {
		_, ipnet, err := net.ParseCIDR(cidr)
		if err != nil {
			continue
		}
		if ipnet.Contains(ip) {
			return true
		}
	}

	return false
}
