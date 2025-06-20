# Caddy DNS IP Check Plugin

## Overview

**Caddy DNS IP Check** is a Caddy v2 plugin module that implements the `caddytls.OnDemandPermission` interface.  
It allows Caddy to check, before issuing an on-demand TLS certificate, whether the requested domain's DNS A record resolves to a specific IP address using a configurable DNS resolver (by default, Cloudflare's 1.1.1.1).  
This is useful for restricting certificate issuance to domains that are properly pointed to your server.

## Motivation

Certificate authorities may penalize you or block your account if you request certificates for domains that are not yet pointing to your server.  
This plugin helps prevent such penalties by ensuring that only domains resolving to your server's IP (as seen by a public DNS resolver) are allowed to obtain certificates.

## Features

- Checks if a domain's A record resolves to a configured IP address.
- Uses a configurable DNS resolver for lookups (default: Cloudflare's 1.1.1.1).
- Provides debug logs for troubleshooting.
- Easy configuration via the Caddyfile.

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) 1.24 or newer
- [xcaddy](https://github.com/caddyserver/xcaddy)

### Build Caddy with the Plugin

Build Caddy with the plugin:

```sh
xcaddy build --with github.com/cristianfd/dnschecker
```

## Usage

Add the following to your global Caddyfile:

```caddyfile
{
    on_demand_tls {
        permission dnschecker {
            targetip 1.2.3.4
            resolver 8.8.8.8
        }
    }
}

domain.com {
	tls {
		on_demand
	}
	respond "Hello world!"
}
```

- Replace `targetip` with the IP address you want the domain to resolve to.
- The `resolver` option is optional. If omitted, the plugin uses Cloudflare's 1.1.1.1 by default.

Example with default resolver:

```caddyfile
{
    on_demand_tls {
        permission tls.permission.dnschecker {
            targetip 1.2.3.4
        }
    }
}
```

## How It Works

When a certificate is requested on-demand, the plugin queries the configured DNS resolver (default: 1.1.1.1) for the domain's A records.  
If any of the returned IPs match the configured `targetip`, the certificate issuance is permitted. Otherwise, it is denied.

## License

This project is licensed under the MIT License.  
You are free to use, modify, and distribute this software.  
Attribution is appreciated but not required.

---

**Author:** Cristian Fabian Dourado