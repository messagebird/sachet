package botgolang

/*
💥 botgolang is zero-configuration library with convenient interface.
Crafted with love in @mail for your awesome bots.
*/

import (
	"context"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
)

const (
	defaultAPIURL = "https://api.icq.net/bot/v1"
	defaultDebug  = false
)

// Bot is the main structure for interaction with API.
// All fields are private, you can configure bot using config arguments in NewBot func.
type Bot struct {
	client  *Client
	updater *Updater
	logger  *logrus.Logger
	Info    *BotInfo
}

// GetInfo returns information about bot:
// id, name, about, avatar
func (b *Bot) GetInfo() (*BotInfo, error) {
	return b.client.GetInfo()
}

// GetChatInfo returns information about chat:
// id, type, title, public, group, inviteLink, admins
func (b *Bot) GetChatInfo(chatID string) (*Chat, error) {
	return b.client.GetChatInfo(chatID)
}

// SendChatActions sends an actions like "typing, looking"
func (b *Bot) SendChatActions(chatID string, actions ...ChatAction) error {
	return b.client.SendChatActions(chatID, actions...)
}

// GetChatAdmins returns chat admins list with fields:
// userID, creator flag
func (b *Bot) GetChatAdmins(chatID string) ([]ChatMember, error) {
	return b.client.GetChatAdmins(chatID)
}

// GetChatMem returns chat members list with fields:
// userID, creator flag, admin flag
func (b *Bot) GetChatMembers(chatID string) ([]ChatMember, error) {
	return b.client.GetChatMembers(chatID)
}

// GetChatBlockedUsers returns chat blocked users list:
// userID
func (b *Bot) GetChatBlockedUsers(chatID string) ([]User, error) {
	return b.client.GetChatBlockedUsers(chatID)
}

// GetChatPendingUsers returns chat join pending users list:
// userID
func (b *Bot) GetChatPendingUsers(chatID string) ([]User, error) {
	return b.client.GetChatPendingUsers(chatID)
}

// BlockChatUser blocks user and removes him from chat.
// If deleteLastMessages is true, the messages written recently will be deleted
func (b *Bot) BlockChatUser(chatID, userID string, deleteLastMessages bool) error {
	return b.client.BlockChatUser(chatID, userID, deleteLastMessages)
}

// UnblockChatUser unblocks user in chat
func (b *Bot) UnblockChatUser(chatID, userID string) error {
	return b.client.UnblockChatUser(chatID, userID)
}

// ResolveChatJoinRequests sends a decision to accept/decline user join to chat
func (b *Bot) ResolveChatJoinRequests(chatID, userID string, accept, everyone bool) error {
	return b.client.ResolveChatPending(chatID, userID, accept, everyone)
}

// SetChatTitle changes chat title
func (b *Bot) SetChatTitle(chatID, title string) error {
	return b.client.SetChatTitle(chatID, title)
}

// SetChatAbout changes chat about
func (b *Bot) SetChatAbout(chatID, about string) error {
	return b.client.SetChatAbout(chatID, about)
}

// SetChatRules changes chat rules
func (b *Bot) SetChatRules(chatID, rules string) error {
	return b.client.SetChatRules(chatID, rules)
}

// GetFileInfo returns information about file:
// id, type, size, filename, url
func (b *Bot) GetFileInfo(fileID string) (*File, error) {
	return b.client.GetFileInfo(fileID)
}

// NewMessage returns new message
func (b *Bot) NewMessage(chatID string) *Message {
	return &Message{
		client: b.client,
		Chat:   Chat{ID: chatID},
	}
}

// NewTextMessage returns new text message
func (b *Bot) NewTextMessage(chatID, text string) *Message {
	return &Message{
		client:      b.client,
		Chat:        Chat{ID: chatID},
		Text:        text,
		ContentType: Text,
	}
}

