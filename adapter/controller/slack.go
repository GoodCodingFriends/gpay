package controller

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/entity"
	"github.com/GoodCodingFriends/gpay/repository"
	"github.com/GoodCodingFriends/gpay/usecase"
	"github.com/nlopes/slack"
)

const (
	cmdTypePay   = "pay"
	cmdTypeClaim = "claim"
)

var (
	ErrUnknownCommand = errors.New("unknown command")
	ErrInvalidUsage   = errors.New("invalid usage")
)

type SlackBot struct {
	logger *log.Logger
	cfg    *config.Config
	rtm    *slack.RTM
	repo   *repository.Repository
}

func NewSlackBot(logger *log.Logger, cfg *config.Config, repo *repository.Repository) *SlackBot {
	l := slack.New(cfg.Controller.Slack.APIToken)
	slack.SetLogger(logger)
	return &SlackBot{
		logger: logger,
		cfg:    cfg,
		rtm:    l.NewRTM(),
		repo:   repo,
	}
}

func (b *SlackBot) Listen() error {
	go b.rtm.ManageConnection()
	for m := range b.rtm.IncomingEvents {
		switch e := m.Data.(type) {
		case *slack.MessageEvent:
			if err := b.handleMessageEvent(e); err != nil {
				return err
			}
		}
	}
	return nil
}

func (b *SlackBot) Stop() error {
	return nil
}

// handleMessageEvent validates the message text and handles the event
// valid formats are like this:
//
//   @gpay send 500 @ktr
//   @gpay send @ktr 500
//   @gpay claim 500 @ktr
//   @gpay claim @ktr 500
//
func (b *SlackBot) handleMessageEvent(e *slack.MessageEvent) error {
	if !strings.HasPrefix(e.Text, b.cfg.Controller.Slack.BotName) {
		b.logger.Println("not gpay command, ignore")
		return nil
	}

	sp := strings.Split(e.Text, " ")
	if len(sp) != 4 {
		// show usage
		return ErrInvalidUsage
	}

	cmdType := sp[1]
	from := entity.UserID(e.User)

	switch cmdType {
	case cmdTypePay:
		return b.handlePayCommand(from, sp[2:])
	case cmdTypeClaim:
		to, amount, err := parsePayCommand(sp[2:])
		if err != nil {
			return err
		}
		tx, err := usecase.Claim(b.repo, &usecase.ClaimParam{
			FromID:  from,
			ToID:    to,
			Amount:  amount,
			Message: "",
		})
		if err != nil {
			return err
		}
		b.logger.Println(tx)
	default:
		return ErrUnknownCommand
	}
	return nil
}

func (b *SlackBot) handlePayCommand(from entity.UserID, sp []string) error {
	to, amount, err := parsePayCommand(sp)
	if err != nil {
		return err
	}
	tx, err := usecase.Pay(b.repo, &usecase.PayParam{
		FromID:  from,
		ToID:    to,
		Amount:  amount,
		Message: "",
	})
	if err != nil {
		return err
	}
	// TODO: handle
	b.logger.Println(tx)
	return nil
}

func parsePayCommand(sp []string) (to entity.UserID, amount entity.Amount, err error) {
	n, err := strconv.Atoi(sp[0])
	// TODO: check <@SOME_ID> or plain string
	if err == nil {
		// format: 500 @ktr
		amount = entity.Amount(n)
		to = entity.UserID(sp[1])
		return
	}

	n, err = strconv.Atoi(sp[1])
	if err == nil {
		// format: @ktr 500
		to = entity.UserID(sp[0])
		amount = entity.Amount(n)
		return
	}

	err = ErrInvalidUsage
	return
}

func parseClaimCommand(sp []string) (to entity.UserID, amount entity.Amount, err error) {
	n, err := strconv.Atoi(sp[0])
	if err == nil {
		// format: 500 @ktr
		amount = entity.Amount(n)
		to = entity.UserID(sp[1])
		return
	}

	n, err = strconv.Atoi(sp[1])
	if err == nil {
		// format: @ktr 500
		to = entity.UserID(sp[0])
		amount = entity.Amount(n)
		return
	}

	err = ErrInvalidUsage
	return
}
