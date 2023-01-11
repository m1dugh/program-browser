package bugcrowd

import (
    "fmt"
    "encoding/json"
    "net/http"
    "errors"
    "io"
    "sync"
    "strconv"
    "strings"
    "github.com/m1dugh/program-browser/pkg/types"
)

const ROOT_URL = `https://bugcrowd.com`

const programsUrl = `%s/programs.json?sort[]=%s&hidden[]=%s&page[]=%d`

type BugcrowdRequester struct {
    throttler *types.RequestThrottler
}

func NewBugcrowdRequester(maxRequests uint) *BugcrowdRequester {
    return &BugcrowdRequester{types.NewRequestThrottler(maxRequests)}
}

func GetUrl(sort string, hidden bool, page uint) string {
    bflag := "false"
    if hidden {
        bflag = "true"
    }
    return fmt.Sprintf(programsUrl, ROOT_URL, sort, bflag, page)
}

type BProgram struct {
    Code string `json:"code"`
    Status string `json:"invitated_status"`
    InProgress bool `json:"in_progress"`
    Logo string `json:"logo"`
    Managed bool `json:"managed?"`
    Name string `json:"name"`
    Participation string `json:"participation"`
    Url string `json:"program_url"`
    SafeHarbor string `json:"safe_harbor_status"`
    Started bool `json:"started?"`
    Ended bool `json:"ended"`
    Category string `json:"industry_name"`
    MinRewards string `json:"min_rewards"`
    MaxRewards string `json:"max_rewards"`
}


func (p *BProgram) ToProgram() *types.Program {
    res := &types.Program{}
    res.Id = fmt.Sprintf("%s-%s", "bugcrowd", p.Code)
    res.Platform = "bugcrowd"
    res.PlatformUrl = p.Url
    res.Name = p.Code
    res.Logo = p.Logo
    res.Started = p.Started
    res.Ended = p.Ended
    res.Category = p.Category
    res.Managed = p.Managed
    min, err := strconv.ParseFloat(p.MinRewards, 32)
    if err == nil {
        res.MinReward = uint(min)
    }
    max, err := strconv.ParseFloat(p.MaxRewards, 32)
    if err == nil {
        res.MaxReward = uint(max)
    }

    var harborStatus string
    if p.SafeHarbor == "full" {
        harborStatus = types.Safe
    } else if p.SafeHarbor == "partial" {
        harborStatus = types.Partial
    } else {
        harborStatus = types.Unsafe
    }
    var status string
    if p.Participation == "public" {
        status = types.Public
    } else {
        status = types.Private
    }
    res.Status = status

    res.SafeHarborStatus = harborStatus

    return res
}

type BResults struct {
    Programs []*BProgram `json:"programs"`
    Meta     struct{
        TotalHits uint `json:"totalHits"`
        Pages uint `json:"totalPages"`
        PageId uint `json:"currentPage"`
    } `json:"meta"`
}

func (requester *BugcrowdRequester) getPrograms(url string) (*BResults, error) {
    requester.throttler.AskRequest()
    res, err := http.Get(url)
    
    if err != nil {
        return nil, errors.New("could not request page")
    }
    
    defer res.Body.Close()

    content, err := io.ReadAll(res.Body)
    if err != nil {
        return nil, errors.New("could not read body")
    }

    results := &BResults{}

    json.Unmarshal(content, results)

    return results, nil
}

func (requester *BugcrowdRequester) getProgramsWorker(url string,
programs *[]*BProgram,
mut *sync.Mutex,
wg *sync.WaitGroup) {
    defer wg.Done()

    results, err := requester.getPrograms(url)
    if err != nil {
        return
    }
    mut.Lock()
    *programs = append(*programs, results.Programs...)
    mut.Unlock()
}

func (requester *BugcrowdRequester) GetProgramList(sort string, hidden bool, pages int) ([]*BProgram, error) {
    if pages == 0 {
        return nil, nil
    }
    var page uint = 1
    var programs []*BProgram
    url := GetUrl(sort, hidden ,page)
    results, err := requester.getPrograms(url)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("could not get program for %s", url))
    }

    pageCount := results.Meta.Pages
    if pages > 0 && uint(pages) < pageCount {
        pageCount = uint(pages)
    }
    programs = make([]*BProgram, 0, results.Meta.TotalHits)
    programs = append(programs, results.Programs...)

    var mut sync.Mutex
    var wg sync.WaitGroup

    for page++;page <= pageCount;page++ {
        url = GetUrl(sort, hidden, page)
        wg.Add(1)
        go requester.getProgramsWorker(url, &programs, &mut, &wg)
    }

    wg.Wait()

    return programs,nil
}

