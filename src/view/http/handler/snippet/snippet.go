package snippet

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/types"
	"github.com/fitant/xbin-api/src/view/http/contract"
	"github.com/go-chi/chi/v5"
)

func CreateE2E(svc service.Service, cfg *config.HTTPServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		snippetID := chi.URLParam(req, "snippetID")

		ephHeader := req.Header.Get("Ephemeral")
		eph, err := strconv.ParseBool(ephHeader)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = svc.CreateE2ESnippet(req.Body, snippetID, eph)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		resp := contract.CreateSnippetResponse{
			URL: fmt.Sprintf(cfg.GetBaseURL(), snippetID),
		}

		raw, _ := json.Marshal(resp)
		req.Header.Add("Content-Type", "application/json")
		w.Write(raw)
	}
}

// TODO: REDO
func Create(svc service.Service, cfg *config.HTTPServerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		data := req.Context().Value(contract.CS).(contract.CreateSnippet)

		snippet := types.Snippet{
			Data: data.Data,
			Metadata: map[string]interface{}{
				"ephemeral": data.Metadata.Ephemeral,
				"language":  data.Metadata.Language,
			},
		}
		raw, _ := json.Marshal(snippet)

		snippetID, err := svc.CreateSnippet(raw, data.Metadata.Ephemeral)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp := contract.CreateSnippetResponse{
			URL: fmt.Sprintf(cfg.GetBaseURL(), snippetID),
		}

		raw, _ = json.Marshal(resp)
		req.Header.Add("Content-Type", "application/json")
		w.Write(raw)
	}
}

// TODO: REDO
func Get(svc service.Service, responseType string) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		snippetID := chi.URLParam(req, "snippetID")

		snippet, err := svc.FetchSnippet(snippetID)
		if err != nil {
			if err == service.ErrNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		if responseType == "raw" {
			x := new(types.Snippet)
			json.Unmarshal(snippet.Snippet, &x)
			w.Write([]byte(x.Data))
			return
		}

		req.Header.Add("Content-Type", "application/json")
		w.Write(snippet.Snippet)
	}
}
