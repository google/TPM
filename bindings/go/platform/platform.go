// Package platform allows Go code to implement the TPM 2.0 Platform Interface.
//
// The TPM 2.0 reference code requires a platform implementation to handle
// non-volatile storage, entropy, timers, and system controls (like locality).
// This package exposes these requirements as Go interfaces and structs that
// can be implemented or configured by the caller.
//
// To use this package, you must create a [Platform] struct and assign it to
// [Current]. If [Current] is not set when the TPM code is run, the callbacks
// will panic.
package platform

// #cgo CFLAGS: -std=c11 -Wall -Wextra -pedantic
//
// #cgo CFLAGS: -I ../../../TPMCmd/tpm/include/
// #cgo CFLAGS: -I ../../../TPMCmd/TpmConfiguration/
import "C"

import (
	"io"
	"math"
	"unsafe"
)

// Variable for the currently active platfom. This is initially nil. If it is
// still nil when the TPM code is run, certain platform functions will panic.
var Current *Platform = nil

// Platform defines the TPM 2.0 Platform Interface specified in:
//
//	tpm/include/platform_interface/tpm_to_platform_interface.h
type Platform struct {
	Properties      Properties                           // Platform properties
	Firmware        Firmware                             // Firmware interface
	Entropy         io.Reader                            // Source of randomness
	Timer           Timer                                // Monotonic timer
	Storage         Storage                              // Non-volatile storage
	System          System                               // System controls (locality, cancel)
	Failure         Failure                              // Failure mode handler
	Debug           io.Writer                            // Debug log destination
	ManufactureData func(data []byte)                    // Callback to get platform manufacture data
	SelfTestEnabled func(fullTest bool, alg uint16) bool // Callback to check if self test is enabled
	PCRConfig       PCRConfig                            // PCR configuration
	ACT             ACT                                  // Authenticated Countdown Timer
}

/*** Entropy ***/

//export _plat__GetEntropy
func _plat__GetEntropy(entropy *uint8, amount uint32) int32 {
	if Current.Entropy == nil {
		panic("_plat__GetEntropy called with Current.Entropy == nil")
	}
	slice := unsafe.Slice(entropy, amount)
	n, err := Current.Entropy.Read(slice)
	if err != nil {
		// Negative return value from ReadEntropy indicates failure.
		return -1
	}
	return int32(min(n, math.MaxInt32))
}

//*** ManufactureData ***/

//export _plat__GetPlatformManufactureData
func _plat__GetPlatformManufactureData(buf *uint8, bufSize uint32) {
	if Current.ManufactureData == nil {
		return
	}
	data := unsafe.Slice(buf, bufSize)
	Current.ManufactureData(data)
}

/*** EnableSelfTests ***/

//export _plat_GetEnabledSelfTest
func _plat_GetEnabledSelfTest(fullTest uint8, pToTestVector *uint8, toTestVectorSize C.size_t) {
	if Current.SelfTestEnabled == nil {
		return
	}
	isFullTest := fullTest != 0

	toTestBits := unsafe.Slice(pToTestVector, toTestVectorSize)
	for i := 0; i < 8*len(toTestBits); i++ {
		byteIdx, bitIdx := i/8, i%8
		mask := byte(1) << bitIdx
		if toTestBits[byteIdx]&mask == 0 {
			continue
		}

		algID := uint16(i)
		if !Current.SelfTestEnabled(isFullTest, algID) {
			toTestBits[byteIdx] &= ^mask
		}
	}
}
