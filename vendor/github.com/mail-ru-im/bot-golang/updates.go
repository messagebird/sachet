package botgolang

import (
	"context"
	"fmt"
	"time"

	dura "github.com/hako/durafmt"
	"github.com/sirupsen/logrus"
)

const (
	sleepTime = time.Second * 3
)

var (
	sleepTimeStr = dura.Parse(sleepTime)
)

type Updater struct {
	logger      *logrus.Logger
	client      *Client
	lastEventID int
	PollTime    int
}

// NewMessageFromPart returns new message based on part message
func (u *Updater) NewMessageFromPayload(message EventPayload) *Message {
	return &Message{
		client:    u.client,
		ID:        message.MsgID,
		Chat:      Chat{ID: message.From.User.ID, Title: message.From.FirstName},
		Text:      message.Text,
		Timestamp: message.Timestamp,
	}
}

func (u *Updater) RunUpdatesCheck(ctx context.Context, ch chan<- Event) {
	_, err := u.GetLastEvents(0)
	if err != nil {
		u.logger.WithFields(logrus.Fields{
			"err": err,
		}).Debug("cannot make initial request to events")
	}

	for {
		select {
		case <-ctx.Done():
			close(ch)
			return
		default:
			events, err := u.GetLastEvents(u.PollTime)
			if err != nil {
				u.logger.WithFields(logrus.Fields{
					"err":            err,
					"retry interval": sleepTimeStr,
				}).Errorf("Failed to get updates, retrying in %s ...", sleepTimeStr)
				time.Sleep(sleepTime)

				continue
			}

			for _, event := range events {
				event.client = u.client
				event.Payload.client = u.client

				ch <- *event
			}
		}
	}
}

func (u *Updater) GetLastEvents(pollTime int) ([]*Event, error) {
	events, err := u.client.GetEvents(u.lastEventID, pollTime)
	if err != nil {
		u.logger.WithFields(logrus.Fields{
			"err":    err,
			"events": events,
		}).Debug("events getting error")
		return events, fmt.Errorf("cannot get events: %s", err)
	}

	count := len(events)
	if count > 0 {
		u.lastEventID = events[count-1].EventID
	}

	return events, nil
}

func NewUpdater(client *Client, pollTime int, logger *logrus.Logger) *Updater {
	if pollTime == 0 {
		pollTime = 60
	}

	return &Updater{
		client:      client,
		lastEventID: 0,
		PollTime:    pollTime,
		logger:      logger,
	}
}
