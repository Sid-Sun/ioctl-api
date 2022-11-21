package model

import (
	"errors"
	"fmt"
	"io"

	"github.com/sid-sun/ioctl-api/src/storageprovider"
	"github.com/sid-sun/ioctl-api/src/types"
)

var ErrNotFound = storageprovider.ErrNotFound
var ErrAlreadyExists = storageprovider.ErrAlreadyExists
var ErrIncorrectSnippetType = errors.New("invalid snipppet type/id provided")

type SnippetController interface {
	NewSnippet(snippet io.Reader, id string, st types.SnippetType) error
	FindSnippet(name string, st types.SnippetType) (*Snippet, error)
}

type Snippet struct {
	ID      string
	Snippet []byte
}

type mongoSnippetController struct {
	sp *storageprovider.S3Provider
}

func NewMongoSnippetController(sp *storageprovider.S3Provider) SnippetController {
	return &mongoSnippetController{
		sp: sp,
	}
}

func (msc *mongoSnippetController) NewSnippet(snippet io.Reader, id string, st types.SnippetType) error {
	switch st {
	case types.EphemeralSnippet:
		id = fmt.Sprintf("ephemeral/%s", id)
	case types.ProlongedSnippet:
		id = fmt.Sprintf("prolonged/%s", id)
	default: // creating static snippets is not supported
		return ErrIncorrectSnippetType
	}
	err := msc.sp.UploadSnippet(snippet, id)
	if err != nil {
		return err
	}
	return nil
}

func (msc *mongoSnippetController) FindSnippet(id string, st types.SnippetType) (*Snippet, error) {
	switch st {
	case types.EphemeralSnippet:
		id = fmt.Sprintf("ephemeral/%s", id)
	case types.StaticSnippet:
		id = fmt.Sprintf("static/%s", id)
	case types.ProlongedSnippet:
		id = fmt.Sprintf("prolonged/%s", id)
	case types.InvalidSnippet:
		return nil, ErrIncorrectSnippetType
	}
	data, err := msc.sp.DownloadSnippet(id)
	if err != nil {
		return nil, err
	}

	// Create new Snippet and return
	return &Snippet{
		ID:      id,
		Snippet: data,
	}, nil
}
