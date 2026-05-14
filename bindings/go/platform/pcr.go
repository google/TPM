package platform

// #include <platform_interface/tpm_to_platform_interface.h>
//
// static inline void set_stateSave(PCR_Attributes* attr, int b) { attr->stateSave = b; }
// static inline void set_doNotIncrement(PCR_Attributes* attr, int b) { attr->doNotIncrementPcrCounter = b; }
// static inline void set_resetLocality(PCR_Attributes* attr, uint8_t loc) { attr->resetLocality = loc; }
// static inline void set_extendLocality(PCR_Attributes* attr, uint8_t loc) { attr->extendLocality = loc; }
import "C"

import "unsafe"

const NUM_PCRS = C.IMPLEMENTATION_PCR

type PCRAttributes struct {
	ShouldSave     bool
	DoNotIncrement bool
	ResetLocality  uint8
	ExtendLocality uint8
}

type PCRConfig interface {
	Attributes(pcrNumber uint32) PCRAttributes
	DefaultActive(alg uint16) bool
	// InitialValue of the specified PCR, written to value.
	//
	// The provided []byte will always have the the correct size, so this
	// method should never fail.
	InitialValue(pcrNumber uint32, alg uint16, locality uint8, value []byte)
}

func (p *Platform) getPCRConfig() PCRConfig {
	if p.PCRConfig == nil {
		return PCClientDefaultConfig
	}
	return p.PCRConfig
}

var PCClientDefaultConfig PCRConfig = pcClientDefaultConfig{}

type pcClientDefaultConfig struct{}

func (c pcClientDefaultConfig) Attributes(pcrNumber uint32) PCRAttributes {
	switch pcrNumber {
	case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15:
		return PCRAttributes{ShouldSave: true, ExtendLocality: 0x1F}
	case 16:
		// Note: The PC Client spec requires DoNotIncrement: true, but the reference C code uses false.
		return PCRAttributes{ShouldSave: false, DoNotIncrement: false, ResetLocality: 0x0F, ExtendLocality: 0x1F}
	case 17, 18:
		return PCRAttributes{ShouldSave: false, DoNotIncrement: false, ResetLocality: 0x10, ExtendLocality: 0x1C}
	case 19:
		return PCRAttributes{ShouldSave: false, DoNotIncrement: false, ResetLocality: 0x10, ExtendLocality: 0x0C}
	case 20:
		// Note: The PC Client spec requires DoNotIncrement: false, but the reference C code uses true.
		return PCRAttributes{ShouldSave: false, DoNotIncrement: true, ResetLocality: 0x14, ExtendLocality: 0x0E}
	case 21, 22:
		return PCRAttributes{ShouldSave: false, DoNotIncrement: true, ResetLocality: 0x14, ExtendLocality: 0x04}
	case 23:
		// Note: The PC Client spec requires DoNotIncrement: true, but the reference C code uses false.
		return PCRAttributes{ShouldSave: false, DoNotIncrement: false, ResetLocality: 0x0F, ExtendLocality: 0x1F}
	default:
		return PCRAttributes{}
	}
}

func (c pcClientDefaultConfig) DefaultActive(alg uint16) bool {
	// Default to active for SHA-1 (0x0004) and SHA-256 (0x000B)
	return alg == 0x0004 || alg == 0x000B
}

func (c pcClientDefaultConfig) InitialValue(pcrNumber uint32, alg uint16, locality uint8, buf []byte) {
	defaultValue := byte(0)
	attr := c.Attributes(pcrNumber)
	if (attr.ResetLocality & 0x10) != 0 {
		defaultValue = 0xFF
	}
	for i := range buf {
		buf[i] = defaultValue
	}
	// HCRTM_PCR is 0
	if pcrNumber == 0 {
		buf[len(buf)-1] = locality
	}
}

//export _platPcr__NumberOfPcrs
func _platPcr__NumberOfPcrs() uint32 {
	return NUM_PCRS
}

//export _platPcr_IsPcrBankDefaultActive
func _platPcr_IsPcrBankDefaultActive(pcrAlg C.TPM_ALG_ID) C.BOOL {
	return b2i(Current.getPCRConfig().DefaultActive(uint16(pcrAlg)))
}

//export _platPcr__GetInitialValueForPcr
func _platPcr__GetInitialValueForPcr(
	pcrNumber uint32,
	pcrAlg uint16,
	startupLocality uint8,
	pcrBuffer *uint8,
	bufferSize uint16,
	pcrLength *uint16,
) C.TPM_RC {
	dst := unsafe.Slice((*byte)(pcrBuffer), bufferSize)
	Current.getPCRConfig().InitialValue(
		uint32(pcrNumber),
		uint16(pcrAlg),
		uint8(startupLocality),
		dst,
	)
	if pcrLength != nil {
		*pcrLength = bufferSize
	}
	return 0
}

//export _platPcr__GetPcrInitializationAttributes
func _platPcr__GetPcrInitializationAttributes(pcrNumber uint32) C.PCR_Attributes {
	a := Current.getPCRConfig().Attributes(pcrNumber)

	var attr C.PCR_Attributes
	C.set_stateSave(&attr, b2i(a.ShouldSave))
	C.set_doNotIncrement(&attr, b2i(a.DoNotIncrement))
	C.set_resetLocality(&attr, C.uint8_t(a.ResetLocality))
	C.set_extendLocality(&attr, C.uint8_t(a.ExtendLocality))
	return attr
}
