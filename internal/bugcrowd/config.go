package bugcrowd

import "github.com/m1dugh/program-browser/internal/config"

type Options struct {
	Filters []config.NameFilter
}

func DefaultOptions() *Options {
	return &Options{}
}
