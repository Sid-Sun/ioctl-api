package middleware

import (
	"context"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/fitant/xbin-api/src/view/http/contract"
)

func WithIngestion() func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ch := req.Header.Get("Content-Type")
			data := contract.CreateSnippet{
				Metadata: contract.Metadata{
					Ephemeral: true,
				},
			}

			var source io.Reader
			source = req.Body

			raw, _ := ioutil.ReadAll(source)

			if len(raw) == 0 {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			if ch == "application/json" {
				err := json.Unmarshal(raw, &data)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if data.Data == "" && data.Metadata.Language == "" {
					data.Data = string(raw)
					data.Metadata.Language = "application/json"
				}
			} else {
				data.Data = string(raw)
				data.Metadata.Language = "plaintext"
			}

			if data.Data == "" || data.Metadata.Language == "" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), contract.CS, data)))
		})
	}
}
