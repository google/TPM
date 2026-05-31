//! Bindings to the TPM2 Reference Implementation in C
//!
//! ## Platform
//! 
//! 
//! 
//! ## Entrypoints
//! 
//! The Reference Implemenation is invoked via functions in the
//! `platform_interface/platform_to_tpm_interface.h` C header file. These
//! functions are exposed via a top-level `unsafe extern "C"` block.
//!
//! ### Manufacture
//!
//! Before any other functions are called, the TPM's NV must be initialized
//! with [`TPM_Manufacture`], and can be re-manufactured by calling
//! [`TPM_TearDown`] then calling [`TPM_Manufacture`] again.
//!
//! These functions return `0` on success, or a non-zero error code on failure.
//!
//! ### Running Commands
//! 
//! [`_TPM_Init`] must be called before executing any commands to properly
//! initialize the TPM's state (using the [`Platform`] implementation).
//! After this, [`ExecuteCommand`] can be used to run a command.
//!   - The request is provided 
//!
//! ### H-CRTM Measurement (optional)
//!
//! Functions to run a Hardware-based Configuration Register Triggered Measurement (H-CRTM) sequence:
//! - [`_TPM_Hash_Start`]: starts the H-CRTM sequence.
//! - [`_TPM_Hash_Data`]: feeds data into the active sequence.
//! - [`_TPM_Hash_End`]: ends the sequence.
//!
//! These functions return `1` (`TRUE`) on success, or `0` (`FALSE`) on failure.
//! [`_TPM_Init`] must be called before these functions are run.
//!
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

/// Value taken from 
pub const MAX_RESPONSE_SIZE: u32 = 0x1000 - 0x80;

pub type BOOL = c_int;

#[allow(clippy::missing_safety_doc)]
unsafe extern "C" {
    /// Initializes the TPM values in preparation for the TPM's first use.
    pub fn TPM_Manufacture(first_time: BOOL) -> c_int;
    /// Prepares the TPM for re-manufacture.
    pub fn TPM_TearDown() -> c_int;

    /// Initializes the TPM reference implementation state.
    pub fn _TPM_Init();
    /// Executes a TPM command.
    pub fn ExecuteCommand(
        request_size: u32,
        request: *const u8,
        response_size: *mut u32,
        response: *mut *mut u8,
    );

    /// Starts an H-CRTM measurement sequence.
    pub fn _TPM_Hash_Start() -> BOOL;
    /// Feeds data into the active H-CRTM measurement sequence.
    pub fn _TPM_Hash_Data(data_size: u32, data: *const u8) -> BOOL;
    /// Ends the active H-CRTM measurement sequence.
    pub fn _TPM_Hash_End() -> BOOL;
}
