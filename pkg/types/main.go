package types

import (
	"strings"
)

type Scheme string

type Host struct {
	Suffix   string `json:"prefix"`
	Wildcard bool   `json:"wildcard"`
}

type Endpoint struct {
	Scheme *string `json:"scheme"`
	Host   Host    `json:"host"`
	Path   *string `json:"path,omitempty"`
}

func NewEndpointFromString(s string) Endpoint {
	var res Endpoint
	index := strings.Index(s, "://")
	if index != -1 {
		protocol := s[:index]
		res.Scheme = &protocol
		s = s[index+3:]
	}

	index = strings.Index(s, "/")

	if index != -1 {
		path := s[index+1:]
		res.Path = &path
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
	if pattern.Scheme != nil && (ep.Scheme == nil || *pattern.Scheme != *ep.Scheme) {
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

	if pattern.Path != nil {
		if ep.Path == nil || !strings.HasPrefix(*ep.Path, *pattern.Path) {
			return false
		}
	}
	return true
}

func (ep *Endpoint) ToString() string {
	var result strings.Builder
	if ep.Scheme != nil {
		result.WriteString(string(*ep.Scheme))
		result.WriteString("://")
	}
	if ep.Host.Wildcard {
		result.WriteString("*.")
	}
	result.WriteString(ep.Host.Suffix)
	if ep.Path != nil {
		result.WriteString(*ep.Path)
	}

	return result.String()
}

type Scope struct {
	AllowedEndpoints []Endpoint `json:"allowedEndpoints"`
	DeniedEndpoints  []Endpoint `json:"deniedEndpoints"`
}

type Program struct {
	Platform   string `json:"platform"`
	PlatformId string `json:"platformId"`
	Scope      Scope  `json:"scope"`
	Name       string `json:"name"`
	URL        string `json:"URL"`
}
