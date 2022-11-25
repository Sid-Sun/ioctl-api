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
var errIncorrectSnippetType = errors.New("invalid snipppet type/id provided")

type SnippetController interface {
	NewSnippet(snippet io.Reader, id string, st types.SnippetType) error
	FindSnippet(name string, st types.SnippetType) (*snippet, error)
}

type snippet struct {
	ID      string
	Snippet []byte
}

type s3SnippetController struct {
	sp *storageprovider.S3Provider
}

func NewS3SnippetController(sp *storageprovider.S3Provider) SnippetController {
	return &s3SnippetController{
		sp: sp,
	}
}

func (msc *s3SnippetController) NewSnippet(snippet io.Reader, id string, st types.SnippetType) error {
	switch st {
	case types.EphemeralSnippet:
		id = fmt.Sprintf("ephemeral/%s", id)
	case types.ProlongedSnippet:
		id = fmt.Sprintf("prolonged/%s", id)
	default: // creating static snippets is not supported
		return errIncorrectSnippetType
	}
	err := msc.sp.UploadSnippet(snippet, id)
	if err != nil {
		return err
	}
	return nil
}

func (msc *s3SnippetController) FindSnippet(id string, st types.SnippetType) (*snippet, error) {
	switch st {
	case types.EphemeralSnippet:
		id = fmt.Sprintf("ephemeral/%s", id)
	case types.StaticSnippet:
		id = fmt.Sprintf("static/%s", id)
	case types.ProlongedSnippet:
		id = fmt.Sprintf("prolonged/%s", id)
	case types.InvalidSnippet:
		return nil, errIncorrectSnippetType
	}
	data, err := msc.sp.DownloadSnippet(id)
	if err != nil {
		return nil, err
	}

	// Create new Snippet and return
	return &snippet{
		ID:      id,
		Snippet: data,
	}, nil
}
