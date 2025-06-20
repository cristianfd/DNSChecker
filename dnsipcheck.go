package dnschecker

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/miekg/dns"

	"go.uber.org/zap"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
)

func init() {
	caddy.RegisterModule(PermissionByDNS{})
}

type OnDemandConfig struct {
	PermissionRaw json.RawMessage `json:"permission,omitempty" caddy:"namespace=tls.permission inline_key=module"`
}

type OnDemandPermission interface {
	CertificateAllowed(ctx context.Context, name string) error
}

type PermissionByDNS struct {
	TargetIP string `json:"targetip"`
	Resolver string `json:"resolver,omitempty"`

	logger   *zap.Logger
	replacer *caddy.Replacer
}

// CaddyModule returns the Caddy module information.
func (PermissionByDNS) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "tls.permission.dnschecker",
		New: func() caddy.Module { return new(PermissionByDNS) },
	}
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (p *PermissionByDNS) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	if !d.Next() {
		return nil
	}
	for d.NextBlock(0) {
		switch d.Val() {
		case "targetip":
			if !d.NextArg() {
				return d.ArgErr()
			}
			p.TargetIP = d.Val()
		case "resolver":
			if !d.NextArg() {
				return d.ArgErr()
			}
			p.Resolver = d.Val()
		default:
			return d.Errf("unrecognized subdirective %s", d.Val())
		}
	}
	return nil
}

func (p *PermissionByDNS) Provision(ctx caddy.Context) error {
	p.logger = ctx.Logger()
	p.replacer = caddy.NewReplacer()
	return nil
}

func lookupHostWithResolver(domain, resolver string) ([]string, error) {
	c := new(dns.Client)
	m := new(dns.Msg)
	m.SetQuestion(dns.Fqdn(domain), dns.TypeA)
	resp, _, err := c.Exchange(m, resolver+":53")
	if err != nil {
		return nil, err
	}
	var ips []string
	for _, ans := range resp.Answer {
		if a, ok := ans.(*dns.A); ok {
			ips = append(ips, a.A.String())
		}
	}
	return ips, nil
}

func (p PermissionByDNS) CertificateAllowed(ctx context.Context, name string) error {
	p.logger.Debug("checking permission for certificate",
		zap.String("domain", name),
		zap.String("targetIP", p.TargetIP),
		zap.String("resolver", p.getResolver()),
	)

	ips, err := lookupHostWithResolver(name, p.getResolver())
	if err != nil {
		p.logger.Error("DNS lookup failed",
			zap.String("domain", name),
			zap.String("error", err.Error()),
		)

		return fmt.Errorf("failed to resolve DNS for %s: %w", name, err)
	}

	for _, ip := range ips {
		p.logger.Debug("Found DNS A record",
			zap.String("domain", name),
			zap.String("ip", ip),
		)

		if ip == p.TargetIP {
			p.logger.Debug("Domain resolves to allowed IP",
				zap.String("domain", name),
				zap.String("ip", ip),
			)
			return nil
		}
	}

	return fmt.Errorf("domain %s does not resolve to allowed IP %s", name, p.TargetIP)
}

// getResolver returns the configured resolver or the default (1.1.1.1)
func (p PermissionByDNS) getResolver() string {
	if p.Resolver != "" {
		return p.Resolver
	}
	return "1.1.1.1"
}

// ErrPermissionDenied is an error that should be wrapped or returned when the
// configured permission module does not allow a certificate to be issued,
// to distinguish that from other errors such as connection failure.
var ErrPermissionDenied = errors.New("certificate not allowed by permission module")

// Interface guards
var (
	_ OnDemandPermission = (*PermissionByDNS)(nil)
	_ caddy.Provisioner  = (*PermissionByDNS)(nil)
)
