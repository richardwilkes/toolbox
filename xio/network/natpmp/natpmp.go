// Copyright (c) 2016-2024 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

// Package natpmp provides an implementation of NAT-PMP.
// See https://tools.ietf.org/html/rfc6886
package natpmp

import (
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"

	"github.com/jackpal/gateway"
	"github.com/richardwilkes/toolbox/atexit"
	"github.com/richardwilkes/toolbox/errs"
)

const (
	protocolVersion = 0
	expiration      = uint32(time.Hour / time.Second)
	tcpFlag         = 0x10000
)

const (
	opExternalAddress = iota
	opMapUDP
	opMapTCP
)

type mapping struct {
	notifyChan chan any
	renew      time.Time
	external   int
}

var (
	once     sync.Once
	gw       net.IP
	lock     sync.RWMutex
	mappings = make(map[int]mapping)
)

// ExternalAddress returns the external address the internet sees you as having.
func ExternalAddress() (net.IP, error) {
	buffer := make([]byte, 2)
	buffer[0] = protocolVersion
	buffer[1] = opExternalAddress
	response, err := call(buffer, 12)
	if err != nil {
		return nil, err
	}
	return response[8:12], nil
}

// MapTCP maps the specified TCP port for external access. It returns the port on the external address that can be used
// to connect to the internal port. If you wish to be notified of changes to the external port mapping, provide a notify
// channel. It will be sent an int containing the updated external port mapping when it changes or an error if a renewal
// fails. The channel will only be sent to if it is ready.
func MapTCP(port int, notifyChan chan any) (int, error) {
	if err := checkPort(port); err != nil {
		return 0, err
	}
	response, err := call(makeMapBuffer(opMapTCP, uint16(port)), 16)
	if err != nil {
		return 0, err
	}
	external := int(binary.BigEndian.Uint16(response[10:12]))
	addMapping(port|tcpFlag, external, notifyChan)
	return external, nil
}

// MapUDP maps the specified UDP port for external access. It returns the port on the external address that can be used
// to connect to the internal port. If you wish to be notified of changes to the external port mapping, provide a notify
// channel. It will be sent an int containing the updated external port mapping when it changes or an error if a renewal
// fails. The channel will only be sent to if it is ready.
func MapUDP(port int, notifyChan chan any) (int, error) {
	if err := checkPort(port); err != nil {
		return 0, err
	}
	response, err := call(makeMapBuffer(opMapUDP, uint16(port)), 16)
	if err != nil {
		return 0, err
	}
	external := int(binary.BigEndian.Uint16(response[10:12]))
	addMapping(port, external, notifyChan)
	return external, nil
}

// UnmapTCP unmaps a previously mapped internal TCP port.
func UnmapTCP(port int) error {
	if err := checkPort(port); err != nil {
		return err
	}
	_, err := call(makeUnmapBuffer(opMapTCP, uint16(port)), 16)
	if err != nil {
		return err
	}
	removeMapping(port | tcpFlag)
	return err
}

// UnmapUDP unmaps a previously mapped internal UDP port.
func UnmapUDP(port int) error {
	if err := checkPort(port); err != nil {
		return err
	}
	_, err := call(makeUnmapBuffer(opMapUDP, uint16(port)), 16)
	if err != nil {
		return err
	}
	removeMapping(port)
	return err
}

func checkPort(port int) error {
	if port > 0 && port < 65536 {
		return nil
	}
	return errs.Newf("port (%d) must be in the range 1-65535", port)
}

func makeMapBuffer(op byte, port uint16) []byte {
	buffer := makeUnmapBuffer(op, port)
	binary.BigEndian.PutUint16(buffer[6:8], port)
	binary.BigEndian.PutUint32(buffer[8:12], expiration)
	return buffer
}

func makeUnmapBuffer(op byte, port uint16) []byte {
	buffer := make([]byte, 12)
	buffer[0] = protocolVersion
	buffer[1] = op
	binary.BigEndian.PutUint16(buffer[4:6], port)
	return buffer
}

