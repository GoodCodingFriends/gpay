package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type commonRequest struct {
	Type string `json:"type"`

	// url_verification
	Challenge string `json:"challenge"`
}

func (r *commonRequest) urlVerificationRequest() *urlVerificationRequest {
	return &urlVerificationRequest{
		challenge: r.Challenge,
	}
}

type urlVerificationRequest struct {
	challenge string
}

func Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		var req commonRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			switch err.(type) {
			case *json.UnmarshalTypeError:
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, "invalid JSON passed: %s", err)
			default:
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, "failed to decode request: %s", err)
			}
			return
		}

		switch req.Type {
		case "url_verification":
			urlVerificationHandler(w, r, req.urlVerificationRequest())
		default:
			w.WriteHeader(http.StatusNotFound)
		}
		return
	})
	return mux
}

func urlVerificationHandler(w http.ResponseWriter, r *http.Request, req *urlVerificationRequest) {
	if _, err := io.WriteString(w, req.challenge); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
