package hackerone

import (
    "testing"
    "log"
    "fmt"
)

func TestGetPrograms(t *testing.T) {
    requester := New(nil)

    variables := directoryRequestVariables{
        First: 25,
    }

    req := h1DirectoryRequest{
        OperationName: "DirectoryQuery",
        Query: DIRECTORY_QUERY,
        Variables: variables,
    }
    teams, err := requester.requestPrograms(req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(teams["teams"])
}
