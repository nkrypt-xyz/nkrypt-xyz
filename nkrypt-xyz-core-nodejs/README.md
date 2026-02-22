# nkrypt-xyz-core-nodejs

Core crypto and utilities for nkrypt.xyz. Node.js ESM module with full unit test coverage.

## Contents

- **constants** – Crypto specs (IV length, algorithms, chunk sizes)
- **buffer-utils** – Base64 encode/decode for ArrayBuffer/Uint8Array
- **crypto-api-utils** – Build/unbuild crypto header (NK001|iv|salt)
- **crypto-utils** – AES-GCM encryption, PBKDF2 key derivation, encrypt/decrypt text and objects

## Usage

```bash
npm install nkrypt-xyz-core-nodejs
```

```ts
import {
  encryptText,
  decryptText,
  createEncryptionKeyFromPassword,
  buildCryptoHeader,
  unbuildCryptoHeader,
} from "nkrypt-xyz-core-nodejs";
```

## Development

```bash
npm install
npm run build
npm run test
```

**Note:** The web-client depends on this package via `file:../nkrypt-xyz-core-nodejs`. Run `npm run build` in this package before building or running the web-client. The web-client's `npm run build` does this automatically.

## Testing

Uses Vitest. Tests run with `WEAKEN_CRYPTO_FOR_TESTING` for deterministic IV/salt in round-trip tests.
