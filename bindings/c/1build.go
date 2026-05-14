// Package c handles compiling the C code in the TPMCmd/tpm/ directory.
//
// "Listen to my tale of woe." - Frankenstein's Monster
//
// # Overview
//
// We aim to compile the entire TPM reference code using only the Go build
// system. Significant contortions are required to achieve this without
// copying C files.
//
// Go's CGO build system is very primitive (to put it politely). While it can
// include headers (.h files) from any location, it can only compile sources
// (.c files) located in the same directory as the Go package. CGO does not
// search subdirectories for sources, so any C files we want to compile must be
// generated in the top level of this directory.
//
// # Notes on CGO Build Behavior
//
//   - **Header Caching**: Go tracks changes to `.c` files in the package
//     directory, but may not detect changes to header files in include paths.
//     If you modify a header file, you must run `go build -a` or
//     `go clean -cache` to force a rebuild.
//   - **Linking Errors**: Linking errors usually appear as missing symbols
//     during the final Go link phase because CGO compiles files independently.
//   - **Threading**: Go 1.26 (via CL 694475) runs CGO compiles in parallel,
//     reducing build times significantly (e.g., from 42s in Go 1.25 to 7s).
//   - **Underscore Files**: Go ignores files starting with `_`. The generation
//     script (regenerate_c.py) strips the leading underscore to
//     allow CGO to compile them.
//   - **Name**: This file is named `1build.go` so that it appears at the top
//     of the file list in editors and directory listings.
//   - **Directory Structure**: `bindings/c/` contains the CGO build file and
//     wrapper files to C sources. The C files access headers in `TPMCmd` at the
//     repository root via relative paths (`../../TPMCmd`).
package c

// #cgo CFLAGS: -std=c11 -Wall -Wextra -pedantic -Werror
// // We do crimes with function pointer casts in the marshalling and hash code.
// #cgo CFLAGS: -Wno-cast-function-type
//
// #cgo CFLAGS: -I ../../TPMCmd/TpmConfiguration/
// #cgo CFLAGS: -I ../../TPMCmd/tpm/include/
// #cgo CFLAGS: -I ../../TPMCmd/tpm/include/private/
// #cgo CFLAGS: -I ../../TPMCmd/tpm/include/private/prototypes/
// #cgo CFLAGS: -I ../../TPMCmd/tpm/cryptolibs/common/include/
// #cgo CFLAGS: -I ../../TPMCmd/tpm/cryptolibs/Ossl/include/
// #cgo CFLAGS: -I ../../TPMCmd/tpm/cryptolibs/TpmBigNum/
// #cgo CFLAGS: -I ../../TPMCmd/tpm/cryptolibs/TpmBigNum/include/
//
// // Link against the system OpenSSL
// #cgo LDFLAGS: -lcrypto
// #cgo CFLAGS: -DSYM_LIB=Ossl
// #cgo CFLAGS: -DHASH_LIB=Ossl
// #cgo CFLAGS: -DMATH_LIB=TpmBigNum
// #cgo CFLAGS: -DBN_MATH_LIB=Ossl
// // Flags to find OpenSSL installation on macOS (default Homebrew location)
// #cgo darwin,amd64 CFLAGS: -I/usr/local/opt/openssl/include
// #cgo darwin,amd64 LDFLAGS: -L/usr/local/opt/openssl/lib
// #cgo darwin,arm64 CFLAGS: -I/opt/homebrew/opt/openssl/include
// #cgo darwin,arm64 LDFLAGS: -L/opt/homebrew/opt/openssl/lib
import "C"
