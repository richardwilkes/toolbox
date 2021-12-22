// Copyright ©2016-2021 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package notifier

import (
	"sort"
	"strings"
	"sync"

	"github.com/richardwilkes/toolbox/errs"
)

// Target defines the method a target of notifications must implement.
type Target interface {
	// HandleNotification is called to deliver a notification.
	HandleNotification(name string, data, producer interface{})
}

// BatchTarget defines the methods a target of notifications that wants to be notified when a batch change occurs must
// implement.
type BatchTarget interface {
	Target
	// BatchMode is called both before and after a series of notifications are about to be broadcast. The target is not
	// guaranteed to have intervening calls to HandleNotification() made to it.
	BatchMode(start bool)
}

// Notifier tracks targets of notifications and provides methods for notifying them.
type Notifier struct {
	recoveryHandler errs.RecoveryHandler
	lock            sync.RWMutex
	batchTargets    map[BatchTarget]bool
	productionMap   map[string]map[Target]int
	nameMap         map[Target]map[string]bool
	currentBatch    []BatchTarget
	batchLevel      int
	enabled         bool
}

// New creates a new notifier.
func New(recoveryHandler errs.RecoveryHandler) *Notifier {
	return &Notifier{
		recoveryHandler: recoveryHandler,
		batchTargets:    make(map[BatchTarget]bool),
		productionMap:   make(map[string]map[Target]int),
		nameMap:         make(map[Target]map[string]bool),
		enabled:         true,
	}
}

// Register a target with this notifier. 'priority' is the relative notification priority, with higher values being
// delivered first. 'names' are the names the target wishes to consume. Names are hierarchical (separated by a .), so
// specifying a name of "foo.bar" will consume not only a produced name of "foo.bar", but all sub-names, such as
// "foo.bar.a", but not "foo.barn" or "foo.barn.a".
func (n *Notifier) Register(target Target, priority int, names ...string) {
	var normalizedNames []string
	for _, name := range names {
		name = normalizeName(name)
		if len(name) > 0 {
			normalizedNames = append(normalizedNames, name)
		}
	}
	if len(normalizedNames) > 0 {
		n.lock.Lock()
		targetNames, ok := n.nameMap[target]
		if !ok {
			targetNames = make(map[string]bool, len(normalizedNames))
			n.nameMap[target] = targetNames
		}
		if batchTarget, ok2 := target.(BatchTarget); ok2 {
			n.batchTargets[batchTarget] = true
		}
		for _, name := range normalizedNames {
			set, ok3 := n.productionMap[name]
			if !ok3 {
				set = make(map[Target]int)
				n.productionMap[name] = set
			}
			set[target] = priority
			targetNames[name] = true
		}
		n.lock.Unlock()
	}
}

func normalizeName(name string) string {
	var buffer strings.Builder
	for _, one := range strings.Split(name, ".") {
		if one != "" {
			if buffer.Len() > 0 {
				buffer.WriteByte('.')
			}
			buffer.WriteString(one)
		}
	}
	return buffer.String()
}

// RegisterFromNotifier adds all registrations from the other notifier into this notifier.
func (n *Notifier) RegisterFromNotifier(other *Notifier) {
	if n == other {
		return
	}
	// To avoid a potential deadlock, we make a copy of the other notifier's data first.
	other.lock.Lock()
	batchTargets := make(map[BatchTarget]bool, len(other.batchTargets))
	for k, v := range other.batchTargets {
		batchTargets[k] = v
	}
	productionMap := make(map[string]map[Target]int, len(other.productionMap))
	for k, v := range other.productionMap {
		m := make(map[Target]int, len(v))
		for k1, v1 := range v {
			m[k1] = v1
		}
		productionMap[k] = m
	}
	nameMap := make(map[Target]map[string]bool, len(other.nameMap))
	for k, v := range other.nameMap {
		m := make(map[string]bool, len(v))
		for k1, v1 := range v {
			m[k1] = v1
		}
		nameMap[k] = m
	}
	other.lock.Unlock()

	n.lock.Lock()
	for k, v := range batchTargets {
		n.batchTargets[k] = v
	}
	for k, v := range productionMap {
		if pm, ok := n.productionMap[k]; ok {
			for k1, v1 := range pm {
				pm[k1] = v1
			}
			n.productionMap[k] = pm
		} else {
			n.productionMap[k] = v
		}
	}
	for k, v := range nameMap {
		if nm, ok := n.nameMap[k]; ok {
			for k1, v1 := range nm {
				nm[k1] = v1
			}
			n.nameMap[k] = nm
		} else {
			n.nameMap[k] = v
		}
	}
	n.lock.Unlock()
}

