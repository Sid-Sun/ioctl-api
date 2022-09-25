package model

import (
	"fmt"
	"io"

	"github.com/fitant/xbin-api/src/storageprovider"
)

var ErrNotFound = storageprovider.ErrNotFound
var ErrAlreadyExists = storageprovider.ErrAlreadyExists

type SnippetController interface {
	NewSnippet(snippet io.Reader, id string, ephemeral bool) error
	FindSnippet(name string, eph bool) (*Snippet, error)
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

func (msc *mongoSnippetController) NewSnippet(snippet io.Reader, id string, ephemeral bool) error {
	if ephemeral {
		id = fmt.Sprintf("ephemeral/%s", id)
	}
	err := msc.sp.UploadSnippet(snippet, id)
	if err != nil {
		return err
	}
	return nil
}

func (msc *mongoSnippetController) FindSnippet(id string, eph bool) (*Snippet, error) {
	if eph {
		id = fmt.Sprintf("ephemeral/%s", id)
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