func (requester *BugcrowdRequester) GetPartialPrograms(sort string, hidden bool, page int) ([]*types.Program, error) {
    bprogs, err := requester.GetProgramList(sort, hidden, page)
    if err != nil {
        return nil, errors.New("GetPartialPrograms: An error occured while getting program list")
    }

    var res []*types.Program = make([]*types.Program, len(bprogs))

    for i, v := range bprogs {
        res[i] = v.ToProgram()
    }

    return res, nil
}


const base_target_url = `/target_groups`


type BTargetGroup struct {
    InScope bool `json:"in_scope"`
    Url string `json:"targets_url"`
}

type BTargetGroups struct {
    Groups []*BTargetGroup `json:"groups"`
}

func _escapeJSON(raw json.RawMessage) (json.RawMessage, error) {
    str, err := strconv.Unquote(strings.Replace(strconv.Quote(string(raw)), `\\u`, `\u`, -1))
    if err != nil {
        return nil, err
    }
    return []byte(str), nil
}

func (requester *BugcrowdRequester) GetTargetGroups(base_url string) ([]*BTargetGroup, error) {

    client := &http.Client{}
    url := fmt.Sprintf("%s%s%s", ROOT_URL, base_url, base_target_url)
    requester.throttler.AskRequest()
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil,errors.New("could not create request")
    }

    req.Header = http.Header{
        "Accept": {"application/json"},
    }

    res, err := client.Do(req)

    if err != nil {
        return nil, errors.New("could not request url")
    }

    defer res.Body.Close()
    content, err := io.ReadAll(res.Body)

    if err != nil {
        return nil, errors.New("could read body")
    }

    var groups BTargetGroups

    err = json.Unmarshal(content, &groups)

    if err != nil {
        return nil, err
    }

    return groups.Groups, nil
}

type _target_tags struct {
    Tags    []struct {
        Name string `json:"name"`
    } `json:"tags"`
}

type BTarget struct {
    Name string `json:"name"`
    Uri  string `json:"uri"`
    Category string `json:"category"`
    target  _target_tags `json:"target"`
}

func (t *BTarget) ToTarget(inScope bool) *types.Target {
    res := &types.Target{}
    res.InScope = inScope
    var uris []string = make([]string, 1)
    if len(t.Uri) == 0 {
        uris[0] = t.Name
    } else {
        uris[0] = t.Uri
    }
    res.URIs = uris
    res.Category = t.Category
    res.Tags = make([]string, len(t.target.Tags))
    for i, v := range t.target.Tags {
        res.Tags[i] = v.Name
    }
    return res
}

type bTargets struct {
    Targets []*BTarget `json:"targets"`
}

func (requester *BugcrowdRequester) GetTarget(url string) ([]*BTarget, error) {
    
    client := &http.Client{}
    url = fmt.Sprintf("%s%s", ROOT_URL, url)
    requester.throttler.AskRequest()
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        return nil,errors.New("could not create request")
    }

    req.Header = http.Header{
        "Accept": {"application/json"},
    }

    res, err := client.Do(req)

    if err != nil {
        return nil, errors.New("could not request url")
    }

    defer res.Body.Close()
    content, err := io.ReadAll(res.Body)

    if err != nil {
        return nil, errors.New("could read body")
    }

    var targets bTargets
    err = json.Unmarshal(content, &targets)

    if err != nil {
        return nil, err
    }

    return targets.Targets, nil
}

func (requester *BugcrowdRequester) FetchTargets(p *types.Program) error {
    groups, err := requester.GetTargetGroups(p.PlatformUrl)   
    if err != nil {
        return err
    }
    
    var targets []types.Target

    for _, g := range groups {
        values, err := requester.GetTarget(g.Url)
        if err != nil {
            continue
        }

        for _, t := range values {
            targets = append(targets, *t.ToTarget(g.InScope))
        }
    }

    p.Targets = targets

    return nil
}

func (requester *BugcrowdRequester) _fetchTargetsWorker(wg *sync.WaitGroup,
p *types.Program) {
    requester.FetchTargets(p)
    wg.Done()
}

func (requester *BugcrowdRequester) FetchAllTargets(programs []*types.Program) {
    var wg sync.WaitGroup

    for _, p := range programs {
        wg.Add(1)
        go requester._fetchTargetsWorker(&wg, p)
    }

    wg.Wait()
}

func (requester *BugcrowdRequester) GetPrograms() ([]*types.Program, error) {
    results, err := requester.GetPartialPrograms("starts-desc", false, -1)
    if err != nil {
        return nil, errors.New("BugcrowdRequested.GetPrograms: an error occured while fetching partial programs")
    } else if len(results) == 0 {
        return nil, errors.New("could not request programs")
    }

    requester.FetchAllTargets(results)

    return results, nil
}

