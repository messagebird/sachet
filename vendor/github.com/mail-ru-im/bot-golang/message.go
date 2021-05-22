package botgolang

import (
	"fmt"
	"os"
	"path/filepath"
)

//go:generate easyjson -all message.go

type MessageContentType uint8

const (
	Unknown MessageContentType = iota
	Text
	OtherFile
	Voice
)

// Message represents a text message
type Message struct {
	client      *Client
	ContentType MessageContentType

	// Id of the message (for editing)
	ID string `json:"msgId"`

	// File contains file attachment of the message
	File *os.File `json:"-"`

	// Id of file to send
	FileID string `json:"fileId"`

	// Text of the message or caption for file
	Text string `json:"text"`

	// Chat where to send the message
	Chat Chat `json:"chat"`

	// Id of replied message
	// You can't use it with ForwardMsgID or ForwardChatID
	ReplyMsgID string `json:"replyMsgId"`

	// Id of forwarded message
	// You can't use it with ReplyMsgID
	ForwardMsgID string `json:"forwardMsgId"`

	// Id of a chat from which you forward the message
	// You can't use it with ReplyMsgID
	// You should use it with ForwardMsgID
	ForwardChatID string `json:"forwardChatId"`

	Timestamp int `json:"timestamp"`

	// The markup for the inline keyboard
	InlineKeyboard *Keyboard `json:"inlineKeyboardMarkup"`
}

func (m *Message) AttachNewFile(file *os.File) {
	m.File = file
	m.ContentType = OtherFile
}

func (m *Message) AttachExistingFile(fileID string) {
	m.FileID = fileID
	m.ContentType = OtherFile
}

func (m *Message) AttachNewVoice(file *os.File) {
	m.File = file
	m.ContentType = Voice
}

func (m *Message) AttachExistingVoice(fileID string) {
	m.FileID = fileID
	m.ContentType = Voice
}

// AttachInlineKeyboard adds a keyboard to the message.
// Note - at least one row should be in the keyboard
// and there should be no empty rows
func (m *Message) AttachInlineKeyboard(keyboard Keyboard) {
	m.InlineKeyboard = &keyboard
}

// Send method sends your message.
// Make sure you have Text or FileID in your message.
func (m *Message) Send() error {
	if m.client == nil {
		return fmt.Errorf("client is not inited, create message with constructor NewMessage, NewTextMessage, etc")
	}

	if m.Chat.ID == "" {
		return fmt.Errorf("message should have chat id")
	}

	switch m.ContentType {
	case Voice:
		if m.FileID != "" {
			return m.client.SendVoiceMessage(m)
		}

		if m.File != nil {
			return m.client.UploadVoice(m)
		}
	case OtherFile:
		if m.FileID != "" {
			return m.client.SendFileMessage(m)
		}

		if m.File != nil {
			return m.client.UploadFile(m)
		}
	case Text:
		return m.client.SendTextMessage(m)
	case Unknown:
		// need to autodetect
		if m.FileID != "" {
			// voice message's fileID always starts with 'I'
			if m.FileID[0] == voiceMessageLeadingRune {
				return m.client.SendVoiceMessage(m)
			}
			return m.client.SendFileMessage(m)
		}

		if m.File != nil {
			if voiceMessageSupportedExtensions[filepath.Ext(m.File.Name())] {
				return m.client.UploadVoice(m)
			}
			return m.client.UploadFile(m)
		}

		if m.Text != "" {
			return m.client.SendTextMessage(m)
		}
	}

	return fmt.Errorf("cannot send message or file without data")
}

// Edit method edits your message.
// Make sure you have ID in your message.
func (m *Message) Edit() error {
	if m.ID == "" {
		return fmt.Errorf("cannot edit message without id")
	}
	return m.client.EditMessage(m)
}

// Delete method deletes your message.
// Make sure you have ID in your message.
func (m *Message) Delete() error {
	if m.ID == "" {
		return fmt.Errorf("cannot delete message without id")
	}

	return m.client.DeleteMessage(m)
}

// Reply method replies to the message.
// Make sure you have ID in the message.
func (m *Message) Reply(text string) error {
	if m.ID == "" {
		return fmt.Errorf("cannot reply to message without id")
	}

	m.ReplyMsgID = m.ID
	m.Text = text

	return m.client.SendTextMessage(m)
}

// Forward method forwards your message to chat.
// Make sure you have ID in your message.
func (m *Message) Forward(chatID string) error {
	if m.ID == "" {
		return fmt.Errorf("cannot forward message without id")
	}

	m.ForwardChatID = m.Chat.ID
	m.ForwardMsgID = m.ID
	m.Chat.ID = chatID

	return m.client.SendTextMessage(m)
}

// Pin message in chat
// Make sure you are admin in this chat
func (m *Message) Pin() error {
	if m.ID == "" {
		return fmt.Errorf("cannot pin message without id")
	}

	return m.client.PinMessage(m)
}

// Unpin message in chat
// Make sure you are admin in this chat
func (m *Message) Unpin() error {
	if m.ID == "" {
		return fmt.Errorf("cannot unpin message without id")
	}

	return m.client.UnpinMessage(m)
}
