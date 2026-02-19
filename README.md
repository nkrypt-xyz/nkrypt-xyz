# nkrypt.xyz

High-performance, secure web server for encrypted file storage with client-side encryption.

## Deployment

Choose your deployment method:

- **[All-in-One (Docker/Podman)](./infrastructure/aio/)** - Single-machine setup for home servers and VPS
- **[Bare Metal](./infrastructure/bare-metal/)** - Low-overhead deployment on Linux
- **[Enterprise (Kubernetes)](./infrastructure/enterprise/)** - Production-grade Terraform + Kubernetes

## Development

See **[Development Guide](./devnotes/)** for local setup and contribution guidelines.

## Stack

**Application**: Go 1.21+, PostgreSQL 16, Redis 7, MinIO  
**Infrastructure**: Terraform, Kubernetes, Helm, GitHub Actions

## License

[GNU General Public License v3.0](LICENSE) © 2022 [Sayem Shafayet](https://sayemshafayet.com)
