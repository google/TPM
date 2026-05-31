# Platform Methods FFI Effective Usage Count

Below is the table containing the effective usage counts of each `_plat` FFI method defined in `macros.rs` in the C reference code, determined by selectively disabling the bindings and analyzing compiler/linker errors during compilation:

| Method | Effective Usage Count |
|---|---|
| `_plat__LocalityGet` | 5 |
| `_plat__IsCanceled` | 6 |
| `_plat__PhysicalPresenceAsserted` | 2 |
| `_plat__WasPowerLost` | 1 |
| `_plat__StartTpmInit` | 1 |
| `_plat__EndOkTpmInit` | 1 |
| `_plat__TearDown` | 0 |
| `_plat__GetPlatformManufactureData` | 1 |
| `_plat__GetEntropy` | 2 |
| `_plat_GetEnabledSelfTest` | 1 |
| `_plat_debug_print` | 4 |
| `_plat_debug_printf` | 18 |
| `_plat__TimerRead` | 2 |
| `_plat__TimerWasReset` | 1 |
| `_plat__TimerWasStopped` | 2 |
| `_plat__ClockRateAdjust` | 6 |
| `_plat__NVEnable` | 1 |
| `_plat__GetNvReadyState` | 2 |
| `_plat__NvMemoryRead` | 1 |
| `_plat__NvGetChangedStatus` | 1 |
| `_plat__NvMemoryWrite` | 1 |
| `_plat__NvMemoryClear` | 3 |
| `_plat__NvMemoryMove` | 1 |
| `_plat__NvCommit` | 1 |
| `_plat__ACT_GetImplemented` | 4 |
| `_plat__ACT_GetRemaining` | 3 |
| `_plat__ACT_GetSignaled` | 4 |
| `_plat__ACT_SetSignaled` | 1 |
| `_plat__ACT_UpdateCounter` | 2 |
| `_plat__ACT_EnableTicks` | 2 |
| `_plat__ACT_Initialize` | 1 |
| `_plat__GetManufacturerCapabilityCode` | 2 |
| `_plat__GetVendorCapabilityCode` | 8 |
| `_plat_GetSpecCapabilityValue` | 1 |
| `_plat__GetVendorTpmType` | 2 |
| `_plat__GetTpmFirmwareVersionHigh` | 4 |
| `_plat__GetTpmFirmwareVersionLow` | 4 |
| `_plat__GetTpmFirmwareSvn` | 2 |
| `_plat__GetTpmFirmwareMaxSvn` | 1 |
| `_plat__GetTpmFirmwareSecret` | 1 |
| `_plat__GetTpmFirmwareSvnSecret` | 1 |
| `_platPcr__NumberOfPcrs` | 2 |
| `_platPcr_IsPcrBankDefaultActive` | 1 |
| `_platPcr__GetInitialValueForPcr` | 1 |
| `_platPcr__GetPcrInitializationAttributes` | 10 |
| `_plat__Fail` | 1 |
| `_plat__InFailureMode` | 31 |
| `_plat__GetFailureCode` | 3 |
| `_plat__GetFailureLocation` | 4 |
| `_plat__GetFailureFunctionName` | 1 |
| `_plat__GetFailureLine` | 1 |
| `_plat__NvVirtual_PopulateNvIndexInfo` | 7 |
| `_plat__NvVirtual_Read` | 1 |
| `_plat__NvVirtual_ReadPublic` | 1 |
| `_plat__NvVirtual_CapGetIndex` | 1 |
| `_plat__NvOperationAcceptsVirtualHandles` | 1 |
| `_plat__IsNvVirtualIndex` | 11 |
