//** Description
//
//    This file contains the NV read and write access methods.  This implementation
//    uses RAM/file and does not manage the RAM/file as NV blocks.
//    The implementation may become more sophisticated over time.
//

//** Includes and Local
#include <stdio.h>
#include <stdarg.h>
#include <time.h>
#include "Platform.h"

#if ENABLE_TPM_DEBUG_PRINT

LIB_EXPORT void _plat_debug_print(const char* str)
{
    printf("%s\n", str);
}

LIB_EXPORT void _plat_debug_print_buffer(const void* buf, const size_t size)
{
    NOT_REFERENCED(buf);
    NOT_REFERENCED(size);
    // not implemented
}

LIB_EXPORT void _plat_debug_print_int32(const char* name, uint32_t value)
{
    printf("%s=0x%04x\n", name, value);
}

LIB_EXPORT void _plat_debug_print_int64(const char* name, uint64_t value)
{
    printf("%s=0x%04x:%04x\n",
           name,
           (uint32_t)(value >> 32),
           (uint32_t)(value & 0xFFFFFFFF));
}

LIB_EXPORT void _plat_debug_printf(const char* fmt, ...)
{
    va_list params;
    va_start(params, fmt);
    vprintf(fmt, params);
}

#endif  // ENABLE_TPM_DEBUG_PRINT
