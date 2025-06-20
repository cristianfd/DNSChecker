# Caddy DNS IP Check Plugin

## Overview

**Caddy DNS IP Check** is a Caddy v2 plugin module that implements the `caddytls.OnDemandPermission` interface.  
It allows Caddy to check, before issuing an on-demand TLS certificate, whether the requested domain's DNS A record resolves to a specific IP address (using Cloudflare's DNS resolver at 1.1.1.1).  
This is useful for restricting certificate issuance to domains that are properly pointed to your server.

## Features

- Checks if a domain's A record resolves to a configured IP address.
- Uses Cloudflare's DNS resolver (1.1.1.1) for lookups.
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
        permission dnschecker 1.2.3.4
    }
}
```

- Replace `1.2.3.4` with the IP address you want the domain to resolve to.

## How It Works

When a certificate is requested on-demand, the plugin queries Cloudflare's DNS resolver (1.1.1.1) for the domain's A records.  
If any of the returned IPs match the configured `allowed_ip`, the certificate issuance is permitted. Otherwise, it is denied.

## License

This project is licensed under the MIT License.  
You are free to use, modify, and distribute this software.  
Attribution is appreciated but not required.

---

**Author:** Cristian Fabian Dourado Alvarez