//! Bindings to the TPM2 Reference Implementation in C
//!
//! This includes the entrypoints to the reference code:
//!   - [`TPM_Manufacture`]: manufacture the NV Data
//!   - [`_TPM_Init`]: initialize the reference code
//!   - [`_TPM_Hash_Start`]/[`_TPM_Hash_Data`]/[`_TPM_Hash_End`]: H-CRTM
//!   - [`ExecuteCommand`]: run a command
//!
//! It also includes the [`Platform`] trait, so that a user can provide the
//! platform-specific functionality required by the Reference Implementation.
//! A `&'static impl Platform` can be registered via [`register_platform!`],
//! or an external platform implementation can be linked in.
#![no_std]
#![feature(c_variadic)] // About to be stabilized, see: https://github.com/rust-lang/rust/issues/44930

mod macros;
mod platform;
mod tests;

use core::ffi::c_int;

pub use platform::*;

pub const MAX_RESPONSE_SIZE: u32 = 0x1000 - 0x80;

pub type BOOL = c_int;

unsafe extern "C" {
    /// Manufactures the TPM's NV Data in preparation for the first use.
    pub fn TPM_Manufacture(firstTime: BOOL) -> c_int;

    /// Initializes the TPM values and subsystem.
    ///
    /// This function must be called before:
    ///   - [`_TPM_Hash_Start`]
    ///   - [`ExecuteCommand`]
    ///
    /// # Safety
    ///   - NV Data must be initialized via [`TPM_Manufacture`] beforehand.
    pub fn _TPM_Init();

    /// Starts an H-CRTM Event Sequence.
    ///
    /// # Safety
    ///   - [`_TPM_Init`] must be called beforehand.
    pub fn _TPM_Hash_Start() -> BOOL;
    /// Sends data to be hashed as part of an H-CRTM Event Sequence.
    ///
    /// # Safety
    ///   - [`_TPM_Hash_Start`] must be called beforehand.
    pub fn _TPM_Hash_Data(dataSize: u32, data: *const u8) -> BOOL;
    /// Completes an H-CRTM Event Sequence.
    ///
    /// # Safety
    ///   - [`_TPM_Hash_Start`]/[`_TPM_Hash_Data`] must be called beforehand.
    pub fn _TPM_Hash_End() -> BOOL;

    /// Executes a raw TPM command.
    ///
    /// When calling this function:
    ///   - `request` / `requestSize` contain the raw request buffer.
    ///   - `response` / `responseSize` contain the raw response buffer.
    ///
    /// On return, `response` / `responseSize` contain the response. Note that
    /// the `response` pointer may contain a new value.
    ///
    /// # Safety
    ///   - `request` must point to `requestSize` bytes.
    ///   - `*response` must point to `responseSize` bytes.
    ///   - `responseSize` must be at least [`MAX_RESPONSE_SIZE`].
    ///   - [`_TPM_Init`] must be called beforehand.
    pub fn ExecuteCommand(
        requestSize: u32,
        request: *const u8,
        responseSize: *mut u32,
        response: *mut *mut u8,
    );
}
