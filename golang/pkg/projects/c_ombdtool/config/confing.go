package issuetool

const (
	SpaceCutSet = "\t\n\v\f\r \u0085\u00A0"

	// Validation things
	ValidParamsCount = 3

	/* Request portions */
	APIURL        = "https://www.omdbapi.com/?apikey=%s&"
	ValidationURL = "https://www.omdbapi.com/apikey.aspx?VERIFYKEY=%s"
)
