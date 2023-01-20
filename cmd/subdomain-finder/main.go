package main

import (
    "github.com/m1dugh/program-browser/pkg/utils"
    programs "github.com/m1dugh/program-browser/pkg/browser"
    "log"
    "regexp"
    "fmt"
    "flag"
    "os"
    "errors"
)

var webRegexp = regexp.MustCompile("[Ww]ebsite")

func main() {

    var settingsPath string
    flag.StringVar(&settingsPath, "settings", "", "The path to the settings.yaml file")

    flag.Parse()

    var options *programs.Options
    if len(settingsPath) == 0 {
        options = programs.DefaultOptions()
    } else if _, err := os.Stat(settingsPath); errors.Is(err, os.ErrNotExist) {
        log.Fatal(fmt.Sprintf("File not found: %s", settingsPath))
    } else {
        body, err := os.ReadFile(settingsPath)
        if err != nil {
            log.Fatal(err)
        }

        options, err = utils.DeserializeOptions(body)
        if err != nil {
            log.Fatal(err)
        }

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
