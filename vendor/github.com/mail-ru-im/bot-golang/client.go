package botgolang

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

type Client struct {
	client  *http.Client
	token   string
	baseURL string
	logger  *logrus.Logger
}

func (c *Client) Do(path string, params url.Values, file *os.File) ([]byte, error) {
	apiURL, err := url.Parse(c.baseURL + path)
	params.Set("token", c.token)

	if err != nil {
		return nil, fmt.Errorf("cannot parse url: %s", err)
	}

	apiURL.RawQuery = params.Encode()
	req, err := http.NewRequest(http.MethodGet, apiURL.String(), nil)
	if err != nil || req == nil {
		return nil, fmt.Errorf("cannot init http request: %s", err)
	}

	if file != nil {
		buffer := &bytes.Buffer{}
		multipartWriter := multipart.NewWriter(buffer)

		fileWriter, err := multipartWriter.CreateFormFile("file", file.Name())
		if err != nil {
			return nil, fmt.Errorf("cannot create multipart writer: %s", err)
		}

		_, err = io.Copy(fileWriter, file)
		if err != nil {
			return nil, fmt.Errorf("cannot copy file into buffer: %s", err)
		}

		if err := multipartWriter.Close(); err != nil {
			return nil, fmt.Errorf("cannot close multipartWriter: %s", err)
		}

		req.Header.Set("Content-Type", multipartWriter.FormDataContentType())
		req.Body = ioutil.NopCloser(buffer)
		req.Method = http.MethodPost
	}

	c.logger.WithFields(logrus.Fields{
		"api_url": apiURL,
	}).Debug("requesting api")

	resp, err := c.client.Do(req)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("request error")
		return []byte{}, fmt.Errorf("cannot make request to bot api: %s", err)
	}

	defer func() {
		if err := resp.Body.Close(); err != nil {
			c.logger.WithFields(logrus.Fields{
				"err": err,
			}).Error("cannot close body")
		}
	}()

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		c.logger.WithFields(logrus.Fields{
			"err": err,
		}).Error("cannot read body")
		return []byte{}, fmt.Errorf("cannot read body: %s", err)
	}

	c.logger.WithFields(logrus.Fields{
		"response": string(responseBody),
	}).Debug("got response from API")

	response := &Response{}

	if err := json.Unmarshal(responseBody, response); err != nil {
		return nil, fmt.Errorf("cannot unmarshal json: %s", err)
	}

	if !response.OK {
		return responseBody, fmt.Errorf("error status from API: %s", response.Description)
	}

	return responseBody, nil
}

func (c *Client) GetInfo() (*BotInfo, error) {
	response, err := c.Do("/self/get", url.Values{}, nil)
	if err != nil {
		return nil, fmt.Errorf("error while receiving information: %s", err)
	}

	info := &BotInfo{}
	if err := json.Unmarshal(response, info); err != nil {
		return nil, fmt.Errorf("error while unmarshalling information: %s", err)
	}

	return info, nil
}

func (c *Client) GetChatInfo(chatID string) (*Chat, error) {
	params := url.Values{
		"chatId": {chatID},
	}
	response, err := c.Do("/chats/getInfo", params, nil)
	if err != nil {
		return nil, fmt.Errorf("error while receiving information: %s", err)
	}

	chat := &Chat{
		client: c,
		ID:     chatID,
	}
	if err := json.Unmarshal(response, chat); err != nil {
		return nil, fmt.Errorf("error while unmarshalling information: %s", err)
	}

	if chat.Type == Private {
		return chat, nil
	}
	return chat, nil
}

func (c *Client) SendChatActions(chatID string, actions ...ChatAction) error {
	actionsMap := make(map[ChatAction]bool)
	filteredActions := make([]ChatAction, 0)
	for _, action := range actions {
		if _, has := actionsMap[action]; !has {
			filteredActions = append(filteredActions, action)
			actionsMap[action] = true
		}
	}
	params := url.Values{
		"chatId":  {chatID},
		"actions": filteredActions,
	}
	_, err := c.Do("/chats/sendActions", params, nil)
	if err != nil {
		return fmt.Errorf("error while receiving information: %s", err)
	}
	return nil
}

