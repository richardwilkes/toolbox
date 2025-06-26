package gateway

import (
	"net"
	"syscall"
	"unsafe"

	"github.com/richardwilkes/toolbox/errs"
)

// Default returns the IP of the default gateway for the current machine, or nil if no gateway is found.
func Default() net.IP {
	type SOCKADDR_INET struct {
		IP [16]byte
		_  [12]byte
	}
	// type SOCKADDR_INET struct {
	// 	Sin6_family int16
	// 	Sin6_port   uint16
	// 	IPv4        [4]byte
	// 	IPv6        [16]byte
	// 	Scope_id    uint32
	// }
	type IP_ADDRESS_PREFIX struct {
		Prefix       SOCKADDR_INET
		PrefixLength uint8
		_            [3]byte
	}
	type MIB_IPFORWARD_ROW2 struct {
		InterfaceLuid        uint64
		InterfaceIndex       uint32
		DestinationPrefix    IP_ADDRESS_PREFIX
		NextHop              SOCKADDR_INET
		SitePrefixLength     uint8
		_                    [3]byte
		ValidLifetime        uint32
		PreferredLifetime    uint32
		Metric               uint32
		Protocol             uint32
		Loopback             uint8
		AutoconfigureAddress uint8
		Publish              uint8
		Immortal             uint8
		Age                  uint32
		Origin               uint32
	}
	type MIB_IPFORWARD_TABLE2 struct {
		NumEntries uint32
		// Followed by NumEntries MIB_IPFORWARD_ROW2
	}
	iphlpapi := syscall.NewLazyDLL("iphlpapi.dll")
	getIpForwardTable2 := iphlpapi.NewProc("GetIpForwardTable2")
	var table *MIB_IPFORWARD_TABLE2
	r1, _, _ := getIpForwardTable2.Call(syscall.AF_UNSPEC, uintptr(unsafe.Pointer(&table)))
	if r1 != 0 || table == nil {
		errs.Log(errs.New("unable to get default routes"))
		return nil
	}
	defer func() {
		iphlpapi.NewProc("FreeMibTable").Call(uintptr(unsafe.Pointer(table)))
	}()
	type routeInfo struct {
		metric int
		gw     net.IP
	}
	var ip net.IP
	best := uint32(0xFFFFFFFF)
	n := table.NumEntries
	rows := (*[1 << 12]MIB_IPFORWARD_ROW2)(unsafe.Pointer(uintptr(unsafe.Pointer(table)) + unsafe.Sizeof(*table)))[:n:n]
	for _, row := range rows {
		if row.SitePrefixLength == 0 {
			if row.Metric < best {
				best = row.Metric
				ip = make(net.IP, net.IPv6len)
				copy(ip, row.NextHop.IP[:16])
			}
		}
	}
	if ip != nil {
		return ip
	}
	return nil
}
