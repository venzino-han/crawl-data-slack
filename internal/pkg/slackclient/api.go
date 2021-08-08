package slackclient

import (
	"github.com/Buzzvil/crawl-data-slack/internal/pkg/crawler"
	"github.com/slack-go/slack"
	"go.uber.org/zap"
)

const CHANNEL_ID = "G015JFASK7Z"

type Client struct {
	logger *zap.Logger
	api    *slack.Client
	mapper mapper
}

func (c Client) Notify(n crawler.Notification) error {
	_, _, err := c.api.PostMessage(n.User.ID, slack.MsgOptionText(n.Event.Message, false))
	return err
}

func (c Client) GetUsers() ([]crawler.User, error) {
	users, err := c.api.GetUsers()
	if err != nil {
		return nil, err
	}

	var activeUsers []slack.User
	for _, u := range users {
		if !u.Deleted && !u.IsBot && !u.IsRestricted {
			activeUsers = append(activeUsers, u)
		}
	}

	return c.mapper.mapSlackUsersToUsers(activeUsers), nil
}

func NewClient(logger *zap.Logger, client *slack.Client) *Client {
	return &Client{
		logger: logger,
		api:    client,
	}
}