func (c *Client) GetChatAdmins(chatID string) ([]ChatMember, error) {
	params := url.Values{
		"chatId": {chatID},
	}

	response, err := c.Do("/chats/getAdmins", params, nil)
	if err != nil {
		return nil, fmt.Errorf("error while receiving admins: %s", err)
	}

	admins := new(AdminsListResponse)
	if err := json.Unmarshal(response, admins); err != nil {
		return nil, fmt.Errorf("error while unmarshalling admins: %s", err)
	}
	return admins.List, nil
}

func (c *Client) GetChatMembers(chatID string) ([]ChatMember, error) {
	params := url.Values{
		"chatId": {chatID},
	}

	response, err := c.Do("/chats/getMembers", params, nil)
	if err != nil {
		return nil, fmt.Errorf("error while receiving members: %s", err)
	}

	members := new(MembersListResponse)
	if err := json.Unmarshal(response, members); err != nil {
		return nil, fmt.Errorf("error while unmarshalling members: %s", err)
	}
	return members.List, nil
}

func (c *Client) GetChatBlockedUsers(chatID string) ([]User, error) {
	params := url.Values{
		"chatId": {chatID},
	}

	response, err := c.Do("/chats/getBlockedUsers", params, nil)
	if err != nil {
		return nil, fmt.Errorf("error while receiving blocked users: %s", err)
	}

	users := new(UsersListResponse)
	if err := json.Unmarshal(response, users); err != nil {
		return nil, fmt.Errorf("error while unmarshalling blocked users: %s", err)
	}
	return users.List, nil
}

func (c *Client) GetChatPendingUsers(chatID string) ([]User, error) {
	params := url.Values{
		"chatId": {chatID},
	}

	response, err := c.Do("/chats/getPendingUsers", params, nil)
	if err != nil {
		return nil, fmt.Errorf("error while receiving pending users: %s", err)
	}

	users := new(UsersListResponse)
	if err := json.Unmarshal(response, users); err != nil {
		return nil, fmt.Errorf("error while unmarshalling pending users: %s", err)
	}
	return users.List, nil
}

func (c *Client) BlockChatUser(chatID, userID string, deleteLastMessages bool) error {
	params := url.Values{
		"chatId":          {chatID},
		"userId":          {userID},
		"delLastMessages": {strconv.FormatBool(deleteLastMessages)},
	}

	response, err := c.Do("/chats/blockUser", params, nil)
	if err != nil {
		return fmt.Errorf("error while blocking user: %s", err)
	}

	users := new(UsersListResponse)
	if err := json.Unmarshal(response, users); err != nil {
		return fmt.Errorf("error while blocking user: %s", err)
	}
	return nil
}

func (c *Client) UnblockChatUser(chatID, userID string) error {
	params := url.Values{
		"chatId": {chatID},
		"userId": {userID},
	}

	response, err := c.Do("/chats/unblockUser", params, nil)
	if err != nil {
		return fmt.Errorf("error while unblocking user: %s", err)
	}

	users := new(UsersListResponse)
	if err := json.Unmarshal(response, users); err != nil {
		return fmt.Errorf("error while unblocking user: %s", err)
	}
	return nil
}

func (c *Client) ResolveChatPending(chatID, userID string, approve, everyone bool) error {
	params := url.Values{
		"chatId":  {chatID},
		"approve": {strconv.FormatBool(approve)},
	}
	if everyone {
		params.Set("everyone", "true")
	} else {
		params.Set("userId", userID)
	}

	if _, err := c.Do("/chats/resolvePending", params, nil); err != nil {
		return fmt.Errorf("error while resolving chat pendings: %s", err)
	}
	return nil
}

func (c *Client) SetChatTitle(chatID, title string) error {
	params := url.Values{
		"chatId": {chatID},
		"title":  {title},
	}

	if _, err := c.Do("/chats/setTitle", params, nil); err != nil {
		return fmt.Errorf("error while setting chat title: %s", err)
	}
	return nil
}

func (c *Client) SetChatAbout(chatID, about string) error {
	params := url.Values{
		"chatId": {chatID},
		"about":  {about},
	}

	if _, err := c.Do("/chats/setAbout", params, nil); err != nil {
		return fmt.Errorf("error while setting chat about: %s", err)
	}
	return nil
}

func (c *Client) SetChatRules(chatID, rules string) error {
	params := url.Values{
		"chatId": {chatID},
		"rules":  {rules},
	}

	if _, err := c.Do("/chats/setRules", params, nil); err != nil {
		return fmt.Errorf("error while setting chat rules: %s", err)
	}
	return nil
}

