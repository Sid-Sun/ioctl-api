package snippet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/sid-sun/ioctl-api/config"
	"github.com/sid-sun/ioctl-api/src/service"
	"github.com/sid-sun/ioctl-api/src/types"
	"github.com/sid-sun/ioctl-api/src/view/http/contract"
)

var regex, _ = regexp.Compile("[0-9a-f]")

func CreateE2E(svc service.Service, cfg *config.HTTPServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		snippetHexUUID := chi.URLParam(req, "snippetHexUUID")
		r := regex.FindAllStringIndex(snippetHexUUID, -1)
		if r == nil || len(snippetHexUUID) != 64 && len(r) != 64 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("hex encoded snippet uuid invalid"))
		}

		ephHeader := req.Header.Get("Ephemeral")
		eph, err := strconv.ParseBool(ephHeader)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = svc.CreateE2ESnippet(req.Body, snippetHexUUID, eph)
		if err != nil {
			if err == service.ErrAlreadyExists {
				w.WriteHeader(http.StatusConflict)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		resp := contract.CreateSnippetResponse{
			URL: fmt.Sprintf(cfg.GetBaseURL(), snippetHexUUID),
		}

		req.Header.Add("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.Encode(resp)
	}
}

func Create(svc service.Service, cfg *config.HTTPServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := req.Context().Value(contract.CS).(contract.CreateSnippet)

		snippet := types.Snippet{
			Data: data.Data,
			Metadata: types.Metadata{
				Ephemeral: data.Metadata.Ephemeral,
				Language:  data.Metadata.Language,
			},
		}

		snippetID, err := svc.CreateSnippet(snippet, data.Metadata.Ephemeral)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		resp := contract.CreateSnippetResponse{
			URL: fmt.Sprintf(cfg.GetBaseURL(), snippetID),
		}

		req.Header.Add("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.Encode(resp)
	}
}

func Get(svc service.Service, responseType string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		snippetID := chi.URLParam(req, "snippetID")

		snippet, err := svc.FetchSnippet(snippetID)
		if err != nil {
			if err == service.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(err.Error()))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if responseType == "raw" {
			w.Write([]byte(snippet.Data))
			return
		}

		req.Header.Add("Content-Type", "application/json")
		e := json.NewEncoder(w)
		e.Encode(snippet)
	}
}
