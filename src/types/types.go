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
	Data     string                 `json:"data"`
	Metadata map[string]interface{} `json:"metadata"`
}
