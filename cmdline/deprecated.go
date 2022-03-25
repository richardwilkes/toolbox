// Copyright ©2016-2022 by Richard A. Wilkes. All rights reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with
// this file, You can obtain one at http://mozilla.org/MPL/2.0/.
//
// This Source Code Form is "Incompatible With Secondary Licenses", as
// defined by the Mozilla Public License, version 2.0.

package cmdline

import "time"

// NewBoolOption creates a new bool Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewBoolOption(val *bool) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewIntOption creates a new int Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewIntOption(val *int) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewIntArrayOption creates a new []int Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewIntArrayOption(val *[]int) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewBoolArrayOption creates a new []bool Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewBoolArrayOption(val *[]bool) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewInt8Option creates a new int8 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewInt8Option(val *int8) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewInt8ArrayOption creates a new []int8 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewInt8ArrayOption(val *[]int8) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewInt16Option creates a new int16 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewInt16Option(val *int16) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewInt16ArrayOption creates a new []int16 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewInt16ArrayOption(val *[]int16) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewInt32Option creates a new int32 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewInt32Option(val *int32) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewInt32ArrayOption creates a new []int32 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewInt32ArrayOption(val *[]int32) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewInt64Option creates a new int64 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewInt64Option(val *int64) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewInt64ArrayOption creates a new []int64 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewInt64ArrayOption(val *[]int64) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUintOption creates a new uint Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUintOption(val *uint32) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUintArrayOption creates a new []uint Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUintArrayOption(val *[]uint32) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUint8Option creates a new uint8 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUint8Option(val *uint8) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUint8ArrayOption creates a new []uint8 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUint8ArrayOption(val *[]uint8) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUint16Option creates a new uint16 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUint16Option(val *uint16) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUint16ArrayOption creates a new []uint16 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUint16ArrayOption(val *[]uint16) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUint32Option creates a new uint32 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUint32Option(val *uint32) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUint32ArrayOption creates a new []uint32 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUint32ArrayOption(val *[]uint32) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUint64Option creates a new uint64 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUint64Option(val *uint64) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewUint64ArrayOption creates a new []uint64 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewUint64ArrayOption(val *[]uint64) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewFloat32Option creates a new float32 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewFloat32Option(val *float32) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewFloat32ArrayOption creates a new []float32 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewFloat32ArrayOption(val *[]float32) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewFloat64Option creates a new float64 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewFloat64Option(val *float64) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewFloat64ArrayOption creates a new []float64 Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewFloat64ArrayOption(val *[]float64) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewStringOption creates a new string Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewStringOption(val *string) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewStringArrayOption creates a new []string Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewStringArrayOption(val *[]string) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewDurationOption creates a new time.Duration Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewDurationOption(val *time.Duration) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}

// NewDurationArrayOption creates a new []time.Duration Option and attaches it to this CmdLine.
// Deprecated: Use .NewGeneralOption() instead. March 25, 2022
func (cl *CmdLine) NewDurationArrayOption(val *[]time.Duration) *Option {
	return cl.NewOption(&GeneralValue{Value: val})
}
