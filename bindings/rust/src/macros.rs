#[doc(hidden)]
#[macro_export]
macro_rules! bind {
    ($trait_name:ident, $c_name:ident ($($arg:ident : $typ:ty),*) -> $ret:ty) => {
        #[unsafe(no_mangle)]
        pub unsafe extern "C" fn $c_name($($arg : $typ),*) -> $ret {
            unsafe { Platform::$trait_name(P, $($arg),*) }
        }
    };
    ($trait_name:ident, $c_name:ident ($($arg:ident : $typ:ty ,)* ...) -> $ret:ty) => {
        #[unsafe(no_mangle)]
        pub unsafe extern "C" fn $c_name($($arg : $typ,)* args: ...) -> $ret {
            unsafe { Platform::$trait_name(P, $($arg,)* args) }
        }
    };
}

/// Registers a [`&'static impl Platform`][crate::Platform] with the reference implementation.
///
/// This macro sets up the required global bindings and FFI symbols expected by the C library.
#[macro_export]
macro_rules! register_platform { ($P:expr) => { const _: () = {
    use core::ffi::{c_char, c_int, c_void};

    use $crate::{BOOL, Platform, SpecCapabilityValue, bind};

    const P: &'static dyn Platform = $P;

    // Status
    bind!(locality_get, _plat__LocalityGet() -> u8);
    bind!(is_canceled, _plat__IsCanceled() -> BOOL);
    bind!(physical_presence_asserted, _plat__PhysicalPresenceAsserted() -> BOOL);

    // Init
    bind!(was_power_lost, _plat__WasPowerLost() -> BOOL);
    bind!(start_tpm_init, _plat__StartTpmInit() -> ());
    bind!(end_ok_tpm_init, _plat__EndOkTpmInit() -> ());

    // Manufacture
    bind!(tear_down, _plat__TearDown() -> ());
    bind!(get_platform_manufacture_data, _plat__GetPlatformManufactureData(buf: *mut u8, buf_size: u32) -> ());

    // Cryptography
    bind!(get_entropy, _plat__GetEntropy(entropy: *mut u8, amount: u32) -> i32);
    bind!(get_enabled_self_test, _plat_GetEnabledSelfTest(full_test: u8, p_to_test_vector: *mut u8, to_test_vector_size: usize) -> ());

    // Debug
    bind!(debug_print, _plat_debug_print(str: *const c_char) -> ());
    bind!(debug_printf, _plat_debug_printf(str: *const c_char, ...) -> ());

    // Timer
    bind!(timer_read, _plat__TimerRead() -> u64);
    bind!(timer_was_reset, _plat__TimerWasReset() -> BOOL);
    bind!(timer_was_stopped, _plat__TimerWasStopped() -> BOOL);
    bind!(clock_rate_adjust, _plat__ClockRateAdjust(adjustment: i32) -> ());

    // NV Memory
    bind!(nv_enable, _plat__NVEnable(plat_parameter: *mut c_void, param_size: usize) -> c_int);
    bind!(get_nv_ready_state, _plat__GetNvReadyState() -> c_int);
    bind!(nv_memory_read, _plat__NvMemoryRead(start_offset: u32, size: u32, data: *mut c_void) -> BOOL);
    bind!(nv_get_changed_status, _plat__NvGetChangedStatus(start_offset: u32, size: u32, data: *mut c_void) -> c_int);
    bind!(nv_memory_write, _plat__NvMemoryWrite(start_offset: u32, size: u32, data: *mut c_void) -> BOOL);
    bind!(nv_memory_clear, _plat__NvMemoryClear(start_offset: u32, size: u32) -> BOOL);
    bind!(nv_memory_move, _plat__NvMemoryMove(src_offset: u32, dst_offset: u32, size: u32) -> BOOL);
    bind!(nv_commit, _plat__NvCommit() -> c_int);

    // ACT
    bind!(act_get_implemented, _plat__ACT_GetImplemented(act: u32) -> BOOL);
    bind!(act_get_remaining, _plat__ACT_GetRemaining(act: u32) -> u32);
    bind!(act_get_signaled, _plat__ACT_GetSignaled(act: u32) -> BOOL);
    bind!(act_set_signaled, _plat__ACT_SetSignaled(act: u32, on: BOOL) -> ());
    bind!(act_update_counter, _plat__ACT_UpdateCounter(act: u32, new_value: u32) -> BOOL);
    bind!(act_enable_ticks, _plat__ACT_EnableTicks(enable: BOOL) -> ());
    bind!(act_initialize, _plat__ACT_Initialize() -> BOOL);

    // Capabilities
    bind!(get_manufacturer_capability_code, _plat__GetManufacturerCapabilityCode() -> u32);
    bind!(get_vendor_capability_code, _plat__GetVendorCapabilityCode(index: c_int) -> u32);
    bind!(get_spec_capability_value, _plat_GetSpecCapabilityValue(return_data: *mut SpecCapabilityValue) -> ());
    bind!(get_vendor_tpm_type, _plat__GetVendorTpmType() -> u32);

    // Firmware
    bind!(get_tpm_firmware_version_high, _plat__GetTpmFirmwareVersionHigh() -> u32);
    bind!(get_tpm_firmware_version_low, _plat__GetTpmFirmwareVersionLow() -> u32);
    bind!(get_tpm_firmware_svn, _plat__GetTpmFirmwareSvn() -> u16);
    bind!(get_tpm_firmware_max_svn, _plat__GetTpmFirmwareMaxSvn() -> u16);
    bind!(get_tpm_firmware_secret, _plat__GetTpmFirmwareSecret(buf_size: u16, buf: *mut u8, out_size: *mut u16) -> c_int);
    bind!(get_tpm_firmware_svn_secret, _plat__GetTpmFirmwareSvnSecret(svn: u16, buf_size: u16, buf: *mut u8, out_size: *mut u16) -> c_int);

    // PCRs
    bind!(number_of_pcrs, _platPcr__NumberOfPcrs() -> u32);
    bind!(is_pcr_bank_default_active, _platPcr_IsPcrBankDefaultActive(pcr_alg: u16) -> BOOL);
    bind!(get_initial_value_for_pcr, _platPcr__GetInitialValueForPcr(pcr_number: u32, pcr_alg: u16, startup_locality: u8, pcr_buffer: *mut u8, buffer_size: u16, pcr_length: *mut u16) -> u32);
    bind!(get_pcr_initialization_attributes, _platPcr__GetPcrInitializationAttributes(pcr_number: u32) -> u32);

    // Failure Mode
    bind!(fail, _plat__Fail(function: *const c_char, line: c_int, location_code: u64, failure_code: c_int) -> ());
    bind!(in_failure_mode, _plat__InFailureMode() -> BOOL);
    bind!(get_failure_code, _plat__GetFailureCode() -> u32);
    bind!(get_failure_location, _plat__GetFailureLocation() -> u64);
    bind!(get_failure_function_name, _plat__GetFailureFunctionName() -> *const c_char);
    bind!(get_failure_line, _plat__GetFailureLine() -> u32);

    // Virtual NV
    bind!(nv_virtual_populate_nv_index_info, _plat__NvVirtual_PopulateNvIndexInfo(handle: u32, public_area: *mut c_void, auth_value: *mut c_void) -> u32);
    bind!(nv_virtual_read, _plat__NvVirtual_Read(data_in: *const c_void, data_out: *mut c_void) -> u32);
    bind!(nv_virtual_read_public, _plat__NvVirtual_ReadPublic(data_in: *const c_void, data_out: *mut c_void) -> u32);
    bind!(nv_virtual_cap_get_index, _plat__NvVirtual_CapGetIndex(handle: u32, count: u32, handle_list: *mut c_void) -> u8);
    bind!(nv_operation_accepts_virtual_handles, _plat__NvOperationAcceptsVirtualHandles(arg_0: u32) -> BOOL);
    bind!(is_nv_virtual_index, _plat__IsNvVirtualIndex(arg_0: u32) -> BOOL);
}; }; }
