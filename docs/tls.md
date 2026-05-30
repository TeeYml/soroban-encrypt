# TLS Configuration

## TLS_MODE Options

| Mode | Description |
|------|-------------|
| `off` | Plain HTTP (default) |
| `auto` | Automatic Let's Encrypt via ACME |
| `manual` | Load cert/key from file paths |

## auto TLS (Let's Encrypt)

```yaml
tls_mode: auto
domain: node1.example.com
acme_cache_dir: /data/acme   # Must be persistent across restarts
```

Requirements:
- Port 443 must be publicly reachable
- Port 80 is used for ACME HTTP-01 challenge (auto-redirected to HTTPS)
- `domain` must match the server's public DNS name

## manual TLS

```yaml
tls_mode: manual
tls_cert: /etc/node/cert.pem
tls_key:  /etc/node/key.pem
```

## Security Posture

- TLS 1.3 minimum enforced
- HSTS header: `max-age=63072000; includeSubDomains`
- HTTP traffic automatically redirected to HTTPS when TLS is active

## Admin Dashboard

The embedded admin dashboard is served at `/admin/` when `ADMIN_DASHBOARD=true`.

> ⚠️ The dashboard is disabled when `TLS_MODE=off` to prevent credentials being sent over plain HTTP.

Set `ADMIN_PASSWORD` to protect the dashboard with HTTP Basic Auth.
