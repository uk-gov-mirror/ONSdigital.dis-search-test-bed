package algorithm

import (
	"fmt"
	"strings"
)

// SearchAlgorithm represents the type of search algorithm to be
// used in the search request.
type SearchAlgorithm string

const (
	// SearchAlgorithmBaseline represents the baseline search algorithm.
	SearchAlgorithmBaseline SearchAlgorithm = "baseline"
	// SearchAlgorithmUnweighted represents the unweighted search algorithm.
	SearchAlgorithmUnweighted SearchAlgorithm = "unweighted"
)

// AllSearchAlgorithms returns every registered search algorithm in canonical
// order. It is the single source of truth for which algorithms can be
// evaluated, so CLI defaults, help text and name validation all derive from
// it and a newly registered algorithm is picked up by each of them.
func AllSearchAlgorithms() []SearchAlgorithm {
	return []SearchAlgorithm{
		SearchAlgorithmBaseline,
		SearchAlgorithmUnweighted,
	}
}

// SearchAlgorithmNames returns the name of every registered search algorithm,
// for CLI help text and error messages.
func SearchAlgorithmNames() []string {
	algorithms := AllSearchAlgorithms()
	names := make([]string, 0, len(algorithms))
	for _, algo := range algorithms {
		names = append(names, string(algo))
	}
	return names
}

// ParseSearchAlgorithm resolves a name to its search algorithm, ignoring
// surrounding whitespace and case. An unrecognised name is an error rather
// than a silent fallback to the baseline algorithm, so a mistyped name can
// never have baseline results reported under it.
func ParseSearchAlgorithm(name string) (SearchAlgorithm, error) {
	candidate := SearchAlgorithm(strings.ToLower(strings.TrimSpace(name)))
	for _, algo := range AllSearchAlgorithms() {
		if candidate == algo {
			return algo, nil
		}
	}

	return "", fmt.Errorf("unknown algorithm %q: valid algorithms are %s",
		name, strings.Join(SearchAlgorithmNames(), ", "))
}

// ParseSearchAlgorithms resolves a list of names to their search algorithms,
// discarding duplicates and preserving canonical order. An empty list, or one
// holding only blank names, resolves to every registered algorithm.
func ParseSearchAlgorithms(names []string) ([]SearchAlgorithm, error) {
	selected := make(map[SearchAlgorithm]struct{}, len(names))
	for _, name := range names {
		if strings.TrimSpace(name) == "" {
			continue
		}

		algo, err := ParseSearchAlgorithm(name)
		if err != nil {
			return nil, err
		}
		selected[algo] = struct{}{}
	}

	if len(selected) == 0 {
		return AllSearchAlgorithms(), nil
	}

	algorithms := make([]SearchAlgorithm, 0, len(selected))
	for _, algo := range AllSearchAlgorithms() {
		if _, ok := selected[algo]; ok {
			algorithms = append(algorithms, algo)
		}
	}

	return algorithms, nil
}

// QueryName represents the name of a specific query type used
// in search requests.
type QueryName string

const (
	// QueryBaselineTerm represents the baseline term query.
	QueryBaselineTerm QueryName = "baseline-term"
	// QueryUnweightedTerm represents the unweighted term query.
	QueryUnweightedTerm QueryName = "unweighted-term"
	// QueryBrowse represents the browse query.
	QueryBrowse QueryName = "browse"
	// QueryCountContentType represents the count content type query.
	QueryCountContentType QueryName = "count-content-type"
	// QueryCountDistinct represents the count distinct items query.
	QueryCountDistinct QueryName = "count-distinct"
	// QueryCountTopicItems represents the count topic items query.
	QueryCountTopicItems QueryName = "count-topic-items"
)

// SortBy represents the sorting criteria for search results.
type SortBy string

const (
	// SortByFirstLetter represents sorting by the first letter of the title.
	SortByFirstLetter SortBy = "first_letter"
	// SortByReleaseDate represents sorting by the release date of the content.
	SortByReleaseDate SortBy = "release_date"
	// SortByReleaseDateAsc represents sorting by the release date of the content in
	// ascending order.
	SortByReleaseDateAsc SortBy = "release_date_asc"
	// SortByRelevance represents sorting by the relevance of the search results.
	SortByRelevance SortBy = "relevance"
	// SortByTitle represents sorting by the title of the content.
	SortByTitle SortBy = "title"
)

// SearchParameters holds the parameters for a search request.
type SearchParameters struct {
	Term           string
	From           int // This is the offset for pagination
	Size           int // This is the count for pagination
	Types          []string
	Index          string // This overrides the index
	SortBy         SortBy
	ReleasedAfter  Date
	ReleasedBefore Date
	Highlight      bool // Whether to highlight search terms in the results
	URIPrefix      string
	Topic          []string
	Now            string // The current timestamp - TODO: Find out the use case of this. Might be for releases.
	DatasetIDs     []string
	CDIDs          []string
	URIs           []string
}
