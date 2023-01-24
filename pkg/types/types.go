package types

import (
	"fmt"
	"time"
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

type TargetCategory string

const (
    Website TargetCategory  = "website"
    API                     = "api"
    GitHub                  = "GitHub"
    SocialMedias            = "social medias"
    IOS                     = "IOS"
    Android                 = "Android"
    Others                  = "others"
)

type Target struct {
    InScope     bool            `json:"in_scope"`
    URIs        []string        `json:"uris"`
    Category    TargetCategory  `json:"category"`
    Tags        []string        `json:"tags"`
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
    StartedAt time.Time `json:"started_at"`
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


type ProgramRequester interface {
    GetPrograms() ([]*Program, error)
    SearchPrograms(name string) ([]*Program, error)
}
