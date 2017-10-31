package network

import (
	"math/rand"
	"net"
	"strings"
	"time"
)

// PrimaryIPAddress returns the primary IP address.
func PrimaryIPAddress() string {
	// Try up to 3 times in case of transient errors
	for i := 0; i < 3; i++ {
		if addresses, err := net.InterfaceAddrs(); err == nil {
			var fallback string
			for _, address := range addresses {
				var ip net.IP
				switch v := address.(type) {
				case *net.IPNet:
					ip = v.IP
				case *net.IPAddr:
					ip = v.IP
				default:
					continue
				}
				if ip.IsGlobalUnicast() {
					if ip.To4() != nil {
						return ip.String()
					}
					if fallback == "" {
						fallback = ip.String()
					}
				}
			}
			if fallback != "" {
				return fallback
			}
		}
		time.Sleep(time.Duration(100+rand.Intn(50)) * time.Millisecond)
	}
	return "127.0.0.1"
}

// PrimaryAddress returns the primary hostname and its associated IP address
// and MAC address.
func PrimaryAddress() (hostname, ipAddress, macAddress string) {
	// Try up to 3 times in case of transient errors
	for i := 0; i < 3; i++ {
		lowest := 1000000
		for address, iface := range ActiveAddresses() {
			if iface.Index < lowest {
				lowest = iface.Index
				hostname = address
				macAddress = iface.HardwareAddr.String()
			}
		}
		if hostname != "" {
			if ips, err := net.LookupIP(hostname); err == nil && len(ips) > 0 {
				for _, ip := range ips {
					if ip.To4() != nil {
						ipAddress = ip.String()
						break
					} else if ipAddress == "" {
						ipAddress = ip.String()
					}
				}
				if ipAddress != "" {
					return hostname, ipAddress, macAddress
				}
			}
		}
		time.Sleep(time.Duration(100+rand.Intn(50)) * time.Millisecond)
	}
	return "localhost", "127.0.0.1", "00:00:00:00:00:00"
}

// ActiveAddresses determines the best address for each active network
// interface. IPv4 addresses will be selected over IPv6 addresses on the same
// interface. Numeric addresses are resolved into names where possible.
func ActiveAddresses() map[string]net.Interface {
	result := make(map[string]net.Interface)
	if ifaces, err := net.Interfaces(); err == nil {
		for _, iface := range ifaces {
			const interesting = net.FlagUp | net.FlagBroadcast
			if iface.Flags&interesting == interesting {
				if name := Address(iface); name != "" {
					result[name] = iface
				}
			}
		}
	}
	return result
}

// Address returns the best address for the network interface. IPv4 addresses
// will be selected over IPv6 addresses on the same interface. Numeric
// addresses are resolved into names where possible. An empty string will be
// returned if the network interface cannot be resolved into an IPv4 or IPv6
// address.
func Address(iface net.Interface) string {
	if addrs, err := iface.Addrs(); err == nil {
		var fallback string
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			if ip.IsGlobalUnicast() {
				ipAddr := ip.String()
				if names, err := net.LookupAddr(ipAddr); err == nil {
					if len(names) > 0 {
						name := names[0]
						if strings.HasSuffix(name, ".") {
							name = name[:len(name)-1]
						}
						if ip.To4() != nil {
							return name
						}
						if fallback == "" {
							fallback = name
						}
						continue
					}
				}
				if ip.To4() != nil {
					return ipAddr
				}
				if fallback == "" {
					fallback = ipAddr
				}
			}
		}
		if fallback != "" {
			return fallback
		}
	}
	return ""
}
