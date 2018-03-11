package controller

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
	"github.com/GoodCodingFriends/gpay/usecase"
	"github.com/nlopes/slack"
)

type interactionHandler struct {
	logger            *log.Logger
	repo              *repository.Repository
	verificationToken string
}

func (h *interactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.logger.Printf("failed to read body: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(buf) < 8 {
		h.logger.Printf("body is invalid: %s\n", string(buf))
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		h.logger.Printf("failed to unescape body: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var msg slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
		h.logger.Printf("failed to unmarshal body: %s\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if msg.Token != h.verificationToken {
		h.logger.Printf("invalid verification token: %s\n", msg.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	invoice, err := h.repo.Invoice.FindByID(context.Background(), entity.InvoiceID(msg.CallbackID))
	if err != nil {
		h.logger.Printf("specified invoice not found: %s\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch msg.Actions[0].Name {
	case actionNameAccept:
		var tx *entity.Transaction
		tx, err = usecase.AcceptInvoice(h.repo, &usecase.AcceptInvoiceParam{
			InvoiceID: invoice.ID,
		})
		if err != nil {
			h.logger.Printf("failed to accept the invoice: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		h.logger.Println(tx)
		responseMessage(w, msg.OriginalMessage, claimAcceptedMessage)
	case actionNameReject:
		err = usecase.RejectInvoice(h.repo, &usecase.RejectInvoiceParam{
			InvoiceID: invoice.ID,
		})
		if err != nil {
			h.logger.Printf("failed to reject the invoice: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		responseMessage(w, msg.OriginalMessage, claimRejectedMessage)
	}
	return
}

func responseMessage(w http.ResponseWriter, original slack.Message, title string) {
	original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Short: false,
		},
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}
