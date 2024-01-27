package issuetool

type CreateRequestBody struct {
	Accept string

	Owner string
	Repo  string
	Title string

	Body      string
	Assignee  string
	Milestone string
	Labels    []string
	Assignees []string
}
