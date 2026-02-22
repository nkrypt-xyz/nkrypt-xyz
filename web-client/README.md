# nkrypt.xyz Web Client

Vue 3 + Quasar SPA for encrypted file storage with client-side encryption.

## Quick Start

### Prerequisites

- **Node.js 20+**
- **npm 6.13+**

### Run Development Server

```bash
npm install
npm run build:core   # Build crypto module (required once before dev or build)
npm run dev
```

The client depends on the [nkrypt-xyz-core-nodejs](../nkrypt-xyz-core-nodejs/) crypto module. Run `npm run build:core` once before dev or build. Production `npm run build` runs it automatically.

Client runs on `http://localhost:9042`. Requires the [web server](../web-server/README.md) to be running for full functionality.

### Build

```bash
npm run build
```

Output: `dist/spa/`

## Scripts

| Command | Description |
|---------|-------------|
| `npm run dev` | Start dev server |
| `npm run build` | Build for production |
| `npm run check` | Lint + type-check |

## Documentation

See [dev-docs/](../dev-docs/) for architecture and contribution guidelines.
