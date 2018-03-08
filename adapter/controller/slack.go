package controller

import (
	"errors"
	"log"
	"regexp"
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

	userIDPattern = regexp.MustCompile(`^<@(.*)>$`)
)

type SlackBot struct {
	logger *log.Logger
	cfg    *config.Config
	client *slack.Client
	repo   *repository.Repository

	// for testing
	disableAPIRequest bool

	// TODO: lock
	idToSlackUser map[string]slack.User
}

func NewSlackBot(logger *log.Logger, cfg *config.Config, repo *repository.Repository) (*SlackBot, error) {
	client := slack.New(cfg.Controller.Slack.APIToken)
	slack.SetLogger(logger)
	bot := &SlackBot{
		logger: logger,
		cfg:    cfg,
		client: client,
		repo:   repo,
	}
	return bot, bot.updateSlackUsers()
}

func (b *SlackBot) Listen() error {
	rtm := b.client.NewRTM()
	go rtm.ManageConnection()
	for m := range rtm.IncomingEvents {
		switch e := m.Data.(type) {
		case *slack.MessageEvent:
			if err := b.handleMessageEvent(e); err != nil {
				b.logger.Printf("handleMessageEvent: %s", err)
			}
		}
	}
	return nil
}

func (b *SlackBot) Stop() error {
	return nil
}

// not used yet
func (b *SlackBot) updateSlackUsers() error {
	b.idToSlackUser = map[string]slack.User{}
	if b.disableAPIRequest {
		return nil
	}

	users, err := b.client.GetUsers()
	if err != nil {
		return err
	}
	for _, u := range users {
		id := u.ID
		if id == "" || u.IsBot {
			continue
		}
		b.idToSlackUser[id] = u
	}
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
		return b.handleClaimCommand(from, sp[2:])
	default:
		return ErrUnknownCommand
	}
}

func (b *SlackBot) handlePayCommand(fromID entity.UserID, sp []string) error {
	p := &parser{idToSlackUser: b.idToSlackUser}
	toID, amount, err := p.parse(sp)
	if err != nil {
		return err
	}
	from, to, err := usecase.FindBothUsersWithUserCreation(b.cfg, b.repo, fromID, toID)
	tx, err := usecase.Pay(b.repo, &usecase.PayParam{
		From:    from,
		To:      to,
		Amount:  amount,
		Message: "",
	})
	if err != nil {
		return err
	}
	// TODO: handle
	b.logger.Printf("%#v\n", tx)
	return nil
}

func (b *SlackBot) handleClaimCommand(fromID entity.UserID, sp []string) error {
	p := &parser{idToSlackUser: b.idToSlackUser}
	toID, amount, err := p.parse(sp[2:])
	if err != nil {
		return err
	}
	from, to, err := usecase.FindBothUsersWithUserCreation(b.cfg, b.repo, fromID, toID)
	tx, err := usecase.Claim(b.repo, &usecase.ClaimParam{
		From:    from,
		To:      to,
		Amount:  amount,
		Message: "",
	})
	if err != nil {
		return err
	}
	b.logger.Println(tx)
	return nil
}

type parser struct {
	idToSlackUser map[string]slack.User
}

// TODO: use more better naming
func (p *parser) parse(args []string) (to entity.UserID, amount entity.Amount, err error) {
	n, err := strconv.Atoi(args[0])
	// TODO: check <@SOME_ID> or plain string
	if err == nil {
		// format: 500 @ktr
		amount = entity.Amount(n)
		to, err = p.normalizeUserID(args[1])
		if err != nil {
			return
		}
		return
	}

	n, err = strconv.Atoi(args[1])
	if err == nil {
		// format: @ktr 500
		to, err = p.normalizeUserID(args[1])
		if err != nil {
			return
		}
		amount = entity.Amount(n)
		return
	}

	err = ErrInvalidUsage
	return
}

func (p *parser) normalizeUserID(s string) (entity.UserID, error) {
	res := userIDPattern.FindStringSubmatch(s)
	if len(res) == 2 {
		if _, ok := p.idToSlackUser[res[1]]; ok {
			return entity.UserID(res[1]), nil
		}
	}
	return entity.UserID(""), errors.New("invalid form")
}
