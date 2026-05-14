package platform

// #include <platform_interface/tpm_to_platform_interface.h>
import "C"

import "encoding/binary"

type Properties struct {
	VendorID     uint32
	VendorString [16]byte
	TPMType      uint32
	TPMSpec      Spec
	PlatformSpec Spec
}

type Spec struct {
	Family  uint32
	Level   uint32
	Version uint32
	Errata  uint32 // Errata was previously called DAY_OF_YEAR
	Year    uint32
}

//export _plat__GetManufacturerCapabilityCode
func _plat__GetManufacturerCapabilityCode() uint32 {
	return Current.Properties.VendorID
}

//export _plat__GetVendorCapabilityCode
func _plat__GetVendorCapabilityCode(index C.int) uint32 {
	start := (index - 1) * 4 // index is ONE-BASED
	chunk := Current.Properties.VendorString[start : start+4]
	return binary.BigEndian.Uint32(chunk)
}

//export _plat__GetVendorTpmType
func _plat__GetVendorTpmType() uint32 {
	return Current.Properties.TPMType
}

//export _plat_GetSpecCapabilityValue
func _plat_GetSpecCapabilityValue(returnData *C.SPEC_CAPABILITY_VALUE) {
	t := &Current.Properties.TPMSpec
	returnData.tpmSpecLevel = C.uint32_t(t.Level)
	returnData.tpmSpecVersion = C.uint32_t(t.Version)
	returnData.tpmSpecYear = C.uint32_t(t.Year)
	returnData.tpmSpecDayOfYear = C.uint32_t(t.Errata)

	p := &Current.Properties.PlatformSpec
	returnData.platformFamily = C.uint32_t(p.Family)
	returnData.platfromLevel = C.uint32_t(p.Level)
	returnData.platformRevision = C.uint32_t(p.Version)
	returnData.platformYear = C.uint32_t(p.Year)
	returnData.platformDayOfYear = C.uint32_t(p.Errata)
}
