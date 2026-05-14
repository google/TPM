package platform

// #include <platform_interface/tpm_to_platform_interface.h>
import "C"

import (
	"time"
)

type Timer interface {
	// Ticks returns the monotonic duration since power-on.
	Ticks() time.Duration

	// WasReset returns true if the timer was reset since the last call.
	// Calling this must clear the internal "reset" state.
	WasReset() bool

	// WasStopped returns true if the timer was stopped since the last call.
	// Calling this must clear the internal "stopped" state.
	WasStopped() bool

	// Adjust changes the rate of the timer.
	Adjust(a Adjustment)
}

// Adjustment is the amount to change the tick timer's rate. Its values
// correspond to the TPM_CLOCK_ADJUST constants.
type Adjustment int

const (
	AdjustCoarseSlower Adjustment = -3
	AdjustMediumSlower Adjustment = -2
	AdjustFineSlower   Adjustment = -1
	AdjustFineFaster   Adjustment = 1
	AdjustMediumFaster Adjustment = 2
	AdjustCoarseFaster Adjustment = 3
)

//export _plat__TimerRead
func _plat__TimerRead() C.uint64_t {
	if Current.Timer == nil {
		panic("_plat__TimerRead called with Current.Timer == nil")
	}
	return C.uint64_t(Current.Timer.Ticks().Milliseconds())
}

//export _plat__TimerWasReset
func _plat__TimerWasReset() C.int {
	if Current.Timer == nil {
		panic("_plat__TimerWasReset called with Current.Timer == nil")
	}
	return b2i(Current.Timer.WasReset())
}

//export _plat__TimerWasStopped
func _plat__TimerWasStopped() C.int {
	if Current.Timer == nil {
		panic("_plat__TimerWasStopped called with Current.Timer == nil")
	}
	return b2i(Current.Timer.WasStopped())
}

//export _plat__ClockRateAdjust
func _plat__ClockRateAdjust(adjustment C._plat__ClockAdjustStep) {
	if Current.Timer == nil {
		panic("_plat__ClockRateAdjust called with Current.Timer == nil")
	}
	Current.Timer.Adjust(Adjustment(adjustment))
}
