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
}

func (r *commonRequest) urlVerificationRequest() *urlVerificationRequest {
	return &urlVerificationRequest{
		challenge: r.Challenge,
	}
}

func (r *commonRequest) appMentionRequest() *appMentionRequest {
	return &appMentionRequest{
		user:           r.User,
		text:           r.Text,
		timestamp:      r.Timestamp,
		channel:        r.Channel,
		eventTimestamp: r.EventTimestamp,
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
		log.Println(r.Header)
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
		case "app_mention":
			appMentionHandler(w, r, req.appMentionRequest())
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

func appMentionHandler(w http.ResponseWriter, r *http.Request, req *appMentionRequest) {
	log.Println(req.user, req.text)
}
