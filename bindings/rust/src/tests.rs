//! Ensure we link correctly with a "always failing" platform implementation.
#![cfg(test)]

use core::{
    ffi::{c_char, c_int, c_void},
    ptr, slice,
};

use super::*;

const FALSE: BOOL = 0;
const TRUE: BOOL = 1;

#[test]
fn test_hcrtm() {
    unsafe { _TPM_Init() };

    assert_eq!(unsafe { _TPM_Hash_Start() }, TRUE);

    let data = [1, 2, 3, 4];
    let (data_ptr, data_len) = (data.as_ptr(), data.len() as u32);
    assert_eq!(unsafe { _TPM_Hash_Data(data_len, data_ptr) }, TRUE);

    assert_eq!(unsafe { _TPM_Hash_End() }, TRUE);
}

#[test]
fn fail_execute_command() {
    let cmd = [];
    let (cmd_ptr, cmd_len) = (cmd.as_ptr(), cmd.len() as u32);
    let mut rsp = [0u8; 4096];
    let (mut rsp_ptr, mut rsp_len) = (rsp.as_mut_ptr(), rsp.len() as u32);

    unsafe { ExecuteCommand(cmd_len, cmd_ptr, &mut rsp_len, &mut rsp_ptr) };

    const FAILURE: [u8; 10] = [
        0x80, 0x01, // tag = 0x8001 (TPM_ST_NO_SESSIONS)
        0x00, 0x00, 0x00, 0x0a, // responseSize = 10
        0x00, 0x00, 0x01, 0x01, // responseCode = 0x101 (TPM_RC_FAILURE)
    ];
    let rsp_out = unsafe { slice::from_raw_parts(rsp_ptr, rsp_len as usize) };
    assert_eq!(rsp_out, &FAILURE);
}

#[test]
fn fail_tpm_manufacture() {
    const MANUF_INVALID_CONFIG: c_int = -1;
    assert_eq!(unsafe { TPM_Manufacture(1) }, MANUF_INVALID_CONFIG);

    const TEARDOWN_OK: c_int = 0;
    assert_eq!(unsafe { TPM_TearDown() }, TEARDOWN_OK);
}

struct FailPlatform;
register_platform!(&FailPlatform);

const TPM_RC_FAILURE: u32 = 0x101;
const TPM_RC_NO_RESULT: u32 = 0x154;
const NV_WRITEFAILURE: c_int = 1;
const NV_INVALID_LOCATION: c_int = -1;

impl Platform for FailPlatform {
    // Status
    unsafe fn locality_get(&self) -> u8 {
        0
    }
    unsafe fn is_canceled(&self) -> BOOL {
        FALSE
    }
    unsafe fn physical_presence_asserted(&self) -> BOOL {
        FALSE
    }

    // Init
    unsafe fn was_power_lost(&self) -> BOOL {
        FALSE
    }
    unsafe fn start_tpm_init(&self) {}
    unsafe fn end_ok_tpm_init(&self) {}

    // Manufacture
    unsafe fn tear_down(&self) {}
    unsafe fn get_platform_manufacture_data(&self, _: *mut u8, _: u32) {}

    // Cryptography
    unsafe fn get_entropy(&self, _: *mut u8, _: u32) -> i32 {
        -1
    }
    unsafe fn get_enabled_self_test(&self, _: u8, _: *mut u8, _: usize) {}

    // Debug
    unsafe fn debug_print(&self, _: *const c_char) {}
    unsafe fn debug_printf(&self, _: *const c_char, _: core::ffi::VaList) {}

    // Timer
    unsafe fn timer_read(&self) -> u64 {
        0
    }
    unsafe fn timer_was_reset(&self) -> BOOL {
        FALSE
    }
    unsafe fn timer_was_stopped(&self) -> BOOL {
        FALSE
    }
    unsafe fn clock_rate_adjust(&self, _: i32) {}

