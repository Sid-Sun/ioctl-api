package types

type EncryptionStack struct {
	Salt *[32]byte
	Key  []byte
	Hash string
	ID   string
}

type SnippetSpec struct {
	Ephemeral  bool   `json:"ephemeral"`
	Version    string `json:"version"`
	Keysalt    string `json:"keysalt"`
	Initvector string `json:"initvector"`
	Ciphertext string `json:"ciphertext"`
}

type Snippet struct {
	Metadata Metadata `json:"metadata"`
	Data     string   `json:"data"`
}

type Metadata struct {
	Ephemeral bool   `json:"ephemeral"`
	ID        string `json:"id"`
	Language  string `json:"language"`
}

type SnippetType int

const (
	StaticSnippet SnippetType = iota + 1
	EphemeralSnippet
	ProlongedSnippet
	InvalidSnippet
)
