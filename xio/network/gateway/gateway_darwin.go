package gateway

import (
	"net"
	"syscall"

	"github.com/richardwilkes/toolbox/errs"
	"golang.org/x/net/route"
)

// Default returns the IP of the default gateway for the current machine, or nil if no gateway is found.
func Default() net.IP {
	rib, err := route.FetchRIB(syscall.AF_INET, route.RIBTypeRoute, 0)
	if err != nil {
		errs.Log(errs.NewWithCause("failed to fetch RIB", err))
		return nil
	}
	var msgs []route.Message
	if msgs, err = route.ParseRIB(route.RIBTypeRoute, rib); err != nil {
		errs.Log(errs.NewWithCause("failed to parse RIB", err))
		return nil
	}
	for _, msg := range msgs {
		if m, ok := msg.(*route.RouteMessage); ok {
			var ip net.IP
			switch sa := m.Addrs[syscall.RTAX_GATEWAY].(type) {
			case *route.Inet4Addr:
				return net.IPv4(sa.IP[0], sa.IP[1], sa.IP[2], sa.IP[3])
			case *route.Inet6Addr:
				ip = make(net.IP, net.IPv6len)
				copy(ip, sa.IP[:])
				return ip
			}
		}
	}
	return nil
}