// Unregister a target.
func (n *Notifier) Unregister(target Target) {
	n.lock.Lock()
	if nameMap, exists := n.nameMap[target]; exists {
		if batchTarget, ok := target.(BatchTarget); ok {
			delete(n.batchTargets, batchTarget)
		}
		for name := range nameMap {
			if set, ok := n.productionMap[name]; ok {
				delete(set, target)
				if len(set) == 0 {
					delete(n.productionMap, name)
				}
			}
		}
		delete(n.nameMap, target)
	}
	n.lock.Unlock()
}

// Enabled returns true if this notifier is currently enabled.
func (n *Notifier) Enabled() bool {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.enabled
}

// SetEnabled sets whether this notifier is enabled or not.
func (n *Notifier) SetEnabled(enabled bool) {
	n.lock.Lock()
	n.enabled = enabled
	n.lock.Unlock()
}

// Notify sends a notification to all interested targets.
func (n *Notifier) Notify(name string, producer interface{}) {
	n.NotifyWithData(name, nil, producer)
}

// NotifyWithData sends a notification to all interested targets. This is a synchronous notification and will not return
// until all interested targets handle the notification.
func (n *Notifier) NotifyWithData(name string, data, producer interface{}) {
	if n.Enabled() {
		name = normalizeName(name)
		if len(name) > 0 {
			targets := make(map[Target]int)
			names := strings.Split(name, ".")
			n.lock.RLock()
			if n.enabled {
				var buffer strings.Builder
				for _, one := range names {
					buffer.WriteString(one)
					one = buffer.String()
					if set, ok := n.productionMap[one]; ok {
						for k, v := range set {
							targets[k] = v
						}
					}
					buffer.WriteByte('.')
				}
			}
			n.lock.RUnlock()
			if len(targets) > 0 {
				list := make([]Target, 0, len(targets))
				for k := range targets {
					list = append(list, k)
				}
				sort.Slice(list, func(i, j int) bool {
					return targets[list[i]] > targets[list[j]]
				})
				for _, target := range list {
					n.notifyTarget(target, name, data, producer)
				}
			}
		}
	}
}

func (n *Notifier) notifyTarget(target Target, name string, data, producer interface{}) {
	defer errs.Recovery(n.recoveryHandler)
	target.HandleNotification(name, data, producer)
}

// BatchLevel returns the current batch level.
func (n *Notifier) BatchLevel() int {
	n.lock.RLock()
	defer n.lock.RUnlock()
	return n.batchLevel
}

// StartBatch informs all BatchTargets that a batch of notifications will be starting. If a previous call to this method
// was made without a call to EndBatch(), then the batch level will be incremented, but no notifications will be made.
func (n *Notifier) StartBatch() {
	var targets []BatchTarget
	n.lock.Lock()
	if n.enabled {
		n.batchLevel++
		if n.batchLevel == 1 && len(n.batchTargets) > 0 {
			n.currentBatch = make([]BatchTarget, 0, len(n.batchTargets))
			for k := range n.batchTargets {
				n.currentBatch = append(n.currentBatch, k)
			}
			targets = n.currentBatch
		}
	}
	n.lock.Unlock()
	for _, target := range targets {
		n.notifyBatchTarget(target, true)
	}
}

func (n *Notifier) notifyBatchTarget(target BatchTarget, start bool) {
	defer errs.Recovery(n.recoveryHandler)
	target.BatchMode(start)
}

// EndBatch informs all BatchTargets that were present when StartBatch() was called that a batch of notifications just
// finished. If batch level is still greater than zero after being decremented, then no notifications will be made.
func (n *Notifier) EndBatch() {
	var targets []BatchTarget
	n.lock.Lock()
	if n.enabled && n.batchLevel > 0 {
		n.batchLevel--
		if n.batchLevel == 0 {
			targets = n.currentBatch
			n.currentBatch = nil
		}
	}
	n.lock.Unlock()
	for _, target := range targets {
		n.notifyBatchTarget(target, false)
	}
}

// Reset removes all targets.
func (n *Notifier) Reset() {
	n.lock.Lock()
	n.batchTargets = make(map[BatchTarget]bool)
	n.productionMap = make(map[string]map[Target]int)
	n.nameMap = make(map[Target]map[string]bool)
	n.currentBatch = nil
	n.batchLevel = 0
	n.lock.Unlock()
}
