package main

import (
    "fmt"
    "log"
    "encoding/json"
    browser "github.com/m1dugh/program-browser/pkg/program-browser"
)

func main() {
    fmt.Println("Hello, World!")
    results, err := browser.GetPrograms();
    if err != nil {
        log.Fatal(err)
    }

    bytes, err := json.Marshal(results)
    fmt.Println(string(bytes))
}
