package programs

import (
    "fmt"
    "errors"
    "github.com/m1dugh/program-browser/internal/bugcrowd"
    . "github.com/m1dugh/program-browser/pkg/types"
)

type Options struct {
    
}

type ProgramBrowser struct {
    requesters []*ProgramRequester
    Options *Options
}

func New(options *Options) *ProgramBrowser {
    return &ProgramBrowser{
        Options: options,
        requesters: nil,
    }
}

func (browser *ProgramBrowser) GetPrograms() ([]*Program, error) {

    options := bugcrowd.DefaultOptions()
    options.MaxPages = 1
    options.FetchTargets = false

    requester := bugcrowd.NewBugcrowdRequester(options)

    results, err := requester.GetPrograms()
    if err != nil {
        return nil, errors.New(fmt.Sprintf("GetPrograms: error while fetching bugcrowd programs: %s", err))
    }
    if len(results) == 0 {
        return nil, errors.New(fmt.Sprintf("GetPrograms: could not find any programs"))
    }

    return results, nil
}

