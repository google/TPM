#ifndef _OIDS_H_
#define _OIDS_H_

// All the OIDs in this file are defined as DER-encoded values with a leading tag
// 0x06 (ASN1_OBJECT_IDENTIFIER), followed by a single length byte. This allows the
// OID size to be determined by looking at octet[1] of the OID (total size is
// OID[1] + 2).

#define ANSI_X962 0x2A, 0x86, 0x48, 0xCE, 0x3D  // 1.2.840.10045

// Encoded to take two additional bytes
#define SM_SCHEME 0x06, 0x08, 0x2A, 0x81, 0x1C, 0xCF, 0x55, 1  // 1.2.156.10197.1
#define NIST_ALG  0x06, 0x09, 0x60, 0x86, 0x48, 1, 101, 3, 4   // 2.16.840.1.101.3.4
// Encoded to take one additional byte
#define NIST_HASH       NIST_ALG, 2                      // 2.16.840.1.101.3.4.2
#define NIST_SIG        NIST_ALG, 3                      // 2.16.840.1.101.3.4.3
#define ECDSA_SHA2      0x06, 0x08, ANSI_X962, 4, 3      // 1.2.840.10045.4.3
#define PRIME_CURVES    0x06, 0x08, ANSI_X962, 3, 1      // 1.2.840.10045.3.1
#define CERTICOM_CURVES 0x06, 0x05, 0x2B, 0x81, 0x04, 0  // 1.3.132.0
#define PKCS1_ALG \
    0x06, 0x09, 0x2A, 0x86, 0x48, 0x86, 0xF7, 0x0D, 1, 1  // 1.2.840.113549.1.1

// These hash OIDs used in a lot of places.
#define OID_SHA1_VALUE 0x06, 0x05, 0x2B, 14, 3, 2, 26  // 1.3.14.3.2.26
#if ALG_SHA1
MAKE_OID(_SHA1);
// Expands to:
//     EXTERN  const BYTE OID_SHA1[] INITIALIZER({OID_SHA1_VALUE})
// which, depending on the setting of EXTERN and INITIALIZER, expands to either:
//     extern const BYTE    OID_SHA1[]
// or
//     const BYTE           OID_SHA1[] = {OID_SHA1_VALUE}
// which is:
//     const BYTE           OID_SHA1[] = {0x06, 0x05, 0x2B, 0x0E, ...}
#endif  // ALG_SHA1

#define OID_SHA256_VALUE NIST_HASH, 1  // 2.16.840.1.101.3.4.2.1
#if ALG_SHA256
MAKE_OID(_SHA256);
#endif  // ALG_SHA256

#define OID_SHA384_VALUE NIST_HASH, 2  // 2.16.840.1.101.3.4.2.2
#if ALG_SHA384
MAKE_OID(_SHA384);
#endif  // ALG_SHA384

#define OID_SHA512_VALUE NIST_HASH, 3  // 2.16.840.1.101.3.4.2.3
#if ALG_SHA512
MAKE_OID(_SHA512);
#endif  // ALG_SHA512

#define OID_SM3_256_VALUE SM_SCHEME, 0x83, 0x11  // 1.2.156.10197.1.401
#if ALG_SM3_256
MAKE_OID(_SM3_256);
#endif  // ALG_SM3_256

#define OID_SHA3_256_VALUE NIST_HASH, 8  // 2.16.840.1.101.3.4.2.8
#if ALG_SHA3_256
MAKE_OID(_SHA3_256);
#endif  // ALG_SHA3_256

#define OID_SHA3_384_VALUE NIST_HASH, 9  // 2.16.840.1.101.3.4.2.9
#if ALG_SHA3_384
MAKE_OID(_SHA3_384);
#endif  // ALG_SHA3_384

#define OID_SHA3_512_VALUE NIST_HASH, 10  // 2.16.840.1.101.3.4.2.10
#if ALG_SHA3_512
MAKE_OID(_SHA3_512);
#endif  // ALG_SHA3_512

// These are used for RSA-PSS
#if ALG_RSA

#  define OID_MGF1_VALUE PKCS1_ALG, 8  // 1.2.840.113549.1.1.8
MAKE_OID(_MGF1);

#  define OID_RSAPSS_VALUE PKCS1_ALG, 10  // 1.2.840.113549.1.1.10
MAKE_OID(_RSAPSS);

// This is the OID to designate the public part of an RSA key.
#  define OID_PKCS1_PUB_VALUE PKCS1_ALG, 1  // 1.2.840.113549.1.1.1
MAKE_OID(_PKCS1_PUB);

// These are used for RSA PKCS1 signature Algorithms
#  define OID_PKCS1_SHA1_VALUE PKCS1_ALG, 5  // 1.2.840.113549.1.1.5
#  if ALG_SHA1
MAKE_OID(_PKCS1_SHA1);
#  endif  // ALG_SHA1

#  define OID_PKCS1_SHA256_VALUE PKCS1_ALG, 11  // 1.2.840.113549.1.1.11
#  if ALG_SHA256
MAKE_OID(_PKCS1_SHA256);
#  endif  // ALG_SHA256

#  define OID_PKCS1_SHA384_VALUE PKCS1_ALG, 12  // 1.2.840.113549.1.1.12
#  if ALG_SHA384
MAKE_OID(_PKCS1_SHA384);
#  endif  // ALG_SHA384

#  define OID_PKCS1_SHA512_VALUE PKCS1_ALG, 13  // 1.2.840.113549.1.1.13
#  if ALG_SHA512
MAKE_OID(_PKCS1_SHA512);
#  endif  // ALG_SHA512