func (c *Client) GetFileInfo(fileID string) (*File, error) {
	params := url.Values{
		"fileId": {fileID},
	}
	response, err := c.Do("/files/getInfo", params, nil)
	if err != nil {
		return nil, fmt.Errorf("error while receiving information: %s", err)
	}

	file := &File{}
	if err := json.Unmarshal(response, file); err != nil {
		return nil, fmt.Errorf("error while unmarshalling information: %s", err)
	}

	return file, nil
}

func (c *Client) GetVoiceInfo(fileID string) (*File, error) {
	return c.GetFileInfo(fileID)
}

func (c *Client) SendTextMessage(message *Message) error {
	params := url.Values{
		"chatId": {message.Chat.ID},
		"text":   {message.Text},
	}

	if message.ReplyMsgID != "" {
		params.Set("replyMsgId", message.ReplyMsgID)
	}

	if message.ForwardMsgID != "" {
		params.Set("forwardMsgId", message.ForwardMsgID)
		params.Set("forwardChatId", message.ForwardChatID)
	}

	if message.InlineKeyboard != nil {
		data, err := json.Marshal(message.InlineKeyboard.GetKeyboard())
		if err != nil {
			return fmt.Errorf("cannot marshal inline keyboard markup: %s", err)
		}

		params.Set("inlineKeyboardMarkup", string(data))
	}

	response, err := c.Do("/messages/sendText", params, nil)
	if err != nil {
		return fmt.Errorf("error while sending text: %s", err)
	}

	if err := json.Unmarshal(response, message); err != nil {
		return fmt.Errorf("cannot unmarshal response from API: %s", err)
	}

	return nil
}

func (c *Client) EditMessage(message *Message) error {
	params := url.Values{
		"msgId":  {message.ID},
		"chatId": {message.Chat.ID},
		"text":   {message.Text},
	}

	if message.InlineKeyboard != nil {
		data, err := json.Marshal(message.InlineKeyboard.GetKeyboard())
		if err != nil {
			return fmt.Errorf("cannot marshal inline keyboard markup: %s", err)
		}

		params.Set("inlineKeyboardMarkup", string(data))
	}

	response, err := c.Do("/messages/editText", params, nil)
	if err != nil {
		return fmt.Errorf("error while editing text: %s", err)
	}

	if err := json.Unmarshal(response, message); err != nil {
		return fmt.Errorf("cannot unmarshal response from API: %s", err)
	}

	return nil
}

func (c *Client) DeleteMessage(message *Message) error {
	params := url.Values{
		"msgId":  {message.ID},
		"chatId": {message.Chat.ID},
	}
	_, err := c.Do("/messages/deleteMessages", params, nil)
	if err != nil {
		return fmt.Errorf("error while deleting message: %s", err)
	}

	return nil
}

func (c *Client) SendFileMessage(message *Message) error {
	params := url.Values{
		"chatId":  {message.Chat.ID},
		"caption": {message.Text},
		"fileId":  {message.FileID},
	}

	if message.ReplyMsgID != "" {
		params.Set("replyMsgId", message.ReplyMsgID)
	}

	if message.ForwardMsgID != "" {
		params.Set("forwardMsgId", message.ForwardMsgID)
		params.Set("forwardChatId", message.ForwardChatID)
	}

	if message.InlineKeyboard != nil {
		data, err := json.Marshal(message.InlineKeyboard.GetKeyboard())
		if err != nil {
			return fmt.Errorf("cannot marshal inline keyboard markup: %s", err)
		}

		params.Set("inlineKeyboardMarkup", string(data))
	}

	response, err := c.Do("/messages/sendFile", params, nil)
	if err != nil {
		return fmt.Errorf("error while making request: %s", err)
	}

	if err := json.Unmarshal(response, message); err != nil {
		return fmt.Errorf("cannot unmarshal response: %s", err)
	}

	return nil
}

