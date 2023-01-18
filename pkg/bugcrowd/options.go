package bugcrowd


type Options struct {
    FetchTargets bool;
    MaxPrograms int;
    Sort string;
    MaxRequests int;
    Hidden bool;
}

func DefaultOptions() *Options {
    return &Options {
        FetchTargets: true,
        MaxPrograms: -1,
        Sort: "starts-desc",
        MaxRequests: 5,
        Hidden: false,
    }
}
