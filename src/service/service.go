package service

import (
	// "encoding/base64"
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

type Service interface {
	CreateSnippet(snippet types.Snippet, ephemeral bool) (string, error)
	CreateE2ESnippet(snippet io.Reader, snippetID string, eph bool) error
	FetchSnippet(id string) (*model.Snippet, error)
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

	snippet.Metadata["id"] = keys.ID
	rawSnippet, _ := json.Marshal(snippet)

	// Deflate snippet -> Encrypt snippet -> encode snippet
	compressedSnippet := utils.Defalte(rawSnippet)
	encryptedSnippet, iv, keysalt := utils.Encrypt(compressedSnippet, keys.Key, keys.Salt)
	snippetSpec := types.SnippetSpec{
		Version:    "v1",
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
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [CreateSnippet] [NewE2ESnippet]", err))
		return err
	}
	return nil
}

func (s *serviceImpl) FetchSnippet(id string) (*model.Snippet, error) {
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
	snip.Snippet = utils.Inflate(decryptedSnippet)

	// Return generated ID instead of the stored hashed ID
	snip.ID = id

	return snip, nil
}
