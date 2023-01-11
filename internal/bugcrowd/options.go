package bugcrowd


type Options struct {
    FetchTargets bool;
    MaxPages int;
    Sort string;
    MaxRequests int;
    Hidden bool;
}

func DefaultOptions() *Options {
    return &Options {
        FetchTargets: true,
        MaxPages: -1,
        Sort: "starts-desc",
        MaxRequests: 5,
        Hidden: false,
    }
}
