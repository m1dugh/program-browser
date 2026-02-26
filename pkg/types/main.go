package types

import (
	"strings"
)

type Scheme string

type Host struct {
	Suffix   string `json:"suffix"`
	Wildcard bool   `json:"wildcard"`
}

type Endpoint struct {
	Scheme string `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	Host   Host   `json:"host" yaml:"host"`
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`
}

func NewEndpointFromString(s string) Endpoint {
	var res Endpoint
	index := strings.Index(s, "://")
	if index != -1 {
		protocol := s[:index]
		res.Scheme = protocol
		s = s[index+3:]
	}

	index = strings.Index(s, "/")

	if index != -1 {
		path := s[index+1:]
		res.Path = path
		s = s[:index]
	}

	if strings.HasPrefix(s, "*.") {
		res.Host.Wildcard = true
		s = s[2:]
	}
	res.Host.Suffix = s

	return res
}

func (ep *Endpoint) CompatibleWith(pattern *Endpoint) bool {
	if pattern.Scheme != "" && (ep.Scheme == "" || pattern.Scheme != ep.Scheme) {
		return false
	}

	if pattern.Host.Wildcard {
		if !strings.HasSuffix(ep.Host.Suffix, pattern.Host.Suffix) {
			return false
		}
	} else {
		if ep.Host.Suffix != pattern.Host.Suffix {
			return false
		}
	}

	if pattern.Path != "" {
		if ep.Path == "" || !strings.HasPrefix(ep.Path, pattern.Path) {
			return false
		}
	}
	return true
}

func (ep *Endpoint) ToString() string {
	var result strings.Builder
	if ep.Scheme != "" {
		result.WriteString(ep.Scheme)
		result.WriteString("://")
	}
	if ep.Host.Wildcard {
		result.WriteString("*.")
	}
	result.WriteString(ep.Host.Suffix)
	if ep.Path != "" {
		result.WriteString(ep.Path)
	}

	return result.String()
}

type Scope struct {
	AllowedEndpoints []Endpoint `json:"allowed_endpoints" yaml:"allowed_endpoints"`
	DeniedEndpoints  []Endpoint `json:"denied_endpoints" yaml:"denied_endpoints"`
}

type Program struct {
	Platform   string `json:"platform" yaml:"platform"`
	PlatformId string `json:"platform_id" yaml:"platform_id"`
	Scope      Scope  `json:"scope" yaml:"scope"`
	Name       string `json:"name" yaml:"name"`
	URL        string `json:"url" yaml:"url"`
}
