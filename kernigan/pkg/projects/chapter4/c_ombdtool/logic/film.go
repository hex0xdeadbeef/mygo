package issuetool

type FilmInfo struct {
	Title     string `json:"Title,omitempty"`
	Year      string `json:"Year,omitempty"`
	Released  string `json:"Released,omitempty"`
	Runtime   string `json:"Runtime,omitempty"`
	Genre     string `json:"Genre,omitempty"`
	PosterURL string `json:"Poster,omitempty"`
}

type Emtpy struct{}