func (c *Client) SendVoiceMessage(message *Message) error {
	params := url.Values{
		"chatId":  {message.Chat.ID},
		"caption": {message.Text},
		"fileId":  {message.FileID},
	}

	if message.ReplyMsgID != "" {
		params.Set("replyMsgId", message.ReplyMsgID)
	}

	if message.ForwardMsgID != "" {
		params.Set("forwardMsgId", message.ForwardMsgID)
		params.Set("forwardChatId", message.ForwardChatID)
	}

	if message.InlineKeyboard != nil {
		data, err := json.Marshal(message.InlineKeyboard.GetKeyboard())
		if err != nil {
			return fmt.Errorf("cannot marshal inline keyboard markup: %s", err)
		}

		params.Set("inlineKeyboardMarkup", string(data))
	}

	response, err := c.Do("/messages/sendVoice", params, nil)
	if err != nil {
		return fmt.Errorf("error while making request: %s", err)
	}

	if err := json.Unmarshal(response, message); err != nil {
		return fmt.Errorf("cannot unmarshal response: %s", err)
	}

	return nil
}

func (c *Client) UploadFile(message *Message) error {
	params := url.Values{
		"chatId":  {message.Chat.ID},
		"caption": {message.Text},
	}

	if message.InlineKeyboard != nil {
		data, err := json.Marshal(message.InlineKeyboard.GetKeyboard())
		if err != nil {
			return fmt.Errorf("cannot marshal inline keyboard markup: %s", err)
		}

		params.Set("inlineKeyboardMarkup", string(data))
	}

	response, err := c.Do("/messages/sendFile", params, message.File)
	if err != nil {
		return fmt.Errorf("error while making request: %s", err)
	}

	if err := json.Unmarshal(response, message); err != nil {
		return fmt.Errorf("cannot unmarshal response: %s", err)
	}

	return nil
}

func (c *Client) UploadVoice(message *Message) error {
	params := url.Values{
		"chatId":  {message.Chat.ID},
		"caption": {message.Text},
	}

	if message.InlineKeyboard != nil {
		data, err := json.Marshal(message.InlineKeyboard.GetKeyboard())
		if err != nil {
			return fmt.Errorf("cannot marshal inline keyboard markup: %s", err)
		}

		params.Set("inlineKeyboardMarkup", string(data))
	}

	response, err := c.Do("/messages/sendVoice", params, message.File)
	if err != nil {
		return fmt.Errorf("error while making request: %s", err)
	}

	if err := json.Unmarshal(response, message); err != nil {
		return fmt.Errorf("cannot unmarshal response: %s", err)
	}

	return nil
}

func (c *Client) GetEvents(lastEventID int, pollTime int) ([]*Event, error) {
	params := url.Values{
		"lastEventId": {strconv.Itoa(lastEventID)},
		"pollTime":    {strconv.Itoa(pollTime)},
	}
	events := &eventsResponse{}

	response, err := c.Do("/events/get", params, nil)
	if err != nil {
		return events.Events, fmt.Errorf("error while making request: %s", err)
	}

	if err := json.Unmarshal(response, events); err != nil {
		return events.Events, fmt.Errorf("cannot parse events: %s", err)
	}

	return events.Events, nil
}

func (c *Client) PinMessage(message *Message) error {
	params := url.Values{
		"chatId": {message.Chat.ID},
		"msgId":  {message.ID},
	}
	_, err := c.Do("/chats/pinMessage", params, nil)
	if err != nil {
		return fmt.Errorf("error while pinning message: %s", err)
	}

	return nil
}

func (c *Client) UnpinMessage(message *Message) error {
	params := url.Values{
		"chatId": {message.Chat.ID},
		"msgId":  {message.ID},
	}
	_, err := c.Do("/chats/unpinMessage", params, nil)
	if err != nil {
		return fmt.Errorf("error while unpinning message: %s", err)
	}

	return nil
}

func (c *Client) SendAnswerCallbackQuery(answer *ButtonResponse) error {
	params := url.Values{
		"queryId":   {answer.QueryID},
		"text":      {answer.Text},
		"url":       {answer.URL},
		"showAlert": {strconv.FormatBool(answer.ShowAlert)},
	}

	_, err := c.Do("/messages/answerCallbackQuery", params, nil)
	if err != nil {
		return fmt.Errorf("error while making request: %s", err)
	}

	return nil
}

func NewClient(baseURL string, token string, logger *logrus.Logger) *Client {
	return &Client{
		token:   token,
		baseURL: baseURL,
		client:  http.DefaultClient,
		logger:  logger,
	}
}
