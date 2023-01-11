package types

import (
    "sync"
    "time"
    "regexp"
    "fmt"
    "github.com/m1dugh/program-browser/internal/utils"
)


const (
    Private string  = "private"
    Public          = "public"
)


const (
    Partial string  = "partial"
    Unsafe          = "unsafe"
    Safe            = "safe"
)

type Target struct {
    InScope bool `json:"in_scope"`
    URIs []string`json:"uris"`
    Category string `json:"category"`
    Tags []string   `json:"tags"`
}

type Program struct {
    Id string `json:"id"`
    Name string `json:"name"`
    Platform string `json:"platform"`
    PlatformUrl string `json:"url"`
    Status string `json:"status"`
    SafeHarborStatus string `json:"safe_harbor_status"`
    Managed bool `json:"managed"`
    Started bool `json:"started"`
    Ended bool `json:"ended"`
    Category string `json:"category"`
    MinReward uint `json:"min_reward"`
    MaxReward uint `json:"max_reward"`
    Targets []Target `json:"targets"`
    Logo string `json:"logo_url"`
}

func (prog *Program) Code() string {
    return fmt.Sprintf("%s-%s", prog.Platform, prog.Name)
}

type Scope struct {
    Include []string    `json:"include"`
    Exclude []string    `json:"exclude"`
}

func (p *Program) GetScope(category *regexp.Regexp) *Scope {
    res := &Scope{}
    var inc []string
    var ex []string
    for _, t := range p.Targets {
        if category != nil && !category.MatchString(t.Category) {
            continue
        }
        if t.InScope {
            for _, uri := range t.URIs {
                inc = append(inc, uri)
            }
        } else {
            for _, uri := range t.URIs {
                ex = append(ex, uri)
            }
        }
    }
    res.Include = inc
    res.Exclude = ex
    return res
}

/// Returns a StringSet of urls and a StringSet of Subdomains
func (scope *Scope) ExtractInfo() (*utils.StringSet, *utils.StringSet) {
    urls := utils.NewStringSet(nil)
    subdomains := utils.NewStringSet(nil)
    for _, s := range scope.Include {
        url := utils.URLMatcher.FindString(s)
        if len(url) > 0 {
            urls.AddWord(url)
        }
        subdomain := utils.SubdomainMatcher.FindString(s)
        if len(subdomain) > 0 {
            subdomains.AddWord(subdomain)
        }
    }

    return urls, subdomains
}

type RequestThrottler struct {
    MaxRequests int
    _requests int
    _lastFlush int64
    mut *sync.Mutex
}

func NewRequestThrottler(maxRequests int) *RequestThrottler {
    res := &RequestThrottler{
        MaxRequests: maxRequests,
        _requests: 0,
        _lastFlush: time.Now().UnixMicro(),
        mut: &sync.Mutex{},
    }
    return res
}

func (r *RequestThrottler) AskRequest() {
    if r.MaxRequests < 0 {
        return
    }
    r.mut.Lock()
    defer r.mut.Unlock()

    timeStampMicro := time.Now().UnixMicro()
    delta := timeStampMicro - r._lastFlush
    if delta > 1000000 {
        r._requests = 1
        r._lastFlush = timeStampMicro
    } else {
        if r._requests < r.MaxRequests {
            r._requests++
        } else {
            for timeStampMicro - r._lastFlush < 1000000 {
                time.Sleep(time.Microsecond)
                timeStampMicro = time.Now().UnixMicro()
            }

            r._requests = 1
            r._lastFlush = timeStampMicro
        }
    }
}


