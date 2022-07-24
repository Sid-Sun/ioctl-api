package contract

import "encoding/json"

type CreateSnippet struct {
	Metadata Metadata `json:"metadata"`
	Data     string   `json:"data"`
}

type CreateE2ESnippet struct {
	Version    string          `json:"version"`
	Ephemeral  bool            `json:"ephemeral"`
	Keysalt    json.RawMessage `json:"keysalt"`
	Initvector json.RawMessage `json:"initvector"`
	Ciphertext json.RawMessage `json:"ciphertext"`
}

type CreateSnippetResponse struct {
	URL string
}

type Metadata struct {
	Ephemeral bool   `json:"ephemeral"`
	Language  string `json:"language"`
}

var CS CreateSnippet