    // NV Memory
    unsafe fn nv_enable(&self, _: *mut c_void, _: usize) -> c_int {
        -1
    }
    unsafe fn get_nv_ready_state(&self) -> c_int {
        NV_WRITEFAILURE
    }
    unsafe fn nv_memory_read(&self, _: u32, _: u32, _: *mut c_void) -> BOOL {
        FALSE
    }
    unsafe fn nv_get_changed_status(&self, _: u32, _: u32, _: *mut c_void) -> c_int {
        NV_INVALID_LOCATION
    }
    unsafe fn nv_memory_write(&self, _: u32, _: u32, _: *mut c_void) -> BOOL {
        FALSE
    }
    unsafe fn nv_memory_clear(&self, _: u32, _: u32) -> BOOL {
        FALSE
    }
    unsafe fn nv_memory_move(&self, _: u32, _: u32, _: u32) -> BOOL {
        FALSE
    }
    unsafe fn nv_commit(&self) -> c_int {
        -1
    }

    // ACT
    unsafe fn act_get_implemented(&self, _: u32) -> BOOL {
        FALSE
    }
    unsafe fn act_get_remaining(&self, _: u32) -> u32 {
        0
    }
    unsafe fn act_get_signaled(&self, _: u32) -> BOOL {
        FALSE
    }
    unsafe fn act_set_signaled(&self, _: u32, _: BOOL) {}
    unsafe fn act_update_counter(&self, _: u32, _: u32) -> BOOL {
        FALSE
    }
    unsafe fn act_enable_ticks(&self, _: BOOL) {}
    unsafe fn act_initialize(&self) -> BOOL {
        FALSE
    }

    // Capabilities
    unsafe fn get_manufacturer_capability_code(&self) -> u32 {
        0
    }
    unsafe fn get_vendor_capability_code(&self, _: c_int) -> u32 {
        0
    }
    unsafe fn get_spec_capability_value(&self, _: *mut SpecCapabilityValue) {}
    unsafe fn get_vendor_tpm_type(&self) -> u32 {
        0
    }

    // Firmware
    unsafe fn get_tpm_firmware_version_high(&self) -> u32 {
        0
    }
    unsafe fn get_tpm_firmware_version_low(&self) -> u32 {
        0
    }
    unsafe fn get_tpm_firmware_svn(&self) -> u16 {
        0
    }
    unsafe fn get_tpm_firmware_max_svn(&self) -> u16 {
        0
    }
    unsafe fn get_tpm_firmware_secret(&self, _: u16, _: *mut u8, _: *mut u16) -> c_int {
        -1
    }
    unsafe fn get_tpm_firmware_svn_secret(&self, _: u16, _: u16, _: *mut u8, _: *mut u16) -> c_int {
        -1
    }

    // PCRs
    unsafe fn number_of_pcrs(&self) -> u32 {
        0
    }
    unsafe fn is_pcr_bank_default_active(&self, _: u16) -> BOOL {
        FALSE
    }
    unsafe fn get_initial_value_for_pcr(
        &self,
        _: u32,
        _: u16,
        _: u8,
        _: *mut u8,
        _: u16,
        _: *mut u16,
    ) -> u32 {
        TPM_RC_FAILURE
    }
    unsafe fn get_pcr_initialization_attributes(&self, _: u32) -> u32 {
        0
    }

    // Failure Mode
    unsafe fn fail(&self, _: *const c_char, _: c_int, _: u64, _: c_int) {}
    unsafe fn in_failure_mode(&self) -> BOOL {
        TRUE
    }
    unsafe fn get_failure_code(&self) -> u32 {
        0
    }
    unsafe fn get_failure_location(&self) -> u64 {
        0
    }
    unsafe fn get_failure_function_name(&self) -> *const c_char {
        ptr::null()
    }
    unsafe fn get_failure_line(&self) -> u32 {
        0
    }

    // Virtual NV
    unsafe fn nv_virtual_populate_nv_index_info(
        &self,
        _: u32,
        _: *mut c_void,
        _: *mut c_void,
    ) -> u32 {
        TPM_RC_NO_RESULT
    }
    unsafe fn nv_virtual_read(&self, _: *const c_void, _: *mut c_void) -> u32 {
        TPM_RC_NO_RESULT
    }
    unsafe fn nv_virtual_read_public(&self, _: *const c_void, _: *mut c_void) -> u32 {
        TPM_RC_NO_RESULT
    }
    unsafe fn nv_virtual_cap_get_index(&self, _: u32, _: u32, _: *mut c_void) -> u8 {
        0
    }
    unsafe fn nv_operation_accepts_virtual_handles(&self, _: u32) -> BOOL {
        FALSE
    }
    unsafe fn is_nv_virtual_index(&self, _: u32) -> BOOL {
        FALSE
    }
}
