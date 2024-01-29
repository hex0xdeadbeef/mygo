package issuetool

const (
	SpaceCutSet = "\t\n\v\f\r \u0085\u00A0"

	// Validation things
	TokenValidationURL = "https://api.github.com/user"
	ClassicTokenPrefix = "ghp_"
	FineGrainedPrefix  = "github_pat_"
	ValidParamsCount   = 5

	/* Request portions */
	CreationLink = "https://api.github.com/repos/%s/%s/issues"
	DeletionLink = "https://api.github.com/repos/%s/%s/issues/%s"

	/* Request headers */
	Accept           = "Accept: application/vnd.github+json"
	GitHubAPIVersion = "X-GitHub-Api-Version: 2022-11-28"
	Authorization    = "Authorization: Bearer "
)
