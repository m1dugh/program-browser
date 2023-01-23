# program-browser
A utility to fetch latest bug bounty programs on various platforms

## features

 - allows to fetch programs from bugcrowd 
 - create a burpsuite compliant scope from targets fetched.


## Install

This program was initially designed as a `GO` library, although it is possible
to use the executables, they are more a showcase of what the lib can do.

To install it in your go project, run the usual `go get` command
```shell
$ go get github.com/m1dugh/program-browser@latest
```

## Getting started

```
package main

import (
    "github.com/m1dugh/program-browser/pkg/browser"
    "github.com/m1dugh/program-browser/pkg/types"
    "log"
)

func main() {
    options := browser.DefaultOptions()

    // returns all the programs that can be found
    results, err := browser.GetPrograms()

    // or
    // returns the programs corresponding to the search term
    // results, err := browser.SearchPrograms("example.com")

    if err != nil {
        log.Fatal(err)
    }

    for _, program := range programs {
        // printing the name
        fmt.Println(program.Name)

        // Extracting the asset scope for API testing and website testing
        webScope := program.GetScope(types.website, types.API)

        // urls are all the urls provided in the scope for the web targets
        // domains are all the subdomains provided in the scope
        // for the web targets
        // urls, domains := webScope.ExtractInfo()

        // Extracting all the targets labelled github
        githubScope := program.GetScope(types.GitHub)
        urls, _ := githubScope.ExtractInfo()

        for _, url := range urls {
            // Prints out all the urls to the provided github repos
            fmt.Println(url)
        }

    }
}
```

For further usage, you can take a look at the [types folder](https://github.com/m1dugh/program-browser/tree/master/pkg/types)
for the usage of the `Scope`.


## Upcoming Features
 - [ ] Link with HackerOne
 - [ ] Link with Intigriti
 - [ ] Proper executable

## Contributing
Feel free to open an issue and make a PR for things to be improved on the
project.

## Licence
This project is under [GNU GPL3 Licence](https://github.com/m1dugh/program-browser/blob/master/LICENSE)
