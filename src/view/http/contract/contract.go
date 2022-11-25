package contract

type CreateSnippet struct {
	Metadata Metadata `json:"metadata"`
	Data     string   `json:"data"`
}

type CreateSnippetResponse struct {
	URL string
}

type Metadata struct {
	Ephemeral bool   `json:"ephemeral"`
	Language  string `json:"language"`
}

var CS CreateSnippet
