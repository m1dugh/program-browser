package types

import (
    "testing"
    "fmt"
)


func TestSimpleScope(t *testing.T) {
    
    includes := []string{
        `^https?://([\w]+\.)*example\.com`,
    }

    scope, err := NewSimpleScope(includes, nil)
    if err != nil {
        t.Errorf("%s", err)
    }

    _testSample(t, scope, "http://example.com", true)

    _testSample(t, scope, "https://www.google.com/?search=https://www.example.com", false)

    _testSample(t, scope, "ftp://example.com", false)
    _testSample(t, scope, "https://www.example.com", true)
}

func TestAdvancedScope(t *testing.T) {
    include := []*ScopeEntry{
        &ScopeEntry{
            Enabled: true,
            Host: `^([\w]+\.)*example.com$`,
            Protocol: `^https?$`,
        },
    }

    exclude := []*ScopeEntry {
        &ScopeEntry{
            Enabled: true,
            Host: `^test.example.com$`,
            Protocol: `^https?$`,
        },
    }

    scope, err := NewScope(include, exclude, true)
    if err != nil {
        t.Errorf("%s", err)
    }
    _testSample(t, scope, "http://example.com", true)

    _testSample(t, scope, "https://www.google.com/?search=https://www.example.com", false)

    _testSample(t, scope, "ftp://example.com", false)
    _testSample(t, scope, "https://www.example.com", true)
    _testSample(t, scope, "https://test.example.com", false)
    _testSample(t, scope, "https://www.test.example.com", true)
}

func TestAddSimpleRule(t *testing.T) {
    scope := NewEmptyScope(true)
    scope.AddSimpleRule(`^https://.*\.google.com/search$`, true)

    _testSample(t, scope, "https://www.google.com/search", true)
    _testSample(t, scope, "https://www.google.com/searche", false)
    _testSample(t, scope, "https://search.google.com/search", true)

}

func _testSample(t *testing.T, scope *Scope, test string, inScope bool) {
    if scope.InScope(test) != inScope {
        if inScope {
            t.Errorf("Expected '%s' to be in scope but got out of scope", test)
        } else {
            t.Errorf("Expected '%s' to be out of scope but got in scope", test)
        }
    }
}

const simpleScope = `
{
    "scope": {
        "advanded": false,
        "include": [
            {
                "enabled": true,
                "url": "hackerone\\.com"
            },
            {
                "enabled": true,
                "url": "^https?://api\\.hackerone\\.com"
            }
        ],
        "exclude": [
            {
                "enabled": true,
                "url": "^https?://docs\\.hackerone\\.com"
            },
            {
                "enabled": true,
                "url": "^http://"
            }
        ]
    }
}
`


func testDeserialize(t *testing.T, body string, advanced bool) {
    scope, err := DeserializeScope([]byte(body))
    if err != nil {
        t.Errorf("Deserialize error: %s", err)
    }

    if scope.Advanced != advanced {
        t.Errorf("expected scope to be not advanced")
    }

    for _, entry := range scope.Include {
        if entry.IsEnabled() != true {
            t.Errorf("Expected entry to be enabled")
        }
    }


    _testSample(t, scope, "https://www.hackerone.com/security", true)
}

func TestDeserializeSimple(t *testing.T) {
    testDeserialize(t, simpleScope, false)
}

const burpAdvancedScope = `
{
    "target": {
        "scope": {
            "advanced": true,
            "include": [
                {
                    "enabled": true,
                    "host": "^www\\.hackerone\\.com$",
                    "protocol": "^https$"
                }
            ],

            "exclude": [
                {
                    "enabled": true,
                    "host": "^docs\\.hackerone\\.com$",
                    "protocol": "any"
                },
                {
                    "enabled": true,
                    "host": "^hackerone\\.com$",
                    "protocol": "^http$"
                }
            ]
        }
    }
}
`
func TestDeserializeBurp(t *testing.T) {
    testDeserialize(t, burpAdvancedScope, true)
}

func TestSimpleEntry(t *testing.T) {
    entry := &ScopeEntry{
        Enabled: true,
        URL: `^https?://hackerone.com`,
    }

    entry.Setup(false)

    if entry.urlReg == nil {
        t.Errorf("Expected url regex to compile")
    }

    var host, protocol, file, url string
    host = "hackerone.com"
    protocol = "https"
    file = "/"
    url = fmt.Sprintf("%s://%s%s", protocol, host, file)

    if !entry.IsValid(host, protocol, file) {
        t.Errorf("Expected `%s` in but go out", url)
    }

    host = "www.hackerone.com"
    url = fmt.Sprintf("%s://%s%s", protocol, host, file)
    if entry.IsValid(host, protocol, file) {
        t.Errorf("Expected `%s` out but go in", url)
    }
}

func TestAdvancedEntry(t *testing.T) {
    entry := &ScopeEntry{
        Enabled: true,
        Host: `^hackerone.com$`,
    }

    entry.Setup(true)

    if entry.hostReg == nil {
        t.Errorf("Expected host regex to compile")
    }

    if entry.fileReg != nil {
        t.Errorf("Expected file regex to be nil")
    }

    if entry.protocolReg != nil {
        t.Errorf("Expected protocol regex to be nil")
    }

    if !entry.IsValid("hackerone.com", "https://", "/") {
        t.Errorf("Expected in but go out")
    }

    if entry.IsValid("www.hackerone.com", "https://", "/") {
        t.Errorf("Expected out but got in")
    }
}


