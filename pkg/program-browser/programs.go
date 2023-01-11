package programs

import (
    "fmt"
    "errors"
    "github.com/m1dugh/program-browser/internal/bugcrowd"
    . "github.com/m1dugh/program-browser/pkg/types"
)

func GetPrograms() ([]*Program, error) {
    requester := bugcrowd.NewBugcrowdRequester(5)

    results, err := requester.GetPrograms()
    if err != nil {
        return nil, errors.New(fmt.Sprintf("GetPrograms: error while fetching bugcrowd programs: %s", err))
    }
    if len(results) == 0 {
        return nil, errors.New(fmt.Sprintf("GetPrograms: could not find any programs"))
    }

    return results, nil
}

