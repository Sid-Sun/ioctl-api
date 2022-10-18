package types

type EncryptionStack struct {
	ID   string
	Key  []byte
	Hash string
	Salt *[32]byte
}

type SnippetSpec struct {
	Version    string `json:"version"`
	Keysalt    string `json:"keysalt"`
	Ephemeral  bool   `json:"ephemeral"`
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