// NewInlineKeyboardMessage returns new text message with inline keyboard
func (b *Bot) NewInlineKeyboardMessage(chatID, text string, keyboard Keyboard) *Message {
	return &Message{
		client:         b.client,
		Chat:           Chat{ID: chatID},
		Text:           text,
		ContentType:    Text,
		InlineKeyboard: &keyboard,
	}
}

// NewFileMessage returns new file message
func (b *Bot) NewFileMessage(chatID string, file *os.File) *Message {
	return &Message{
		client:      b.client,
		Chat:        Chat{ID: chatID},
		File:        file,
		ContentType: OtherFile,
	}
}

// NewFileMessageByFileID returns new message with previously uploaded file id
func (b *Bot) NewFileMessageByFileID(chatID, fileID string) *Message {
	return &Message{
		client:      b.client,
		Chat:        Chat{ID: chatID},
		FileID:      fileID,
		ContentType: OtherFile,
	}
}

// NewVoiceMessage returns new voice message
func (b *Bot) NewVoiceMessage(chatID string, file *os.File) *Message {
	return &Message{
		client:      b.client,
		Chat:        Chat{ID: chatID},
		File:        file,
		ContentType: Voice,
	}
}

// NewVoiceMessageByFileID returns new message with previously uploaded voice file id
func (b *Bot) NewVoiceMessageByFileID(chatID, fileID string) *Message {
	return &Message{
		client:      b.client,
		Chat:        Chat{ID: chatID},
		FileID:      fileID,
		ContentType: Voice,
	}
}

// NewMessageFromPart returns new message based on part message
func (b *Bot) NewMessageFromPart(message PartMessage) *Message {
	return &Message{
		client:    b.client,
		ID:        message.MsgID,
		Chat:      Chat{ID: message.From.User.ID, Title: message.From.FirstName},
		Text:      message.Text,
		Timestamp: message.Timestamp,
	}
}

// NewButtonResponse returns new ButtonResponse
func (b *Bot) NewButtonResponse(queryID, url, text string, showAlert bool) *ButtonResponse {
	return &ButtonResponse{
		client:    b.client,
		QueryID:   queryID,
		Text:      text,
		URL:       url,
		ShowAlert: showAlert,
	}
}

func (b *Bot) NewChat(id string) *Chat {
	return &Chat{
		client: b.client,
		ID:     id,
	}
}

// SendMessage sends a message, passed as an argument.
// This method fills the argument with ID of sent message and returns an error if any.
func (b *Bot) SendMessage(message *Message) error {
	message.client = b.client
	return message.Send()
}

// EditMessage edit a message passed as an argument.
func (b *Bot) EditMessage(message *Message) error {
	return b.client.EditMessage(message)
}

// GetUpdatesChannel returns a channel, which will be filled with events.
// You can pass cancellable context there and stop receiving events.
// The channel will be closed after context cancellation.
func (b *Bot) GetUpdatesChannel(ctx context.Context) <-chan Event {
	updates := make(chan Event)

	go b.updater.RunUpdatesCheck(ctx, updates)

	return updates
}

// NewBot returns new bot object.
// All communications with bot API must go through Bot struct.
// In general you don't need to configure this bot, therefore all options are optional arguments.
func NewBot(token string, opts ...BotOption) (*Bot, error) {
	logger := logrus.New()
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	apiURL, debug := defaultAPIURL, defaultDebug
	for _, option := range opts {
		switch option.Type() {
		case "api_url":
			apiURL = option.Value().(string)
		case "debug":
			debug = option.Value().(bool)
		}
	}

	if debug {
		logger.SetLevel(logrus.DebugLevel)
	}

	client := NewClient(apiURL, token, logger)
	updater := NewUpdater(client, 0, logger)

	info, err := client.GetInfo()
	if err != nil {
		return nil, fmt.Errorf("cannot get info about bot: %s", err)
	}

	return &Bot{
		client:  client,
		updater: updater,
		logger:  logger,
		Info:    info,
	}, nil
}
