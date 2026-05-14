package platform

// #include <platform_interface/tpm_to_platform_interface.h>
//
// typedef const char cchar_t;
import "C"

import (
	"fmt"
)

type Failure interface {
	// Fail puts the TPM into failure mode.
	Fail(info FailureInfo)
	// FailureInfo returns information about why a TPM entered failure mode.
	// Returns nil if the TPM is not in failure mode.
	FailureInfo() *FailureInfo
}

// FailureInfo holds information about a TPM failure.
//
// *FailureInfo implements the Failure interface, acting as a simple
// container that stores the last failure.
type FailureInfo struct {
	Code      uint32
	Location  uint64
	Line      uint32
	cFuncName *C.cchar_t
}

func (fi *FailureInfo) FunctionName() string {
	return C.GoString(fi.cFuncName)
}

func (fi *FailureInfo) Fail(info FailureInfo) {
	*fi = info
}

func (fi *FailureInfo) FailureInfo() *FailureInfo {
	return fi
}

func (p *Platform) getFailureInfo() *FailureInfo {
	if p.Failure == nil {
		return nil
	}
	return p.Failure.FailureInfo()
}

//export _plat__Fail
func _plat__Fail(function *C.cchar_t, line C.int, locationCode C.uint64_t, failureCode C.int) {
	info := FailureInfo{
		Code:      uint32(failureCode),
		Location:  uint64(locationCode),
		Line:      uint32(line),
		cFuncName: function,
	}
	if Current.Failure == nil {
		panic(fmt.Sprintf("%+v", info))
	}
	Current.Failure.Fail(info)
}

//export _plat__InFailureMode
func _plat__InFailureMode() C.int {
	return b2i(Current.getFailureInfo() != nil)
}

//export _plat__GetFailureCode
func _plat__GetFailureCode() uint32 {
	if info := Current.getFailureInfo(); info != nil {
		return info.Code
	}
	return 0
}

//export _plat__GetFailureLocation
func _plat__GetFailureLocation() C.uint64_t {
	if info := Current.getFailureInfo(); info != nil {
		return C.uint64_t(info.Location)
	}
	return 0
}

//export _plat__GetFailureFunctionName
func _plat__GetFailureFunctionName() *C.cchar_t {
	if info := Current.getFailureInfo(); info != nil {
		return info.cFuncName
	}
	return nil
}

//export _plat__GetFailureLine
func _plat__GetFailureLine() uint32 {
	if info := Current.getFailureInfo(); info != nil {
		return info.Line
	}
	return 0
}
