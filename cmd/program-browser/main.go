package main

import (
    "fmt"
    "log"
    "encoding/json"
    programs "github.com/m1dugh/program-browser/pkg/program-browser"
)

func main() {
    browser := programs.New(nil);

    results, err := browser.GetPrograms()
    if err != nil {
        log.Fatal(err)
    }

    bytes, err := json.Marshal(results)
    fmt.Println(string(bytes))
}
