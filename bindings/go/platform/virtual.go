package platform

// #include <platform_interface/tpm_to_platform_interface.h>
import "C"

/*** Virtual NV Stubs (always return TPM_RC_NO_RESULT / NO / FALSE) ***/

//export _plat__NvVirtual_PopulateNvIndexInfo
func _plat__NvVirtual_PopulateNvIndexInfo(handle C.TPM_HANDLE, publicArea *C.TPMS_NV_PUBLIC, authValue *C.TPM2B_AUTH) C.TPM_RC {
	_, _, _ = handle, publicArea, authValue
	return C.TPM_RC_NO_RESULT
}

//export _plat__NvVirtual_Read
func _plat__NvVirtual_Read(dataIn *C.NV_Read_In, dataOut *C.NV_Read_Out) C.TPM_RC {
	_, _ = dataIn, dataOut
	return C.TPM_RC_NO_RESULT
}

//export _plat__NvVirtual_ReadPublic
func _plat__NvVirtual_ReadPublic(dataIn *C.NV_ReadPublic_In, dataOut *C.NV_ReadPublic_Out) C.TPM_RC {
	_, _ = dataIn, dataOut
	return C.TPM_RC_NO_RESULT
}

//export _plat__NvVirtual_CapGetIndex
func _plat__NvVirtual_CapGetIndex(handle C.TPMI_DH_OBJECT, count uint32, handleList *C.TPML_HANDLE) C.TPMI_YES_NO {
	_, _, _ = handle, count, handleList
	return C.NO
}

//export _plat__NvOperationAcceptsVirtualHandles
func _plat__NvOperationAcceptsVirtualHandles(_ C.TPM_CC) C.BOOL {
	return C.FALSE
}

//export _plat__IsNvVirtualIndex
func _plat__IsNvVirtualIndex(_ C.TPM_HANDLE) C.BOOL {
	return C.FALSE
}
