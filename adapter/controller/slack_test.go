package controller

import (
	"io/ioutil"
	"log"
	"strings"
	"testing"

	"github.com/GoodCodingFriends/gpay/config"
	"github.com/GoodCodingFriends/gpay/repository/repositorytest"
	"github.com/nlopes/slack"
	"github.com/stretchr/testify/require"
)

func setupSlackBot(t *testing.T) *SlackBot {
	logger := log.New(ioutil.Discard, "", log.LstdFlags)
	cfg, err := config.Process()
	require.NoError(t, err)

	repo := repositorytest.NewInMemory()

	return NewSlackBot(logger, cfg, repo)
}

func TestSlackBot_handleMessageEvent(t *testing.T) {
	s := setupSlackBot(t)
	e := &slack.MessageEvent{}

	t.Run("not gpay command", func(t *testing.T) {
		e.Text = "foo"
		err := s.handleMessageEvent(e)
		require.NoError(t, err)
	})

	t.Run("len(sp) is not 4", func(t *testing.T) {
		e.Text = "gpay foo"
		err := s.handleMessageEvent(e)
		require.Equal(t, ErrInvalidUsage, err)
	})

	t.Run("unknown command", func(t *testing.T) {
		e.Text = "gpay kumiko reina shuichi"
		err := s.handleMessageEvent(e)
		require.Equal(t, ErrUnknownCommand, err)
	})

	t.Run("unknown command", func(t *testing.T) {
		e.Text = "gpay kumiko reina shuichi"
		err := s.handleMessageEvent(e)
		require.Equal(t, ErrUnknownCommand, err)
	})
}

func Test_parsePayCommand(t *testing.T) {
	cases := []struct {
		in     string
		hasErr bool
	}{
		{"500 @ktr", false},
		{"@ktr 500", false},
		{"500 ktr", false},
		{"ktr 500", false},
		{"ktr 50o", true},
	}

	for _, c := range cases {
		_, _, err := parsePayCommand(strings.Split(c.in, " "))
		if c.hasErr {
			require.Error(t, err)
		}
	}
}
