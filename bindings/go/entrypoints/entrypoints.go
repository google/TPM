// Package entrypoints contains bindings to the TPM's public interface.
//
// This package provides the direct entry points to the TPM 2.0 reference code,
// including initialization, manufacturing, and executing commands.
//
// Before using these entry points, you must ensure that the [platform] package
// has been initialized by setting [platform.Current].
//
// Typical usage involves calling [Manufacture] and [Init], and then
// using [ExecuteCommand] to send raw byte commands to the TPM.
package entrypoints

// #cgo CFLAGS: -std=c11 -Wall -Wextra -pedantic
//
// #cgo CFLAGS: -I ../../../TPMCmd/tpm/include/
// #cgo CFLAGS: -I ../../../TPMCmd/TpmConfiguration/
//
// // TODO: Fix headers to not need extra #include files
// #include <TpmConfiguration/TpmProfile.h>
// #include <tpm_public/BaseTypes.h>
// #include <platform_interface/platform_to_tpm_interface.h>
import "C"

import (
	"errors"
	"iter"
	"runtime"
	"unsafe"

	// Depend on the internal package to force a build of the TPMCmd/tpm/ code.
	_ "github.com/google/TPM/bindings/c"
)

// Error values for [HashSequence]
var (
	ErrHashStart = errors.New("call to _TPM_Hash_Start() failed")
	ErrHashData  = errors.New("call to _TPM_Hash_Data() failed")
	ErrHashEnd   = errors.New("call to _TPM_Hash_End() failed")
)

// Error values for [Manufacture]
var (
	ErrManufactureAlreadyDone   = errors.New("manufacturing is already complete")
	ErrManufactureInvalidConfig = errors.New("manufacturing failed: invalid config")
	ErrManufactureNvNotReady    = errors.New("manufacturing failed: NV System not available")
	ErrManufactureUnknown       = errors.New("manufacturing failed: unknown error")
)

// Send a _TPM_Init indication.
func Init() {
	C._TPM_Init()
}

// Execute a command and return a response
//
// Caller must ensure that:
//
//	len(cmd) <= math.MaxUint32
func ExecuteCommand(cmd []byte) []byte {
	cmdSize := C.uint32_t(len(cmd))
	cmdData := (*C.uint8_t)(unsafe.SliceData(cmd))

	// ExecuteCommand requires the response buffer to be pre-allocated
	rsp := make([]byte, C.MAX_RESPONSE_SIZE)
	rspSize := C.uint32_t(len(rsp))
	rspData := (*C.uint8_t)(unsafe.SliceData(rsp))
	// Pin the rsp buffer so it isn't moved during [ExecuteCommand]
	var pinner runtime.Pinner
	pinner.Pin(rspData)
	defer pinner.Unpin()

	C.ExecuteCommand(cmdSize, cmdData, &rspSize, &rspData)

	// If ExecuteCommand wrote to our response buffer, just slice it.
	if rspData == (*C.uint8_t)(unsafe.SliceData(rsp)) {
		return rsp[:rspSize]
	}
	// If ExecuteCommand used its own response buffer, copy it into Go memory.
	return C.GoBytes(unsafe.Pointer(rspData), C.int(rspSize))
}

// HashSequence sends an H-CRTM indication to the TPM.
//
// [Init] must be called before this function.
func HashSequence(data iter.Seq[[]byte]) error {
	if C._TPM_Hash_Start() != C.TRUE {
		return ErrHashStart
	}
	for buf := range data {
		dataSize := C.uint32_t(len(buf))
		dataPointer := (*C.uint8_t)(unsafe.SliceData(buf))
		if C._TPM_Hash_Data(dataSize, dataPointer) != C.TRUE {
			return ErrHashData
		}
	}
	if C._TPM_Hash_End() != C.TRUE {
		return ErrHashEnd
	}
	return nil
}

// Manufacture a TPM
func Manufacture() error {
	switch C.TPM_Manufacture(1) {
	case 0:
		return nil
	case 1:
		return ErrManufactureAlreadyDone
	case -1:
		return ErrManufactureInvalidConfig
	case -2:
		return ErrManufactureNvNotReady
	default:
		return ErrManufactureUnknown
	}
}
