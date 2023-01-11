package utils

import (
    "regexp"
)


const (
    param string = `(\?[\-\w=\.~\;\[\]&]+)?`
    localUrl = `(/[\w\.~=\-]+)+/?`
    subdomain = `([\w\-]+\.)+[a-z]{2,7}`
)

var (
    SubdomainMatcher = regexp.MustCompile(subdomain)
    URLMatcher = regexp.MustCompile(`https?://`+ subdomain + localUrl + param)
    EmailMatcher = regexp.MustCompile(`[\w\-\.]+@` + subdomain)
)
