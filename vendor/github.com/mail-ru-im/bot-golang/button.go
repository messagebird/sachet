package botgolang

//go:generate easyjson -all button.go

// Button represents a button in inline keyboard
// Make sure you have URL or CallbackData in your Button.
type Button struct {
	// Button text
	Text string `json:"text"`

	// URL to be opened
	// You can't use it with CallbackData
	URL string `json:"url,omitempty"`

	// Data that identify the button
	// You can't use it with URL
	CallbackData string `json:"callbackData,omitempty"`

	// Style of a button
	Style ButtonStyle `json:"style,omitempty"`
}

// ButtonStyle represent a style of a Button
type ButtonStyle string

const (
	ButtonPrimary   ButtonStyle = "primary"
	ButtonAttention ButtonStyle = "attention"
)

// WithStyle sets ButtonStyle for Button
func (v Button) WithStyle(style ButtonStyle) Button {
	v.Style = style
	return v
}

// NewURLButton returns new button with URL field
func NewURLButton(text string, url string) Button {
	return Button{
		Text: text,
		URL:  url,
	}
}

// NewCallbackButton returns new button with CallbackData field
func NewCallbackButton(text string, callbackData string) Button {
	return Button{
		Text:         text,
		CallbackData: callbackData,
	}
}

// ButtonResponse represents a data that is returned when a button is clicked
type ButtonResponse struct {
	client *Client

	// Id of the query
	QueryID string `json:"queryId"`

	// Text of the response message
	Text string `json:"text"`

	// Display alert?
	ShowAlert bool `json:"showAlert"`

	// URL to be opened
	URL string `json:"url"`

	// CallbackData of the query (id of the pressed button).
	CallbackData string `json:"callbackData"`
}

// Send method sends your response message.
// Make sure you have QueryID in your ButtonResponse.
func (cl *ButtonResponse) Send() error {
	return cl.client.SendAnswerCallbackQuery(cl)
}
