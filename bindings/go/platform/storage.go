package platform

// #include <platform_interface/tpm_to_platform_interface.h>
import "C"

import (
	"slices"
	"unsafe"
)

// Storage handles the platform's storage of non-volatile data.
//
// The [Storage.OnEnable], [Storage.OnReady], and [Storage.OnCommit] handlers
// are optional. If not provided, they default to successful operations.
type Storage struct {
	// Data holds the in-memory representation of the NV data.
	Data [C.NV_MEMORY_SIZE]byte

	// OnEnable loads the non-volatile data into memory.
	//
	// Returns 0 on success, >0 on recoverable failure, <0 on unrecoverable failure.
	OnEnable func() int

	// OnReady returns the state of the non-volatile storage system.
	//
	// This allows handling of potential NV errors before a command runs.
	OnReady func() ReadyState

	// OnCommit persists the in-memory data to non-volatile storage.
	//
	// Returns 0 on success, !=0 on failure.
	OnCommit func() int
}

type ReadyState int

const (
	NvReady        ReadyState = 0
	NvWriteFailure ReadyState = 1
	NvRateLimit    ReadyState = 2
)

//export _plat__NVEnable
func _plat__NVEnable(platParameter unsafe.Pointer, paramSize C.size_t) C.int {
	if Current.Storage.OnEnable != nil {
		return C.int(Current.Storage.OnEnable())
	}
	return 0
}

//export _plat__GetNvReadyState
func _plat__GetNvReadyState() C.int {
	if Current.Storage.OnReady != nil {
		return C.int(Current.Storage.OnReady())
	}
	return C.int(NvReady)
}

func getNvSlice(startOffset, size C.uint) []byte {
	data := Current.Storage.Data[:]
	if uint64(startOffset)+uint64(size) > uint64(len(data)) {
		return nil
	}
	return data[startOffset : startOffset+size]
}

//export _plat__NvMemoryRead
func _plat__NvMemoryRead(startOffset, size C.uint, data unsafe.Pointer) C.int {
	src := getNvSlice(startOffset, size)
	dst := unsafe.Slice((*byte)(data), size)
	if src == nil {
		return C.FALSE
	}
	copy(dst, src)
	return C.TRUE
}

//export _plat__NvGetChangedStatus
func _plat__NvGetChangedStatus(startOffset, size C.uint, data unsafe.Pointer) C.int {
	src := getNvSlice(startOffset, size)
	test := unsafe.Slice((*byte)(data), size)
	if src == nil {
		return C.NV_INVALID_LOCATION
	} else if slices.Equal(src, test) {
		return C.NV_IS_SAME
	} else {
		return C.NV_HAS_CHANGED
	}
}

//export _plat__NvMemoryWrite
func _plat__NvMemoryWrite(startOffset, size C.uint, data unsafe.Pointer) C.int {
	src := unsafe.Slice((*byte)(data), size)
	dst := getNvSlice(startOffset, size)
	if dst == nil {
		return C.FALSE
	}
	copy(dst, src)
	return C.TRUE
}

//export _plat__NvMemoryClear
func _plat__NvMemoryClear(startOffset, size C.uint) C.int {
	dst := getNvSlice(startOffset, size)
	if dst == nil {
		return C.FALSE
	}
	// Fill with 0xFF bytes as that's what the reference code does.
	for i := range dst {
		dst[i] = 0xFF
	}
	return C.TRUE
}

//export _plat__NvMemoryMove
func _plat__NvMemoryMove(srcOffset, dstOffset C.uint, size C.uint) C.int {
	src := getNvSlice(srcOffset, size)
	dst := getNvSlice(dstOffset, size)
	if src == nil || dst == nil {
		return C.FALSE
	}
	copy(dst, src)
	return C.TRUE
}

//export _plat__NvCommit
func _plat__NvCommit() C.int {
	if Current.Storage.OnCommit != nil {
		return C.int(Current.Storage.OnCommit())
	}
	return 0
}
