package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"
)

const (
	IssuesURL = "https://api.github.com/search/issues"
)

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items      []*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

func FilePrintIssues() {
	if len(os.Args) >= 3 {
		fileName := os.Args[1]
		// Open a file to write
		file, err := os.Create(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "While creating a file, error has occured: %s\n", err)
			os.Exit(1)
		}
		defer file.Close()

		terms := os.Args[2:]

		// Create a map for formatted output
		issuesMap := getIssuesMap(terms)
		sortedKeys := getSortedKeys(issuesMap)

		for _, key := range sortedKeys {
			fmt.Fprintf(file, "%s\n", key)
			for _, issue := range issuesMap[key] {
				fmt.Fprintf(file, "%-5d | %9.15s | %-20s | %-50s | %25v\n", issue.Number, issue.Title, issue.User.Login, issue.User.HTMLURL, issue.CreatedAt)
			}
			fmt.Fprintln(file)
		}

		fmt.Fprintf(os.Stdout, "Issues have written!\n")
		return
	}

	fmt.Fprintf(os.Stderr, "There are no sufficient arguments\n")
	return
}

func getIssuesMap(terms []string) map[string][]Issue {

	// Get some written issues
	issuesSearchResult, err := SearchIssues(terms)
	if err != nil {
		fmt.Fprintf(os.Stderr, "There is no terms\n")
		os.Exit(1)
	}

	ageCategories := map[string][]Issue{"0. Less than month": nil, "1. Less than year:": nil, "2. Out of year": nil}

	for _, issue := range issuesSearchResult.Items {
		difference := getMonthDifference(issue.CreatedAt)
		ageCategories[difference] = append(ageCategories[difference], *issue)
	}

	return ageCategories
}

func SearchIssues(terms []string) (*IssuesSearchResult, error) {
	query := url.QueryEscape(strings.Join(terms, " "))

	response, err := http.Get(IssuesURL + "?q=" + query)
	if err != nil {
		return nil, err
	}

	if response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("Search query failed: %s", response.Status)
	}

	var result IssuesSearchResult
	// .NewDecoder(response.Body) defines a source the info will be read from
	// .Decode(&result) defines a destiny a data will be sent to
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		response.Body.Close()
		return nil, err
	}
	response.Body.Close()
	return &result, nil
}

func getSortedKeys(issueMap map[string][]Issue) []string {
	sortedKeys := make([]string, 0, len(issueMap))

	for key := range issueMap {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	return sortedKeys
}

func getMonthDifference(createdAt time.Time) string {
	const (
		monthsInYear = 12
	)

	curYear, curMonth, _ := time.Now().Date()
	createdYear, createdMonth, _ := createdAt.Date()

	differenthInMonth := (curYear-createdYear)*monthsInYear + int(curMonth) - int(createdMonth)
	if differenthInMonth <= 0 {
		return "0. Less than month"
	} else if differenthInMonth < 12 {
		return "1. Less than year:"
	} else {
		return "2. Out of year"
	}
}

/* 		fmt.Fprintf(file, "%d Issues found:\n", issues.TotalCount)
for _, issue := range issues.Items {
	fmt.Fprintf(file, "%-5d | %9.15s | %.55s | %s | %50v\n", issue.Number, issue.Title, issue.User.Login, issue.User.HTMLURL, issue.CreatedAt)
} */
