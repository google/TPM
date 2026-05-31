use std::{fs, io};

fn main() -> io::Result<()> {
    let mut build = cc::Build::new();
    // Build for C11
    build.std("c11").flag_if_supported("-pedantic");
    // Enable most warnings
    build.warnings(true).extra_warnings(true);
    build.flag("-Wno-cast-function-type");

    // Setup include directories
    build.includes([
        "../../TPMCmd/TpmConfiguration/",
        "../../TPMCmd/tpm/include/",
        "../../TPMCmd/tpm/include/private/",
        "../../TPMCmd/tpm/include/private/prototypes/",
        "../../TPMCmd/tpm/cryptolibs/common/include/",
        "../../TPMCmd/tpm/cryptolibs/Ossl/include/",
        "../../TPMCmd/tpm/cryptolibs/TpmBigNum/",
        "../../TPMCmd/tpm/cryptolibs/TpmBigNum/include/",
    ]);

    // Add all .c source files from ../c
    for entry in fs::read_dir("../c")? {
        let path = entry?.path();
        if path.extension().is_some_and(|ext| ext == "c") {
            build.file(path);
        }
    }

    // Use OpenSSL for Cryptography
    build.define("FOO", "barrrrrrr");
    build.define("SYM_LIB", "Ossl");
    build.define("HASH_LIB", "Ossl");
    build.define("MATH_LIB", "TpmBigNum");
    build.define("BN_MATH_LIB", "Ossl");
    // Link against OpenSSL
    println!("cargo:rustc-link-lib=crypto");

    // Compile
    build.compile("tpm2_ref");
    Ok(())
}
