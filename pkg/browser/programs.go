package browser

import (
    "github.com/m1dugh/program-browser/pkg/bugcrowd"
    . "github.com/m1dugh/program-browser/pkg/types"
    "sync"
)

type Options struct {
    BugcrowdOptions *bugcrowd.Options   `yaml:"bugcrowd-options"`
    Bugcrowd        bool                `yaml:"bugcrowd"`
}


func DefaultOptions() *Options {
    return &Options{
        Bugcrowd: true,
        BugcrowdOptions: nil,
    }
}

type ProgramBrowser struct {
    requesters []ProgramRequester
    Options *Options
}

func New(options *Options) *ProgramBrowser {
    if options == nil {
        options = DefaultOptions()
    }
    requesters := make([]ProgramRequester, 0)

    if options.Bugcrowd {
        requesters = append(requesters, bugcrowd.New(options.BugcrowdOptions))
    }

    return &ProgramBrowser{
        Options: options,
        requesters: requesters,
    }
}

func (browser *ProgramBrowser) GetPrograms() ([]*Program, error) {

    var results []*Program = make([]*Program, 0)

    var mut sync.Mutex
    var wg sync.WaitGroup
    for _, requester := range browser.requesters {
        wg.Add(1)
        go func(results *[]*Program, mut *sync.Mutex, wg *sync.WaitGroup) {
            defer wg.Done()
            programs, err := requester.GetPrograms()
            if err != nil {
                return
            }
            mut.Lock()
            *results = append(*results, programs...)
            mut.Unlock()
        }(&results, &mut, &wg)
    }

    wg.Wait()

    return results, nil
}

