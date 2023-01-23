package types

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
    "errors"
    "github.com/m1dugh/program-browser/internal/utils"
)

type ScopeEntry struct {
    Advanced bool       `json:"-"`
    Enabled bool        `json:"enabled"`
    Host string         `json:"host,omitempty"`
    Protocol string     `json:"protocol,omitempty"`
    File string         `json:"file,omitempty"`
    URL string          `json:"url,omitempty"`
    hostReg *regexp.Regexp
    protocolReg *regexp.Regexp
    fileReg *regexp.Regexp
    urlReg *regexp.Regexp
}

func (s *ScopeEntry) IsEnabled() bool {
    return s.Enabled
}

func (s *ScopeEntry) Setup(advanced bool) error {
    s.Advanced = advanced
    if advanced {
        if len(s.Protocol) == 0 {
            s.Protocol = "any"
        }

        protocol := s.Protocol
        if strings.ToLower(protocol) == "any" {
            protocol = `^[a-z]{2,7}$`
        }
        reg, err := regexp.Compile(protocol)
        if err != nil {
            return err
        }
        s.protocolReg = reg

        if len(s.Host) > 0 {
            reg, err := regexp.Compile(s.Host)
            if err != nil {
                return err
            }
            s.hostReg = reg
        }

        if len(s.File) > 0 {
            reg, err := regexp.Compile(s.File)
            if err != nil {
                return err
            }
            s.fileReg = reg
        }
    } else {
        if len(s.URL) > 0 {
            reg, err := regexp.Compile(s.URL)
            if err != nil {
                return err
            }
            s.urlReg = reg
        }
    }

    return nil
}

/// returns the url corresponding to the regex
func (s *ScopeEntry) ToURL() string {
    var url string
    if !s.Advanced {
        url = s.URL
    } else {
        if s.protocolReg != nil {
            url = fmt.Sprintf("%s://%s%s", s.Protocol, s.Host, s.File)
        } else {
            url = s.Host + s.File
        }
    }

    url = strings.ReplaceAll(url, "\\", "")
    url = strings.ReplaceAll(url, "*.", "")
    url = strings.ReplaceAll(url, "^", "")
    url = strings.ReplaceAll(url, "$", "")

    return url
}

func (s *ScopeEntry) IsValid(host, protocol, file string) bool {
    if s.Advanced {
        if s.hostReg != nil && !s.hostReg.MatchString(host) {
            return false
        }

        if s.protocolReg != nil && !s.protocolReg.MatchString(protocol) {
            return false
        }

        if s.fileReg != nil && !s.fileReg.MatchString(file) {
            return false
        }
    } else {
        url := fmt.Sprintf("%s://%s%s", protocol, host, file)
        if s.urlReg != nil && !s.urlReg.MatchString(url) {
            return false
        }
    }

    return true
}

type Scope struct {
    Advanced bool           `json:"advanced"`
    Exclude []*ScopeEntry   `json:"exclude"`
    Include []*ScopeEntry   `json:"include"`
}

func NewEmptyScope(advanced bool) *Scope {
    scope, _ := NewScope(make([]*ScopeEntry, 0), make([]*ScopeEntry, 0), advanced)
    return scope
}

func NewSimpleScope(include []string, exclude []string) (*Scope, error) {
    scope, err := NewScope(make([]*ScopeEntry, 0, len(include)), make([]*ScopeEntry, 0, len(exclude)), false)
    if err != nil {
        return nil, err
    }

    for _, s := range include {
        entry := &ScopeEntry{
            Enabled: true,
            URL: s,
        }
        err = scope.AddRule(entry, true)
    }

    for _, s := range exclude {
        entry := &ScopeEntry{
            Enabled: true,
            URL: s,
        }
        err = scope.AddRule(entry, false)
    }

    return scope, err
}

type burpTarget struct {
    Scope *Scope `json:"scope"`
}

type burpScope struct {
    Scope   *Scope      `json:"scope,omitempty"`
    Target  *burpTarget `json:"target"`
}

func (s *Scope) setup() error {
    for _, entry := range s.Include {
        err := entry.Setup(s.Advanced)
        if err != nil {
            return err
        }
    }

    for _, entry := range s.Exclude {
        err := entry.Setup(s.Advanced)
        if err != nil {
            return err
        }
    }

    return nil
}

func (s *burpScope) getScope() (*Scope, error) {
    var res *Scope
    if s.Target != nil && s.Target.Scope != nil {
        res = s.Target.Scope
    } else {
        res = s.Scope
    }

    err := res.setup()
    if err != nil {
        return nil, err
    }

    return res, nil

}

