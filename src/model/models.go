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
	NewSnippet(snippet io.Reader, hexuuid string, st types.SnippetType) error
	FindSnippet(name string, st types.SnippetType) (*Snippet, error)
}

type Snippet struct {
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

func (msc *s3SnippetController) NewSnippet(snippet io.Reader, hexuuid string, st types.SnippetType) error {
	switch st {
	case types.EphemeralSnippet:
		hexuuid = fmt.Sprintf("ephemeral/%s", hexuuid)
	case types.ProlongedSnippet:
		hexuuid = fmt.Sprintf("prolonged/%s", hexuuid)
	default: // creating static snippets is not supported
		return ErrIncorrectSnippetType
	}
	err := msc.sp.UploadSnippet(snippet, hexuuid)
	if err != nil {
		return err
	}
	return nil
}

func (msc *s3SnippetController) FindSnippet(hexuuid string, st types.SnippetType) (*Snippet, error) {
	switch st {
	case types.EphemeralSnippet:
		hexuuid = fmt.Sprintf("ephemeral/%s", hexuuid)
	case types.StaticSnippet:
		hexuuid = fmt.Sprintf("static/%s", hexuuid)
	case types.ProlongedSnippet:
		hexuuid = fmt.Sprintf("prolonged/%s", hexuuid)
	case types.InvalidSnippet:
		return nil, ErrIncorrectSnippetType
	}
	data, err := msc.sp.DownloadSnippet(hexuuid)
	if err != nil {
		return nil, err
	}

	// Create new Snippet and return
	return &Snippet{
		ID:      hexuuid,
		Snippet: data,
	}, nil
}
