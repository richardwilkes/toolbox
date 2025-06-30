// Copyright (c) 2016-2025 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package gateway

import (
	"bytes"
	"encoding/binary"
	"net"
	"os"
	"syscall"
	"unsafe"

	"github.com/richardwilkes/toolbox/v2/errs"
)

type routeInfo struct {
	gw     net.IP
	metric int
}

// Default returns the IP of the default gateway for the current machine, or nil if no gateway is found.
func Default() net.IP {
	var best *routeInfo
	if v4Routes, err := getRoutes(syscall.AF_INET); err != nil {
		errs.Log(errs.NewWithCause("unable to get default routes for IPv4", err))
	} else {
		for _, r := range v4Routes {
			if best == nil || r.metric < best.metric {
				best = &r
			}
		}
	}
	if v6Routes, err := getRoutes(syscall.AF_INET6); err != nil {
		errs.Log(errs.NewWithCause("unable to get default routes for IPv6", err))
	} else {
		for _, r := range v6Routes {
			if best == nil || r.metric < best.metric {
				best = &r
			}
		}
	}
	if best == nil {
		return nil
	}
	return best.gw
}

func getRoutes(family int) ([]routeInfo, error) {
	fd, err := syscall.Socket(syscall.AF_NETLINK, syscall.SOCK_RAW, syscall.NETLINK_ROUTE)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer func() {
		_ = syscall.Close(fd) //nolint:errcheck // Ignore close error
	}()
	var req bytes.Buffer
	hdr := syscall.NlMsghdr{
		Len:   uint32(syscall.NLMSG_HDRLEN + syscall.SizeofRtMsg),
		Type:  syscall.RTM_GETROUTE,
		Flags: syscall.NLM_F_DUMP | syscall.NLM_F_REQUEST,
		Seq:   1,
		Pid:   uint32(os.Getpid()),
	}
	_ = binary.Write(&req, binary.LittleEndian, hdr)                                  //nolint:errcheck // Can't fail
	_ = binary.Write(&req, binary.LittleEndian, syscall.RtMsg{Family: uint8(family)}) //nolint:errcheck // Can't fail
	sa := &syscall.SockaddrNetlink{Family: syscall.AF_NETLINK}
	if err = syscall.Sendto(fd, req.Bytes(), 0, sa); err != nil {
		return nil, errs.Wrap(err)
	}
	buf := make([]byte, 16384)
	var n int
	if n, _, err = syscall.Recvfrom(fd, buf, 0); err != nil {
		return nil, errs.Wrap(err)
	}
	var msgs []syscall.NetlinkMessage
	if msgs, err = syscall.ParseNetlinkMessage(buf[:n]); err != nil {
		return nil, errs.Wrap(err)
	}
	var routes []routeInfo
	for _, m := range msgs {
		if m.Header.Type == syscall.NLMSG_DONE {
			break
		}
		if m.Header.Type != syscall.RTM_NEWROUTE {
			continue
		}
		rt := (*syscall.RtMsg)(unsafe.Pointer(&m.Data[0]))
		if int(rt.Family) != family {
			continue
		}
		attrs := m.Data[syscall.SizeofRtMsg:]
		dstLen := rt.Dst_len
		var gw net.IP
		metric := 0
		for len(attrs) >= 4 {
			attr := (*syscall.RtAttr)(unsafe.Pointer(&attrs[0]))
			if attr.Len < 4 || int(attr.Len) > len(attrs) {
				break
			}
			value := attrs[4:attr.Len]
			switch attr.Type {
			case syscall.RTA_DST:
				dstLen = rt.Dst_len
			case syscall.RTA_GATEWAY:
				if family == syscall.AF_INET && len(value) >= 4 {
					gw = net.IPv4(value[0], value[1], value[2], value[3])
				} else if family == syscall.AF_INET6 && len(value) >= 16 {
					gw = make(net.IP, net.IPv6len)
					copy(gw, value[:16])
				}
			case syscall.RTA_PRIORITY:
				if len(value) >= 4 {
					metric = int(binary.LittleEndian.Uint32(value[:4]))
				}
			}
			attrs = attrs[attr.Len:]
		}
		if dstLen == 0 && gw != nil {
			routes = append(routes, routeInfo{gw: gw, metric: metric})
		}
	}
	return routes, nil
}