func DeserializeScope(body []byte) (*Scope, error) {
    var scope burpScope
    err := json.Unmarshal(body, &scope)
    if err != nil {
        return nil, err
    }

    res, err := scope.getScope()
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Could not deserialize scope: %s", err))
    }
    return res, nil
}

func NewScope(include []*ScopeEntry, exclude []*ScopeEntry, advanced bool) (*Scope, error) {
    res := &Scope{
        Exclude: exclude,
        Include: include,
        Advanced: advanced,
    }

    err := res.setup()
    if err != nil {
        return nil, err
    }
    return res, nil
}

func (s *Scope) AddRule(entry *ScopeEntry, in bool) error {
    if in {
        s.Include = append(s.Include, entry)
    } else {
        s.Exclude = append(s.Exclude, entry)
    }
    return entry.Setup(s.Advanced)
}

/// splits a URL into 3 parts
/// The host, the protocol, and the file
/// if one of the parts is missing, the corresponding string will be empty
func splitURL(url string) (string, string, string) {
    var protocol, host, file string
    splits := strings.SplitN(url, "://", 2)
    if len (splits) <= 1 {
        splits = strings.SplitN(splits[0], "/", 2)
    } else {
        protocol = splits[0] 
        splits = strings.SplitN(splits[1], "/", 2)
    }

    host = splits[0]
    if len (splits) > 1 {
        file = fmt.Sprintf("/%s", splits[1])
    }

    return host, protocol, file
}

func (s *Scope) AddSimpleRule(url string, in bool) error {
    var entry *ScopeEntry
    if s.Advanced {
        host, protocol, file := splitURL(url)
        if len(file) > 0 {
            file = "^" + file
        }
        if len(protocol) > 0 {
            protocol = protocol + "$"
        }
        if len(host) > 0 {
            host = fmt.Sprintf("^%s$", host)
        }

        entry = &ScopeEntry{
            Enabled: true,
            Host: host,
            Protocol: protocol,
            File: file,
        }

    } else {
        entry = &ScopeEntry{
            Enabled: true,
            URL: url,
        }
    }

    return s.AddRule(entry, in)
}

func (s *Scope) InScope(url string) bool {

    host, protocol, file := splitURL(url)
    if len(protocol) == 0 {
        return false
    }
    if len(file) == 0 {
        file = "/"
    }

    for _, entry := range s.Exclude {
        if entry.IsEnabled() && entry.IsValid(host, protocol, file) {
            return false
        }
    }

    for _, entry := range s.Include {
        if entry.IsEnabled() && entry.IsValid(host, protocol, file) {
            return true
        }
    }

    return false
}

func wildcardToRegex(uri string) string {
    uri = strings.ReplaceAll(uri, ".", "\\.")
    uri = strings.ReplaceAll(uri, "*", ".*")
    uri = strings.ReplaceAll(uri, ".*\\.", ".*\\.?")

    return uri
}

func (p *Program) GetScope(category *regexp.Regexp) *Scope {
    res, _ := NewScope(make([]*ScopeEntry, 0), make([]*ScopeEntry, 0), true)
    for _, t := range p.Targets {
        if category != nil && !category.MatchString(t.Category) {
            continue
        }
        if t.InScope {
            for _, uri := range t.URIs {
                res.AddSimpleRule(wildcardToRegex(uri), true)
            }
        } else {
            for _, uri := range t.URIs {
                res.AddSimpleRule(wildcardToRegex(uri), false)
            }
        }
    }

    return res
}

/// returns a list of bytes representing the Burp Suite corresponding
/// scope
func (s *Scope) ToBurpScope() ([]byte, error) {
    target := &burpTarget{
        Scope: s,
    }

    burpScope := burpScope{
        Target: target,
    };

    body, err := json.MarshalIndent(burpScope, "", "\t")
    if err != nil {
        return nil, err
    }
    return body, nil

}

/// Returns a StringSet of urls and a StringSet of Subdomains
func (scope *Scope) ExtractInfo() (*utils.StringSet, *utils.StringSet) {
    urls := utils.NewStringSet(nil)
    subdomains := utils.NewStringSet(nil)

    for _, s := range scope.Include {
        url := utils.URLMatcher.FindString(s.ToURL())
        if len(url) > 0 {
            urls.AddWord(url)
        }
        subdomain := utils.SubdomainMatcher.FindString(s.ToURL())
        if len(subdomain) > 0 {
            subdomains.AddWord(subdomain)
        }
    }

    return urls, subdomains
}
