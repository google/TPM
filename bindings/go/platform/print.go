package platform

// #include <platform_interface/tpm_to_platform_interface.h>
//
// typedef const char cchar_t;
import "C"

import (
	"io"
)

//export _plat_debug_print
func _plat_debug_print(str *C.cchar_t) {
	if Current.Debug == nil {
		return
	}
	s := C.GoString(str)
	s += "\n"
	io.WriteString(Current.Debug, s)
}
