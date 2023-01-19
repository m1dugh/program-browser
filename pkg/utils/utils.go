package utils

import (
	programs "github.com/m1dugh/program-browser/pkg/program-browser"
	yaml "gopkg.in/yaml.v3"
)

func DeserializeOptions(body []byte) (*programs.Options, error) {
    options := &programs.Options{}
    err := yaml.Unmarshal(body, options)

    if err != nil {
        return nil, err
    }

    return options, nil
}