#  define OID_PKCS1_SM3_256_VALUE SM_SCHEME, 0x83, 0x78  // 1.2.156.10197.1.504
#  if ALG_SM3_256
MAKE_OID(_PKCS1_SM3_256);
#  endif  // ALG_SM3_256

#  define OID_PKCS1_SHA3_256_VALUE NIST_SIG, 14  // 2.16.840.1.101.3.4.3.14
#  if ALG_SHA3_256
MAKE_OID(_PKCS1_SHA3_256);
#  endif  // ALG_SHA3_256

#  define OID_PKCS1_SHA3_384_VALUE NIST_SIG, 15  // 2.16.840.1.101.3.4.3.15
#  if ALG_SHA3_384
MAKE_OID(_PKCS1_SHA3_384);
#  endif  // ALG_SHA3_384

#  define OID_PKCS1_SHA3_512_VALUE NIST_SIG, 16  // 2.16.840.1.101.3.4.3.16
#  if ALG_SHA3_512
MAKE_OID(_PKCS1_SHA3_512);
#  endif  // ALG_SHA3_512

#endif  // ALG_RSA

#if ALG_ECDSA

#  define OID_ECDSA_SHA1_VALUE 0x06, 0x07, ANSI_X962, 4, 1  // 1.2.840.10045.4.1
#  if ALG_SHA1
MAKE_OID(_ECDSA_SHA1);
#  endif  // ALG_SHA1

#  define OID_ECDSA_SHA256_VALUE ECDSA_SHA2, 2  // 1.2.840.10045.4.3.2
#  if ALG_SHA256
MAKE_OID(_ECDSA_SHA256);
#  endif  // ALG_SHA256

#  define OID_ECDSA_SHA384_VALUE ECDSA_SHA2, 3  // 1.2.840.10045.4.3.3
#  if ALG_SHA384
MAKE_OID(_ECDSA_SHA384);
#  endif  // ALG_SHA384

#  define OID_ECDSA_SHA512_VALUE ECDSA_SHA2, 4  // 1.2.840.10045.4.3.4
#  if ALG_SHA512
MAKE_OID(_ECDSA_SHA512);
#  endif  // ALG_SHA512

#  define OID_ECDSA_SM3_256_VALUE SM_SCHEME, 0x83, 0x75  // 1.2.156.10197.1.501
#  if ALG_SM3_256
MAKE_OID(_ECDSA_SM3_256);
#  endif  // ALG_SM3_256

#  define OID_ECDSA_SHA3_256_VALUE NIST_SIG, 10  // 2.16.840.1.101.3.4.3.10
#  if ALG_SHA3_256
MAKE_OID(_ECDSA_SHA3_256);
#  endif  // ALG_SHA3_256

#  define OID_ECDSA_SHA3_384_VALUE NIST_SIG, 11  // 2.16.840.1.101.3.4.3.11
#  if ALG_SHA3_384
MAKE_OID(_ECDSA_SHA3_384);
#  endif  // ALG_SHA3_384

#  define OID_ECDSA_SHA3_512_VALUE NIST_SIG, 12  // 2.16.840.1.101.3.4.3.12
#  if ALG_SHA3_512
MAKE_OID(_ECDSA_SHA3_512);
#  endif  // ALG_SHA3_512

#endif  // ALG_ECDSA

#if ALG_ECC

#  define OID_ECC_PUBLIC_VALUE 0x06, 0x07, ANSI_X962, 2, 1  // 1.2.840.10045.2.1
MAKE_OID(_ECC_PUBLIC);

#  define OID_ECC_NIST_P192_VALUE PRIME_CURVES, 1  // 1.2.840.10045.3.1.1
#  if ECC_NIST_P192
MAKE_OID(_ECC_NIST_P192);
#  endif  // ECC_NIST_P192

#  define OID_ECC_NIST_P224_VALUE CERTICOM_CURVES, 33  // 1.3.132.0.33
#  if ECC_NIST_P224
MAKE_OID(_ECC_NIST_P224);
#  endif  // ECC_NIST_P224

#  define OID_ECC_NIST_P256_VALUE PRIME_CURVES, 7  // 1.2.840.10045.3.1.7
#  if ECC_NIST_P256
MAKE_OID(_ECC_NIST_P256);
#  endif  // ECC_NIST_P256

#  define OID_ECC_NIST_P384_VALUE CERTICOM_CURVES, 34  // 1.3.132.0.34
#  if ECC_NIST_P384
MAKE_OID(_ECC_NIST_P384);
#  endif  // ECC_NIST_P384

#  define OID_ECC_NIST_P521_VALUE CERTICOM_CURVES, 35  // 1.3.132.0.35
#  if ECC_NIST_P521
MAKE_OID(_ECC_NIST_P521);
#  endif  // ECC_NIST_P521

// No OIDs defined for these anonymous curves
#  define OID_ECC_BN_P256_VALUE 0x00
#  if ECC_BN_P256
MAKE_OID(_ECC_BN_P256);
#  endif  // ECC_BN_P256

#  define OID_ECC_BN_P638_VALUE 0x00
#  if ECC_BN_P638
MAKE_OID(_ECC_BN_P638);
#  endif  // ECC_BN_P638

#  define OID_ECC_SM2_P256_VALUE SM_SCHEME, 0x82, 0x2D  // 1.2.156.10197.1.301
#  if ECC_SM2_P256
MAKE_OID(_ECC_SM2_P256);
#  endif  // ECC_SM2_P256

#  if ECC_BN_P256
#    define OID_ECC_BN_P256 NULL
#  endif  // ECC_BN_P256

#endif  // ALG_ECC

#define OID_SIZE(OID) (OID[1] + 2)

#endif  // !_OIDS_H_
