package main

import (
    "github.com/m1dugh/program-browser/pkg/bugcrowd"
    programs "github.com/m1dugh/program-browser/pkg/program-browser"
    "log"
    "regexp"
    "fmt"
)

var webRegexp = regexp.MustCompile("[Ww]ebsite")

func main() {
    boptions := bugcrowd.DefaultOptions()
    boptions.MaxPrograms = 4
    boptions.FetchTargets = true
    options := &programs.Options{
        Bugcrowd: true,
        BugcrowdOptions: boptions,
    }

    browser := programs.New(options)

    results, err := browser.GetPrograms()
    if err != nil {
        log.Fatal(err)
    }

    for _, program := range results {
        scope := program.GetScope(webRegexp)
        _, domains := scope.ExtractInfo()
        for _, domain := range domains.ToArray() {
            fmt.Println(domain)
        }
    }
}
