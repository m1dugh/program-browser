package config

import (
	"regexp"
	"strings"

	"github.com/m1dugh/program-browser/internal/utils"
)

type NameFilter struct {
	Exact string `yaml:"exact"`
	Regex string `yaml:"regex"`
	Glob string `yaml:"glob"`
	Insensitive bool `yaml:"insensitive"`

	compiledRegex *regexp.Regexp
}

func (filter *NameFilter) configure() error {

	var err error
	
	if filter.Insensitive {
		filter.Regex = strings.ToLower(filter.Regex)
		filter.Glob = strings.ToLower(filter.Glob)
		filter.Exact = strings.ToLower(filter.Exact)
	}

	if len(filter.Regex) > 0 {
		if filter.compiledRegex, err = regexp.Compile(filter.Regex); err != nil {
			return err
		}
	}

	return nil
}

func (filter *NameFilter) CheckName(s string) bool {
	if filter.Insensitive {
		s = strings.ToLower(s)
	}

	if filter.compiledRegex != nil {
		return filter.compiledRegex.MatchString(s)
	}

	if len(filter.Glob) != 0 {
		return utils.FnMatch(filter.Glob, s)
	}
	
	if len(filter.Exact) != 0 {
		return filter.Exact == s
	}

	return false
}
