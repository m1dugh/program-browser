package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"

	programs "github.com/m1dugh/program-browser/pkg/browser"
	"github.com/m1dugh/program-browser/pkg/types"
	"github.com/m1dugh/program-browser/pkg/utils"
)

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
        scope := program.GetScope(types.Website, types.API)
        _, domains := scope.ExtractInfo()
        for _, domain := range domains.ToArray() {
            fmt.Println(domain)
        }
    }
}
