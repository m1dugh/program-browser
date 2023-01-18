package main

import (
    "fmt"
    "log"
    "encoding/json"
    programs "github.com/m1dugh/program-browser/pkg/program-browser"
    "github.com/m1dugh/program-browser/pkg/bugcrowd"
)

func main() {
    boptions := bugcrowd.DefaultOptions()
    boptions.MaxPrograms = 4
    boptions.FetchTargets = true
    options := &programs.Options{
        Bugcrowd: true,
        BugcrowdOptions: boptions,
    }
    browser := programs.New(options);

    results, err := browser.GetPrograms()
    if err != nil {
        log.Fatal(err)
    }

    bytes, err := json.Marshal(results)
    fmt.Println(string(bytes))
}
