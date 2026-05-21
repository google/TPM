use std::{ffi::OsStr, fs::read_dir};

fn main() {
    let mut build = cc::Build::new();

    // Add all .c files from ../c
    for entry in read_dir("../c").unwrap() {
        let path = entry.unwrap().path();
        if path.extension().and_then(OsStr::to_str) == Some("c") {
            build.file(path);
        }
    }

    // Add includes
    build.include("../../TPMCmd/TpmConfiguration/");
    build.include("../../TPMCmd/tpm/include/");
    build.include("../../TPMCmd/tpm/include/private/");
    build.include("../../TPMCmd/tpm/include/private/prototypes/");
    build.include("../../TPMCmd/tpm/cryptolibs/common/include/");
    build.include("../../TPMCmd/tpm/cryptolibs/Ossl/include/");
    build.include("../../TPMCmd/tpm/cryptolibs/TpmBigNum/");
    build.include("../../TPMCmd/tpm/cryptolibs/TpmBigNum/include/");

    // Add defines
    build.define("SYM_LIB", "Ossl");
    build.define("HASH_LIB", "Ossl");
    build.define("MATH_LIB", "TpmBigNum");
    build.define("BN_MATH_LIB", "Ossl");

    // Add flags
    build.flag("-std=c11");
    build.warnings(true);
    build.extra_warnings(true);
    build.flag("-pedantic");
    build.flag("-Wno-cast-function-type");

    // Compile
    build.compile("tpm2_ref");

    // Link against OpenSSL
    println!("cargo:rustc-link-lib=crypto");
}
