package issuetool

const (
	LogFilePath = "pkg/projects/issuetool"
	SpaceCutSet = "\t\n\v\f\r \u0085\u00A0"

	TokenValidationURL = "https://api.github.com/user"
	ClassicTokenPrefix = "ghp_"
	FineGrainedPrefix  = "github_pat_"

	TestToken = ""

	ValidParamsCount = 4

	/* Request portions */
	LinkPrefix = "https://api.github.com/repos/"
	LinkSuffix = "/issues/"

	// Request headers
	Accept           = "Accept: application/vnd.github+json"
	GitHubAPIVersion = "X-GitHub-Api-Version: 2022-11-28"
	Authorization    = "Authorization: Bearer "
)
