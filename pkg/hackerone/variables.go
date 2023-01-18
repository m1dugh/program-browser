package hackerone

const DIRECTORY_QUERY = `
query DirectoryQuery($first: Int, $cursor: String, $secureOrderBy: FiltersTeamFilterOrder, $where: FiltersTeamFilterInput) {
    teams(first: $first, after: $cursor, secure_order_by: $secureOrderBy, where: $where) {
        pageInfo {
            endCursor
            hasNextPage
            __typename
        }
        edges {
            node {
                id
                bookmarked
                ...TeamTableResolvedReports
                ...TeamTableAvatarAndTitle
                ...TeamTableLaunchDate
                ...TeamTableMinimumBounty
                ...TeamTableAverageBounty
                ...BookmarkTeam
                __typename
            }
            __typename
        }
        __typename
    }
}

fragment TeamTableResolvedReports on Team {
    id
    resolved_report_count
    __typename
}

fragment TeamTableAvatarAndTitle on Team {
    id
    profile_picture(size: medium)
    name
    handle
    submission_state
    triage_active
    publicly_visible_retesting
    state
    allows_bounty_splitting
    external_program {
        id
        __typename
    }
    ...TeamLinkWithMiniProfile
    __typename
}

fragment TeamLinkWithMiniProfile on Team {
    id
    handle
    name
    __typename
}

fragment TeamTableLaunchDate on Team {
    id
    launched_at
    __typename
}

fragment TeamTableMinimumBounty on Team {
    id
    currency
    base_bounty
    __typename
}

fragment TeamTableAverageBounty on Team {
    id
    currency
    average_bounty_lower_amount
    average_bounty_upper_amount
    __typename
}

fragment BookmarkTeam on Team {
    id
    bookmarked
    __typename
}
`

type directoryRequestVariables struct {
    First int `json:"first"`      
    Cursor int `json:"cursor"`
}
