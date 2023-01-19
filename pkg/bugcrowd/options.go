package bugcrowd

const (
    DefaultSkipScope = false
    DefaultMaxPrograms = -1
    DefaultSort = "starts-desc"
    DefaultMaxRequests = 5
    DefaultHidden = false
)

type Options struct {
    SkipScope bool      `yaml:"skip-scope"`
    MaxPrograms int     `yaml:"max-programs"`
    Sort string         `yaml:"sort"`
    MaxRequests int     `yaml:"max-requests"`
    Hidden bool         `yaml:"hidden"`
}

func DefaultOptions() *Options {
    return &Options {
        SkipScope: DefaultSkipScope,
        MaxPrograms: DefaultMaxPrograms,
        Sort: DefaultSort,
        MaxRequests: DefaultMaxRequests,
        Hidden: DefaultHidden,
    }
}
