package utils

import (
	"fmt"
	"net"
	"strings"
)

// PrintServerAddresses 启动成功后打印本机与局域网访问地址，便于复制到客户端
func PrintServerAddresses(port int) {
	fmt.Println()
	fmt.Println("========================================")
	fmt.Println("  Linyu server started")
	fmt.Println("  Copy an address below to the client:")
	fmt.Println("----------------------------------------")
	fmt.Printf("  http://127.0.0.1:%d\n", port)
	for _, ip := range LocalIPv4s() {
		fmt.Printf("  http://%s:%d\n", ip, port)
	}
	fmt.Println("========================================")
	fmt.Println()
}

// LocalIPv4s 返回本机可用局域网 IPv4（跳过回环、链路本地与常见虚拟网卡）
func LocalIPv4s() []string {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	seen := make(map[string]struct{})
	var ips []string
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 跳过常见虚拟网卡，优先展示真实局域网地址
		name := strings.ToLower(iface.Name)
		if strings.Contains(name, "docker") ||
			strings.Contains(name, "veth") ||
			strings.Contains(name, "br-") ||
			strings.Contains(name, "vmware") ||
			strings.Contains(name, "virtualbox") ||
			strings.Contains(name, "hyper-v") {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil {
				continue
			}
			ip4 := ip.To4()
			if ip4 == nil || ip4.IsLoopback() || ip4.IsLinkLocalUnicast() {
				continue
			}
			s := ip4.String()
			if _, ok := seen[s]; ok {
				continue
			}
			seen[s] = struct{}{}
			ips = append(ips, s)
		}
	}
	return ips
}
