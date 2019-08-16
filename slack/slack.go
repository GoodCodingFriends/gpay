package slack

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
)

type commonRequest struct {
	Type string `json:"type"`

	// url_verification
	Challenge string `json:"challenge"`

	// app_mention
	User           string      `json:"user"`
	Text           string      `json:"text"`
	Timestamp      json.Number `json:"ts"`
	Channel        string      `json:"channel"`
	EventTimestamp json.Number `json:"event_ts"`

	// event_callback
	Event *event `json:"event"`
}

type event struct {
	Type           string      `json:"type"`
	User           string      `json:"user"`
	Text           string      `json:"text"`
	Timestamp      json.Number `json:"ts"`
	Channel        string      `json:"channel"`
	EventTimestamp json.Number `json:"event_ts"`
}

func (r *commonRequest) urlVerificationRequest() *urlVerificationRequest {
	return &urlVerificationRequest{
		challenge: r.Challenge,
	}
}

type urlVerificationRequest struct {
	challenge string
}

type appMentionRequest struct {
	user           string
	text           string
	timestamp      json.Number
	channel        string
	eventTimestamp json.Number
}

func Router(r chi.Router) {
	r.Use(authenticationMiddleware)
	r.HandleFunc("/events", func(w http.ResponseWriter, r *http.Request) {
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
		case "event_callback":
			switch req.Event.Type {
			case "app_mention":
				appMentionHandler(w, r, req.Event)
			}
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	})
}

func urlVerificationHandler(w http.ResponseWriter, r *http.Request, req *urlVerificationRequest) {
	if _, err := io.WriteString(w, req.challenge); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func appMentionHandler(w http.ResponseWriter, r *http.Request, e *event) {
	log.Println(e.User, e.Text)
}
