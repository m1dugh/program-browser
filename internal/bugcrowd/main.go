package bugcrowd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"

	pbTypes "github.com/m1dugh/program-browser/pkg/types"
)

const platform string = "bugcrowd"

type BugcrowdApi struct {
	client  *http.Client
	baseURL string
	options *Options
}

func NewBugcrowdApi(opts *Options) *BugcrowdApi {
	if opts == nil {
		opts = DefaultOptions()
	}
	return &BugcrowdApi{
		client:  &http.Client{},
		baseURL: "https://bugcrowd.com",
		options: opts,
	}
}

func (api *BugcrowdApi) prepareProgramsRequest(req *http.Request, page int) *http.Request {

	req.Header.Add("Content-Type", "application/json")
	query := req.URL.Query()
	query.Add("category", "bug_bounty")
	query.Add("page", strconv.Itoa(page))
	req.URL.RawQuery = query.Encode()

	return req
}

type programsResult struct {
	Engagements []struct {
		Brief string `json:"briefUrl"`
		Name  string `json:"name"`
	} `json:"engagements"`
	Metadata struct {
		Limit int `json:"limit"`
		Count int `json:"totalCount"`
	} `json:"paginationMeta"`
}

func (api *BugcrowdApi) fetchPrograms() (chan *string, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s", api.baseURL, "engagements.json"), nil)
	if err != nil {
		return nil, err
	}

	result := make(chan *string)

	go func() error {

		for page := 1; true; page++ {

			var programs programsResult

			req = api.prepareProgramsRequest(req, page)
			if err != nil {
				return err
			}

			res, err := api.client.Do(req)
			if err != nil {
				return err
			}

			if err := json.NewDecoder(res.Body).Decode(&programs); err != nil {
				return err
			}

			if len(programs.Engagements) == 0 {
				result <- nil
				break
			}

			for _, val := range programs.Engagements {
				programValid := len(api.options.Filters) == 0

				for i := 0; i < len(api.options.Filters) && !programValid; i++ {
					programValid = api.options.Filters[i].CheckName(val.Name)
				}

				if programValid {
					result <- &val.Brief
				}
			}

		}
		return nil
	}()

	return result, nil
}

type bugcrowdScopeEntry struct {
	Targets []struct {
		Uri  string `json:"uri"`
		Name string `json:"name"`
	} `json:"targets"`
	InScope bool `json:"inScope"`
}

func convertScope(entries []bugcrowdScopeEntry) pbTypes.Scope {
	var result pbTypes.Scope

	for _, entry := range entries {
		for _, target := range entry.Targets {
			if !strings.Contains(target.Uri, target.Name) {
				continue
			}
			ep := pbTypes.NewEndpointFromString(target.Name)
			if entry.InScope {
				result.AllowedEndpoints = append(result.AllowedEndpoints, ep)
			} else {
				result.DeniedEndpoints = append(result.DeniedEndpoints, ep)
			}
		}
	}

	return result
}

type bugcrowdProgram struct {
	Data struct {
		Brief struct {
			Id   string `json:"id"`
			Name string `json:"name"`
		} `json:"brief"`

		Scope []bugcrowdScopeEntry `json:"scope"`
	} `json:"data"`
}

type programEndpoints struct {
	BriefApi struct {
		BriefVersionDocument string `json:"getBriefVersionDocument"`
	} `json:"engagementBriefApi"`
}

func (api *BugcrowdApi) retrieveScopeEndpoint(endpoint string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s/changelog/", api.baseURL, endpoint), nil)
	if err != nil {
		return "", err
	}

	res, err := api.client.Do(req)
	if err != nil {
		return "", err
	}

	tokenizer := html.NewTokenizer(res.Body)
	for {
		tt := tokenizer.Next()

		switch tt {
		case html.ErrorToken:
			return "", nil

		case html.StartTagToken, html.SelfClosingTagToken:
			t := tokenizer.Token()

			if t.Data == "div" {
				for _, attr := range t.Attr {
					if attr.Key == "data-api-endpoints" {
						return attr.Val, nil
					}
				}
			}
		}
	}
}

func (api *BugcrowdApi) fetchProgram(endpoint string) (pbTypes.Program, error) {

	res, err := api.retrieveScopeEndpoint(endpoint)
	if err != nil {
		return pbTypes.Program{}, nil
	}

	var endpoints programEndpoints

	if err := json.Unmarshal([]byte(res), &endpoints); err != nil {
		return pbTypes.Program{}, err
	}

	bcProg, err := api.retrieveProgramInfo(endpoints.BriefApi.BriefVersionDocument)
	if err != nil {
		return pbTypes.Program{}, nil
	}

	scope := convertScope(bcProg.Data.Scope)

	return pbTypes.Program{
		PlatformId: bcProg.Data.Brief.Id,
		Platform:   platform,
		Name:       bcProg.Data.Brief.Name,
		Scope:      scope,
	}, nil
}

func (api *BugcrowdApi) retrieveProgramInfo(endpoint string) (bugcrowdProgram, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/%s.json", api.baseURL, endpoint), nil)
	if err != nil {
		return bugcrowdProgram{}, err
	}

	res, err := api.client.Do(req)
	if err != nil {
		return bugcrowdProgram{}, err
	}

	var result bugcrowdProgram
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return bugcrowdProgram{}, err
	}

	return result, nil
}

func (api *BugcrowdApi) FetchPrograms() (chan *pbTypes.Program, error) {

	ch, err := api.fetchPrograms()
	if err != nil {
		return nil, err
	}

	result := make(chan *pbTypes.Program)

	go func() error {
		var ep *string
		for ep = <-ch; ep != nil; ep = <-ch {
			prg, err := api.fetchProgram(*ep)
			if err != nil {
				log.Fatal(err)
				result <- nil
				return err
			}
			result <- &prg
		}
		result <- nil
		return nil
	}()

	return result, nil
}
