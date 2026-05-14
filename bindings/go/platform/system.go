package platform

// #include <platform_interface/tpm_to_platform_interface.h>
import "C"

type System interface {
	// Locality returns the locality (0-4 or 32-255) of the current command.
	Locality() uint8
	// Canceled returns true if the current command should be canceled.
	Canceled() bool
	// PhysicalPresence returns true if Physical Presense is asserted.
	PhysicalPresence() bool
	// StartInit is called when a _TPM_Init indication starts.
	StartInit()
	// EndOkInit is called when a _TPM_Init indication completes successfully.
	EndOkInit()
	// TearDown is called when the TPM is being prepared for re-manufacture.
	TearDown()
}

// b2i converts a Go bool to a C int
func b2i(b bool) C.int {
	if b {
		return 1
	}
	return 0
}

//export _plat__LocalityGet
func _plat__LocalityGet() uint8 {
	if Current.System != nil {
		return Current.System.Locality()
	}
	return 0 // Default locality is 0.
}

//export _plat__IsCanceled
func _plat__IsCanceled() C.int {
	// By default, no commands are cancelled
	return b2i(Current.System != nil && Current.System.Canceled())
}

//export _plat__WasPowerLost
func _plat__WasPowerLost() C.int {
	// We just always assume power was lost if we are calling Init again.
	return C.TRUE
}

//export _plat__PhysicalPresenceAsserted
func _plat__PhysicalPresenceAsserted() C.int {
	// By default, physical presence is not asserted
	return b2i(Current.System != nil && Current.System.PhysicalPresence())
}

//export _plat__StartTpmInit
func _plat__StartTpmInit() {
	if Current.System != nil {
		Current.System.StartInit()
	}
}

//export _plat__EndOkTpmInit
func _plat__EndOkTpmInit() {
	if Current.System != nil {
		Current.System.EndOkInit()
	}
}

//export _plat__TearDown
func _plat__TearDown() {
	if Current.System != nil {
		Current.System.TearDown()
	}
}