func addMapping(internal, external int, notifyChan chan any) {
	lock.Lock()
	mappings[internal] = mapping{
		external:   external,
		renew:      time.Now().Add(50 * time.Minute),
		notifyChan: notifyChan,
	}
	lock.Unlock()
}

func removeMapping(internal int) {
	lock.Lock()
	delete(mappings, internal)
	lock.Unlock()
}

func call(msg []byte, resultSize int) ([]byte, error) {
	once.Do(setupGateway)
	if gw == nil {
		return nil, errs.New("No gateway found")
	}
	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   gw,
		Port: 5351,
	})
	if err != nil {
		return nil, errs.Wrap(err)
	}
	defer func() { err = conn.Close() }()
	timeout := time.Now().Add(30 * time.Second)
	err = conn.SetDeadline(timeout)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	result := make([]byte, resultSize)
	for time.Now().Before(timeout) {
		if _, err = conn.Write(msg); err != nil {
			return nil, errs.Wrap(err)
		}
		var n int
		var remote *net.UDPAddr
		if n, remote, err = conn.ReadFromUDP(result); err != nil {
			var nerr net.Error
			if errors.As(err, &nerr) && nerr.Timeout() {
				continue
			}
			return nil, errs.Wrap(err)
		}
		if !remote.IP.Equal(gw) {
			continue
		}
		if n != resultSize {
			return nil, errs.Newf("unexpected result size (received %d, expected %d)", n, resultSize)
		}
		if result[0] != 0 {
			return nil, errs.Newf("unknown protocol version (%d)", result[0])
		}
		expectedOp := msg[1] | 0x80
		if result[1] != expectedOp {
			return nil, errs.Newf("unexpected opcode (received %d, expected %d)", result[1], expectedOp)
		}
		code := binary.BigEndian.Uint16(result[2:4])
		switch code {
		case 0:
			return result, nil
		case 1:
			return nil, errs.New("unsupported version")
		case 2:
			return nil, errs.New("not authorized")
		case 3:
			return nil, errs.New("network failure")
		case 4:
			return nil, errs.New("out of resources")
		case 5:
			return nil, errs.New("unsupported opcode")
		default:
			return nil, errs.Newf("unknown result code %d", code)
		}
	}
	return nil, errs.Newf("timed out trying to contact gateway")
}

func setupGateway() {
	var err error
	if gw, err = gateway.DiscoverGateway(); err == nil {
		atexit.Register(cleanup)
		go renewals()
	}
}

func renewals() {
	for {
		time.Sleep(time.Minute)
		now := time.Now()
		lock.RLock()
		renew := make(map[int]mapping, len(mappings))
		for k, v := range mappings {
			if !now.Before(v.renew) {
				renew[k] = v
			}
		}
		lock.RUnlock()
		for port, v := range renew {
			var external int
			var err error
			if port&tcpFlag != 0 {
				external, err = MapTCP(port&^tcpFlag, v.notifyChan)
			} else {
				external, err = MapUDP(port, v.notifyChan)
			}
			if v.notifyChan != nil {
				if err != nil {
					var portType string
					if port&tcpFlag != 0 {
						portType = "TCP"
					} else {
						portType = "UDP"
					}
					select {
					case v.notifyChan <- errs.NewWithCausef(err, "mapping renewal for %s port %d failed", portType, port):
					default:
					}
				} else if v.external != external {
					select {
					case v.notifyChan <- external:
					default:
					}
				}
			}
		}
	}
}

func cleanup() {
	lock.RLock()
	ports := make([]int, len(mappings))
	for port := range mappings {
		ports = append(ports, port)
	}
	lock.RUnlock()
	for _, port := range ports {
		var err error
		if port&tcpFlag != 0 {
			err = UnmapTCP(port &^ tcpFlag)
		} else {
			err = UnmapUDP(port)
		}
		if err != nil {
			errs.Log(err)
		}
	}
}
