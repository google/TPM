#include <stdarg.h>
#include <stdio.h>

#include <platform_interface/tpm_to_platform_interface.h>

// Implement _plat_debug_printf
void _plat_debug_printf(const char* format, ...) {
    va_list args;
    va_start(args, format);
    char buf[1024];
    vsnprintf(buf, sizeof(buf), format, args);
    va_end(args);
    _plat_debug_print(buf);
}
