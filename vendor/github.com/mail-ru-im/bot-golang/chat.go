package botgolang

//go:generate easyjson -all chat.go

type ChatAction = string

const (
	TypingAction  ChatAction = "typing"
	LookingAction ChatAction = "looking"
)

type ChatType = string

const (
	Private ChatType = "private"
	Group   ChatType = "group"
	Channel ChatType = "channel"
)

type Chat struct {
	client *Client
	// Id of the chat
	ID string `json:"chatId"`

	// Type of the chat: channel, group or private
	Type ChatType `json:"type"`

	// First name of the user
	FirstName string `json:"firstName"`

	// Last name of the user
	LastName string `json:"lastName"`

	// Nick of the user
	Nick string `json:"nick"`

	// User about or group/channel description
	About string `json:"about"`

	// Rules of the group/channel
	Rules string `json:"rules"`

	// Title of the chat
	Title string `json:"title"`

	// Flag that indicates that requested chat is the bot
	IsBot bool `json:"isBot"`

	// Is this chat public?
	Public bool `json:"public"`

	// Is this chat has join moderation?
	JoinModeration bool `json:"joinModeration"`

	// You can send this link to all your friends
	InviteLink string `json:"inviteLink"`
}

func (c *Chat) resolveID() string {
	switch c.Type {
	case Private:
		return c.Nick
	default:
		return c.ID
	}
}

// Send bot actions to the chat
//
// You can call this method every time you change the current actions,
// or every 10 seconds if the actions have not changed. After sending a
// request without active action, you should not re-notify of their absence.
func (c *Chat) SendActions(actions ...ChatAction) error {
	return c.client.SendChatActions(c.resolveID(), actions...)
}

// Get chat administrators list
func (c *Chat) GetAdmins() ([]ChatMember, error) {
	return c.client.GetChatAdmins(c.ID)
}

// Get chat members list
func (c *Chat) GetMembers() ([]ChatMember, error) {
	return c.client.GetChatMembers(c.ID)
}

// Get chat blocked users list
func (c *Chat) GetBlockedUsers() ([]User, error) {
	return c.client.GetChatBlockedUsers(c.ID)
}

// Get chat join pending users list
func (c *Chat) GetPendingUsers() ([]User, error) {
	return c.client.GetChatPendingUsers(c.ID)
}

// Block user and remove him from chat.
// If deleteLastMessages is true, the messages written recently will be deleted
func (c *Chat) BlockUser(userID string, deleteLastMessages bool) error {
	return c.client.BlockChatUser(c.ID, userID, deleteLastMessages)
}

// Unblock user in chat (but not add him back)
func (c *Chat) UnblockUser(userID string) error {
	return c.client.UnblockChatUser(c.ID, userID)
}

// ResolveJoinRequest resolve specific user chat join request
func (c *Chat) ResolveJoinRequest(userID string, accept bool) error {
	return c.client.ResolveChatPending(c.ID, userID, accept, false)
}

// ResolveAllJoinRequest resolve all chat join requests
func (c *Chat) ResolveAllJoinRequests(accept bool) error {
	return c.client.ResolveChatPending(c.ID, "", accept, true)
}

// SetTitle changes chat title
func (c *Chat) SetTitle(title string) error {
	return c.client.SetChatTitle(c.ID, title)
}

// SetAbout changes chat about
func (c *Chat) SetAbout(about string) error {
	return c.client.SetChatAbout(c.ID, about)
}

// SetRules changes chat rules
func (c *Chat) SetRules(rules string) error {
	return c.client.SetChatRules(c.ID, rules)
}
