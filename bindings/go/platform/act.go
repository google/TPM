package platform

// #include <platform_interface/tpm_to_platform_interface.h>
import "C"

type ACT interface {
	Initialize()
	EnableTicks()
	DisableTicks()
	Remaining() uint32
	UpdateCounter(uint32) bool
	Signaled() bool
	SetSignaled(bool)
}

//export _plat__ACT_GetImplemented
func _plat__ACT_GetImplemented(act uint32) C.int {
	// We only implement TPM_RH_ACT_0
	return b2i(act == 0)
}

//export _plat__ACT_GetRemaining
func _plat__ACT_GetRemaining(act uint32) uint32 {
	if act == 0 && Current.ACT != nil {
		return Current.ACT.Remaining()
	}
	return 0 // All ACT timeouts default to 0.
}

//export _plat__ACT_GetSignaled
func _plat__ACT_GetSignaled(act uint32) C.int {
	if act == 0 && Current.ACT != nil {
		return b2i(Current.ACT.Signaled())
	}
	return C.FALSE // Return FALSE to indicate no ACT has signaled.
}

//export _plat__ACT_SetSignaled
func _plat__ACT_SetSignaled(act uint32, on C.int) {
	if act == 0 && Current.ACT != nil {
		Current.ACT.SetSignaled(on != 0)
	}
}

//export _plat__ACT_UpdateCounter
func _plat__ACT_UpdateCounter(act uint32, newValue uint32) C.int {
	if act == 0 && Current.ACT != nil {
		return b2i(Current.ACT.UpdateCounter(newValue))
	}
	return C.TRUE // Pretend update is pending so TPM does not retry.
}

//export _plat__ACT_EnableTicks
func _plat__ACT_EnableTicks(enable C.int) {
	if Current.ACT != nil {
		if enable != 0 {
			Current.ACT.EnableTicks()
		} else {
			Current.ACT.DisableTicks()
		}
	}
}

//export _plat__ACT_Initialize
func _plat__ACT_Initialize() C.int {
	if Current.ACT != nil {
		Current.ACT.Initialize()
	}
	return C.TRUE // Return TRUE to indicate success.
}
