package affiliate

import (
	"fmt"
	"net/url"
	"strings"
)

// Rule describes a safe affiliate URL rewrite for a known source.
type Rule struct {
	SourceID       string
	AllowedDomains []string
	Param          string
	Value          string
	Enabled        bool
}

// Result carries both canonical evidence URL and user-facing presentation URL.
type Result struct {
	CanonicalURL string `json:"canonicalUrl"`
	PresentedURL string `json:"presentedUrl"`
	Applied      bool   `json:"applied"`
	Reason       string `json:"reason,omitempty"`
}

// ApplyQueryParam applies a simple affiliate query-param rule to safe product URLs.
func ApplyQueryParam(raw string, rule Rule) (Result, error) {
	res := Result{CanonicalURL: raw, PresentedURL: raw}
	if !rule.Enabled {
		res.Reason = "affiliate_rule_disabled"
		return res, nil
	}
	if rule.Param == "" || rule.Value == "" {
		return res, fmt.Errorf("affiliate rule for %s missing param/value", rule.SourceID)
	}
	u, err := url.Parse(raw)
	if err != nil {
		return res, err
	}
	if u.Scheme != "https" {
		res.Reason = "unsupported_url_scheme"
		return res, nil
	}
	if !domainAllowed(u.Hostname(), rule.AllowedDomains) {
		res.Reason = "domain_not_allowed"
		return res, nil
	}
	if unsafePath(u.Path) {
		res.Reason = "unsafe_url_path"
		return res, nil
	}
	q := u.Query()
	q.Set(rule.Param, rule.Value)
	u.RawQuery = q.Encode()
	res.PresentedURL = u.String()
	res.Applied = true
	return res, nil
}

func domainAllowed(host string, allowed []string) bool {
	host = strings.ToLower(host)
	for _, domain := range allowed {
		domain = strings.ToLower(domain)
		if host == domain || strings.HasSuffix(host, "."+domain) {
			return true
		}
	}
	return false
}

func unsafePath(path string) bool {
	p := strings.ToLower(path)
	for _, bad := range []string{"/cart", "/checkout", "/login", "/account", "/payment"} {
		if strings.Contains(p, bad) {
			return true
		}
	}
	return false
}
