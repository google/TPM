package platform

// #include <platform_interface/tpm_to_platform_interface.h>
import "C"

import "unsafe"

type Firmware struct {
	Version   uint64
	SVN       uint16
	MaxSVN    uint16
	Secret    Secret
	SVNSecret func(svn uint16) Secret
}

// SECRET_LEN is the len of the target buffer passed to [Secret.Fill].
const SECRET_LEN = C.PRIMARY_SEED_SIZE

type Secret interface {
	Fill(buf []byte) int
}

//export _plat__GetTpmFirmwareVersionHigh
func _plat__GetTpmFirmwareVersionHigh() uint32 {
	return uint32(Current.Firmware.Version >> 32)
}

//export _plat__GetTpmFirmwareVersionLow
func _plat__GetTpmFirmwareVersionLow() uint32 {
	return uint32(Current.Firmware.Version & 0xFFFFFFFF)
}

//export _plat__GetTpmFirmwareSvn
func _plat__GetTpmFirmwareSvn() uint16 {
	return Current.Firmware.SVN
}

//export _plat__GetTpmFirmwareMaxSvn
func _plat__GetTpmFirmwareMaxSvn() uint16 {
	return Current.Firmware.MaxSVN
}

func secretFill(s Secret, bufSize uint16, buf *uint8, outSize *uint16) C.int {
	if s == nil {
		return -1
	}
	dst := unsafe.Slice(buf, bufSize)
	if bytesFilled := s.Fill(dst); bytesFilled >= 0 && bytesFilled <= len(dst) {
		*outSize = uint16(bytesFilled)
		return 0
	}
	return -1
}

//export _plat__GetTpmFirmwareSecret
func _plat__GetTpmFirmwareSecret(bufSize uint16, buf *uint8, outSize *uint16) C.int {
	return secretFill(Current.Firmware.Secret, bufSize, buf, outSize)
}

//export _plat__GetTpmFirmwareSvnSecret
func _plat__GetTpmFirmwareSvnSecret(svn uint16, bufSize uint16, buf *uint8, outSize *uint16) C.int {
	var secret Secret
	if Current.Firmware.SVNSecret != nil {
		secret = Current.Firmware.SVNSecret(svn)
	}
	return secretFill(secret, bufSize, buf, outSize)
}
