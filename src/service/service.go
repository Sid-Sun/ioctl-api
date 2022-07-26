package service

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/sid-sun/ioctl-api/config"
	"github.com/sid-sun/ioctl-api/src/model"
	"github.com/sid-sun/ioctl-api/src/types"
	"github.com/sid-sun/ioctl-api/src/utils"
)

var ErrNotFound = model.ErrNotFound
var ErrAlreadyExists = model.ErrAlreadyExists

type Service interface {
	CreateSnippet(snippet types.Snippet, ephemeral bool) (string, error)
	CreateE2ESnippet(snippet io.Reader, snippetHexUUID string, eph bool) error
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
	// Compress user Snippet instead of entire object (more effective and cleaner) then B64 encode it
	// Without B64, JSON Marshal produces a large output for non-text data due to encoding challenges (think: escapes)
	// B64 overcomes this by avoiding these and thus JSOn is of a proportional length
	snippet.Data = base64.RawURLEncoding.EncodeToString(utils.Deflate([]byte(snippet.Data)))
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

	// Srsly Golang, can we just have ternary?
	st := types.EphemeralSnippet
	if !ephemeral {
		st = types.ProlongedSnippet
	}
	err = s.sc.NewSnippet(bytes.NewReader(data), keys.Hash, st)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [CreateSnippet] [NewSnippet]", err))
		return "", err
	}

	return keys.ID, nil
}

func (s *serviceImpl) CreateE2ESnippet(snippet io.Reader, snippetHexUUID string, ephemeral bool) error {
	st := types.EphemeralSnippet
	if !ephemeral {
		st = types.ProlongedSnippet
	}
	d, err := ioutil.ReadAll(snippet)
	if err != nil {
		return err
	}
	if len(d) == 0 {
		return errors.New("body cannot be empty")
	}
	err = s.sc.NewSnippet(bytes.NewReader(d), snippetHexUUID, st)
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
	snip, err := s.sc.FindSnippet(encodedID, checkNoteType(id))
	if err != nil {
		if err != model.ErrNotFound {
			utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [FindSnippet]", err))
		}
		return nil, err
	}

	snippetSpec := new(types.SnippetSpec)
	err = json.Unmarshal(snip.Snippet, &snippetSpec)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [Unmarshal] [SnippetSpec]", err))
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
			utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [Unmarshal] [V1] [decryptedSnippet]", err))
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
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [Unmarshal] [V2] [decryptedSnippet]", err))
		return nil, err
	}

	decodedData, err := base64.RawURLEncoding.DecodeString(snippet.Data)
	if err != nil {
		utils.Logger.Error(fmt.Sprintf("%s : %v", "[Service] [FetchSnippet] [Unmarshal] [V2] [DecodeString]", err))
		return nil, err
	}

	snippet.Data = string(utils.Inflate(decodedData))

	return &snippet, nil
}
