package hackerone

import (
    "net/http"
    "github.com/m1dugh/program-browser/pkg/types"
    "encoding/json"
    "errors"
    "fmt"
    "io"
    "bytes"
)

const URL = "https://hackerone.com/graphql"

type HackerOneOptions struct {
    MaxRequests int
}

type HackerOneRequester struct {
    Client *http.Client
    throttler *types.RequestThrottler
    Options *HackerOneOptions
}

func New(options *HackerOneOptions) *HackerOneRequester {
    return &HackerOneRequester{
        Client: &http.Client{},
        throttler: types.NewRequestThrottler(-1),
        Options: options,
    }
}

type h1DirectoryRequest struct {
    OperationName string `json:"operationName"`
    Query string `json:"query"`
    Variables interface{} `json:"variables"`
}

func (r *HackerOneRequester) requestPrograms(request h1DirectoryRequest) (map[string]interface{}, error) {

    payload, err := json.Marshal(request)
    if err != nil {
        return nil, errors.New("Could not marshal json payload")
    }

    response, err := r.Client.Post(URL, "application/json", bytes.NewReader(payload))
    if err != nil || response.StatusCode != 200 {
        return nil, errors.New(fmt.Sprintf("Could not request %s, received %d", URL, response.StatusCode))
    }

    body, err := io.ReadAll(response.Body)
    if err != nil {
        return nil, errors.New("Could not read the body")
    }

    var teams map[string]interface{}
    err = json.Unmarshal(body, &teams)
    if err != nil {
        return nil, errors.New("Could not serialize body as json")
    }

    return teams, nil
}

func (r *HackerOneRequester) GetPrograms() ([]*types.Program, error) {

    return nil, nil
}
