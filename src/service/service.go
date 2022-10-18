package service

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/model"
	"github.com/fitant/xbin-api/src/types"
	"github.com/fitant/xbin-api/src/utils"
)

var ErrNotFound = model.ErrNotFound
var ErrAlreadyExists = model.ErrAlreadyExists

type Service interface {
	CreateSnippet(snippet types.Snippet, ephemeral bool) (string, error)
	CreateE2ESnippet(snippet io.Reader, snippetID string, eph bool) error
	FetchSnippet(id string) (*types.Snippet, error)
}

type serviceImpl struct {
	sc        model.SnippetController
	overrides map[string]string
}

func NewSnippetService(sc model.SnippetController, cfg config.Service) Service {
	encryptionKeys = make(chan types.EncryptionStack, 3)
	encryptionKeysEphemeral = make(chan types.EncryptionStack, 3)
	go populateEncryptionStack(2)
	return &serviceImpl{
		sc:        sc,
		overrides: cfg.Overrides,
	}
}

func (s *serviceImpl) CreateSnippet(snippet types.Snippet, ephemeral bool) (string, error) {
	var keys types.EncryptionStack
	switch ephemeral {
	case true:
		keys = <-encryptionKeysEphemeral
	case false:
		keys = <-encryptionKeys
	}

	snippet.Metadata.ID = keys.ID
	// Snippet Spec v2
	// Compress user Snippet isntead of entire object (more effective and cleaner) then B64 encode it
	// Without B64, JSON Marshal produces a large output for non-text data due to encoding challanges (think: escapes)
	// B64 overcomes this by avoiding these and thus JSOn is of a proportional length
	snippet.Data = base64.RawURLEncoding.EncodeToString(utils.Defalte([]byte(snippet.Data)))
	rawSnippet, err := json.Marshal(snippet)
	if err != nil {
		return "", err
	}

	// Encrypt snippet -> encode snippet
	encryptedSnippet, iv, keysalt := utils.Encrypt(rawSnippet, keys.Key, keys.Salt)
	snippetSpec := types.SnippetSpec{
		Version:    "v2",
		Ephemeral:  ephemeral,
		Ciphertext: base64.RawURLEncoding.EncodeToString(encryptedSnippet),
		Initvector: base64.RawURLEncoding.EncodeToString(iv),
		Keysalt:    base64.RawURLEncoding.EncodeToString(keysalt),
	}

	data, err := json.Marshal(snippetSpec)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [CreateSnippet] [Marshal]", err))
		return "", err
	}

	err = s.sc.NewSnippet(bytes.NewReader(data), keys.Hash, ephemeral)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [CreateSnippet] [NewSnippet]", err))
		return "", err
	}

	return keys.ID, nil
}

func (s *serviceImpl) CreateE2ESnippet(snippet io.Reader, snippetID string, eph bool) error {
	err := s.sc.NewSnippet(snippet, snippetID, eph)
	if err != nil {
		if err != model.ErrAlreadyExists {
			utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [CreateSnippet] [NewE2ESnippet]", err))
		}
		return err
	}
	return nil
}

func (s *serviceImpl) FetchSnippet(id string) (*types.Snippet, error) {
	if s.overrides[id] != "" {
		id = s.overrides[id]
	}
	hashedID := utils.HashID([]byte(id))
	encodedID := hex.EncodeToString(hashedID)
	snip, err := s.sc.FindSnippet(encodedID, checkIfEphemeral(id))
	if err != nil {
		if err != model.ErrNotFound {
			utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [FindSnippet]", err))
		}
		return nil, err
	}

	snippetSpec := new(types.SnippetSpec)
	err = json.Unmarshal(snip.Snippet, &snippetSpec)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [Unmarshal]", err))
		return nil, err
	}

	ciphertext, err := base64.RawURLEncoding.DecodeString(snippetSpec.Ciphertext)
	if err != nil {
		return nil, err
	}
	iv, err := base64.RawURLEncoding.DecodeString(snippetSpec.Initvector)
	if err != nil {
		return nil, err
	}
	salt, err := base64.RawURLEncoding.DecodeString(snippetSpec.Keysalt)
	if err != nil {
		return nil, err
	}

	decryptedSnippet := utils.Decrypt(ciphertext, salt, iv, []byte(id))
	
	var snippet types.Snippet
	if snippetSpec.Version == "v1" {
		decompressedJSON := utils.Inflate(decryptedSnippet)
		err = json.Unmarshal(decompressedJSON, &snippet)
		if err != nil {
			return nil, err
		}
		// This was not *always* set on v1 snippets, lets set it
		snippet.Metadata.ID = id
		return &snippet, nil
	}

	// In v2 we only compress/decompress user data of json instead of complete JSON
	// And B64 encode the user data post compression
	err = json.Unmarshal(decryptedSnippet, &snippet)
	if err != nil {
		return nil, err
	}

	decodedData, err := base64.RawURLEncoding.DecodeString(snippet.Data)
	if err != nil {
		return nil, err
	}

	snippet.Data = string(utils.Inflate(decodedData))

	return &snippet, nil
}
