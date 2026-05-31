use core::ffi::{VaList, c_char, c_int, c_void};

use crate::BOOL;

/// Values returned by [`Platform::get_spec_capability_value`].
#[repr(C)]
#[derive(Debug, Copy, Clone, PartialEq, Eq)]
pub struct SpecCapabilityValue {
    pub tpm_spec_level: u32,
    pub tpm_spec_version: u32,
    pub tpm_spec_year: u32,
    pub tpm_spec_day_of_year: u32,

    pub platform_family: u32,
    pub platform_level: u32,
    pub platform_revision: u32,
    pub platform_year: u32,
    pub platform_day_of_year: u32,
}

/// Platform-specific functionality called by the Core TPM library.
#[allow(clippy::missing_safety_doc)]
pub trait Platform {
    // Status
    unsafe fn locality_get(&self) -> u8;
    unsafe fn is_canceled(&self) -> BOOL;
    unsafe fn physical_presence_asserted(&self) -> BOOL;

    // Init
    unsafe fn was_power_lost(&self) -> BOOL;
    unsafe fn start_tpm_init(&self);
    unsafe fn end_ok_tpm_init(&self);

    // Manufacture
    unsafe fn get_platform_manufacture_data(&self, buf: *mut u8, buf_size: u32);
    unsafe fn tear_down(&self);

    // Cryptography
    unsafe fn get_entropy(&self, entropy: *mut u8, amount: u32) -> i32;
    unsafe fn get_enabled_self_test(
        &self,
        full_test: u8,
        test_vector: *mut u8,
        test_vector_size: usize,
    );

    // Debug
    unsafe fn debug_print(&self, str: *const c_char);
    unsafe fn debug_printf(&self, str: *const c_char, args: VaList);

    // Timer
    unsafe fn timer_read(&self) -> u64;
    unsafe fn timer_was_reset(&self) -> BOOL;
    unsafe fn timer_was_stopped(&self) -> BOOL;
    unsafe fn clock_rate_adjust(&self, adjustment: i32);

    // NV Memory
    unsafe fn nv_enable(&self, plat_parameter: *mut c_void, param_size: usize) -> c_int;
    unsafe fn get_nv_ready_state(&self) -> c_int;
    unsafe fn nv_memory_read(&self, start_offset: u32, size: u32, data: *mut c_void) -> BOOL;
    unsafe fn nv_get_changed_status(
        &self,
        start_offset: u32,
        size: u32,
        data: *mut c_void,
    ) -> c_int;
    unsafe fn nv_memory_write(&self, start_offset: u32, size: u32, data: *mut c_void) -> BOOL;
    unsafe fn nv_memory_clear(&self, start_offset: u32, size: u32) -> BOOL;
    unsafe fn nv_memory_move(&self, src_offset: u32, dst_offset: u32, size: u32) -> BOOL;
    unsafe fn nv_commit(&self) -> c_int;

    // ACT
    unsafe fn act_get_implemented(&self, act: u32) -> BOOL;
    unsafe fn act_get_remaining(&self, act: u32) -> u32;
    unsafe fn act_get_signaled(&self, act: u32) -> BOOL;
    unsafe fn act_set_signaled(&self, act: u32, on: BOOL);
    unsafe fn act_update_counter(&self, act: u32, new_value: u32) -> BOOL;
    unsafe fn act_enable_ticks(&self, enable: BOOL);
    unsafe fn act_initialize(&self) -> BOOL;

    // Capabilities
    unsafe fn get_manufacturer_capability_code(&self) -> u32;
    unsafe fn get_vendor_capability_code(&self, index: c_int) -> u32;
    unsafe fn get_spec_capability_value(&self, return_data: *mut SpecCapabilityValue);
    unsafe fn get_vendor_tpm_type(&self) -> u32;

    // Firmware
    unsafe fn get_tpm_firmware_version_high(&self) -> u32;
    unsafe fn get_tpm_firmware_version_low(&self) -> u32;
    unsafe fn get_tpm_firmware_svn(&self) -> u16;
    unsafe fn get_tpm_firmware_max_svn(&self) -> u16;
    unsafe fn get_tpm_firmware_secret(
        &self,
        buf_size: u16,
        buf: *mut u8,
        out_size: *mut u16,
    ) -> c_int;
    unsafe fn get_tpm_firmware_svn_secret(
        &self,
        svn: u16,
        buf_size: u16,
        buf: *mut u8,
        out_size: *mut u16,
    ) -> c_int;

    // PCRs
    unsafe fn number_of_pcrs(&self) -> u32;
    unsafe fn is_pcr_bank_default_active(&self, pcr_alg: u16) -> BOOL;
    unsafe fn get_initial_value_for_pcr(
        &self,
        pcr_number: u32,
        pcr_alg: u16,
        startup_locality: u8,
        pcr_buffer: *mut u8,
        buffer_size: u16,
        pcr_length: *mut u16,
    ) -> u32;
    unsafe fn get_pcr_initialization_attributes(&self, pcr_number: u32) -> u32;

    // Failure Mode
    unsafe fn fail(
        &self,
        function: *const c_char,
        line: c_int,
        location_code: u64,
        failure_code: c_int,
    );
    unsafe fn in_failure_mode(&self) -> BOOL;
    unsafe fn get_failure_code(&self) -> u32;
    unsafe fn get_failure_location(&self) -> u64;
    unsafe fn get_failure_function_name(&self) -> *const c_char;
    unsafe fn get_failure_line(&self) -> u32;

    // Virtual NV
    unsafe fn nv_virtual_populate_nv_index_info(
        &self,
        handle: u32,
        public_area: *mut c_void,
        auth_value: *mut c_void,
    ) -> u32;
    unsafe fn nv_virtual_read(&self, data_in: *const c_void, data_out: *mut c_void) -> u32;
    unsafe fn nv_virtual_read_public(&self, data_in: *const c_void, data_out: *mut c_void) -> u32;
    unsafe fn nv_virtual_cap_get_index(
        &self,
        handle: u32,
        count: u32,
        handle_list: *mut c_void,
    ) -> u8;
    unsafe fn nv_operation_accepts_virtual_handles(&self, arg_0: u32) -> BOOL;
    unsafe fn is_nv_virtual_index(&self, arg_0: u32) -> BOOL;
}
