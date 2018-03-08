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
	"github.com/k0kubun/pp"
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
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var msg slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &msg); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if msg.Token != h.verificationToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	invoice, err := h.repo.Invoice.FindByID(context.Background(), entity.InvoiceID(msg.CallbackID))
	if err != nil {
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
			h.logger.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		pp.Println(tx)
		responseMessage(w, msg.OriginalMessage, claimAcceptedMessage)
	case actionNameReject:
		err = usecase.RejectInvoice(h.repo, &usecase.RejectInvoiceParam{
			InvoiceID: invoice.ID,
		})
		if err != nil {
			h.logger.Println(err)
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
