// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package network provides network-related utilities.
package network

import (
	"math/rand/v2"
	"net"
	"sort"
	"strings"
	"time"

	"github.com/richardwilkes/toolbox/collection"
	"github.com/richardwilkes/toolbox/txt"
)

// Constants for common network addresses.
const (
	IPv4LoopbackAddress = "127.0.0.1"
	IPv6LoopbackAddress = "::1"
	LocalHost           = "localhost"
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
		//nolint:gosec // Yes, it is ok to use a weak prng here
		time.Sleep(time.Duration(100+rand.IntN(50)) * time.Millisecond)
	}
	return IPv4LoopbackAddress
}

// PrimaryAddress returns the primary hostname and its associated IP address and MAC address.
func PrimaryAddress() (hostname, ipAddress, macAddress string) {
	// Try up to 3 times in case of transient errors
	for i := 0; i < 3; i++ {
		lowest := 1000000
		for address, iFace := range ActiveAddresses() {
			if iFace.Index < lowest {
				lowest = iFace.Index
				hostname = address
				macAddress = iFace.HardwareAddr.String()
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
		//nolint:gosec // Yes, it is ok to use a weak prng here
		time.Sleep(time.Duration(100+rand.IntN(50)) * time.Millisecond)
	}
	return LocalHost, IPv4LoopbackAddress, "00:00:00:00:00:00"
}

// ActiveAddresses determines the best address for each active network interface. IPv4 addresses will be selected over
// IPv6 addresses on the same interface. Numeric addresses are resolved into names where possible.
func ActiveAddresses() map[string]net.Interface {
	result := make(map[string]net.Interface)
	if iFaces, err := net.Interfaces(); err == nil {
		for _, iFace := range iFaces {
			const interesting = net.FlagUp | net.FlagBroadcast
			if iFace.Flags&interesting == interesting {
				if name := Address(iFace); name != "" {
					result[name] = iFace
				}
			}
		}
	}
	return result
}

// Address returns the best address for the network interface. IPv4 addresses will be selected over IPv6 addresses on
// the same interface. Numeric addresses are resolved into names where possible. An empty string will be returned if the
// network interface cannot be resolved into an IPv4 or IPv6 address.
func Address(iFace net.Interface) string {
	if addrs, err := iFace.Addrs(); err == nil {
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
				var names []string
				if names, err = net.LookupAddr(ipAddr); err == nil {
					if len(names) > 0 {
						name := strings.TrimSuffix(names[0], ".")
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

// AddressesForHost returns the addresses/names for the given host. If an IP number is passed in, then it will be
// returned. If a host name is passed in, the host name plus the IP address(es) it resolves to will be returned. If the
// empty string is passed in, then the host names and IP addresses for all active interfaces will be returned.
func AddressesForHost(host string) []string {
	ss := collection.NewSet[string]()
	if host == "" { // All address on machine
		if iFaces, err := net.Interfaces(); err == nil {
			for _, iFace := range iFaces {
				const interesting = net.FlagUp | net.FlagBroadcast
				if iFace.Flags&interesting == interesting {
					var addrs []net.Addr
					if addrs, err = iFace.Addrs(); err == nil {
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
								ss.Add(ip.String())
								var names []string
								if names, err = net.LookupAddr(ip.String()); err == nil {
									for _, name := range names {
										ss.Add(strings.TrimSuffix(name, "."))
									}
								}
							}
						}
					}
				}
			}
		}
	} else {
		ss.Add(host)
		if net.ParseIP(host) == nil {
			if ips, err := net.LookupIP(host); err == nil && len(ips) > 0 {
				for _, ip := range ips {
					ss.Add(ip.String())
				}
			}
		}
	}
	for _, one := range []string{"::", IPv6LoopbackAddress, IPv4LoopbackAddress} {
		if ss.Contains(one) {
			delete(ss, one)
			ss.Add(LocalHost)
		}
	}
	addrs := ss.Values()
	sort.Slice(addrs, func(i, j int) bool {
		isName1 := net.ParseIP(addrs[i]) == nil
		isName2 := net.ParseIP(addrs[j]) == nil
		if isName1 == isName2 {
			return txt.NaturalLess(addrs[i], addrs[j], true)
		}
		return isName1
	})
	return addrs
}
