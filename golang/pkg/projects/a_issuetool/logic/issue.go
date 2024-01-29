package issuetool

import "time"

type Issue struct {
	Title     string    `json:"title,omitempty"`
	Body      string    `json:"body,omitempty"`
	Milestone int       `json:"milestone,omitempty"`
	Labels    []string  `json:"labels,omitempty"`
	Assignees []string  `json:"assignees,omitempty"`
	HTMLURL   string    `json:"html_url,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	State     string    `json:"state,omitempty"`
	User      *User
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}

type Empty struct{}
