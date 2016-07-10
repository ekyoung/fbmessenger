package fbmessenger

/*------------------------------------------------------
Send API
------------------------------------------------------*/

// Text message is a fluent helper method for creating a SendRequest containing a text message.
func TextMessage(text string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Text: text,
		},
	}
}

// ImageMessage is a fluent helper method for creating a SendRequest containing a message with
// an image attachment that has a URL payload.
func ImageMessage(url string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "image",
				Payload: MediaPayload{
					Url: url,
				},
			},
		},
	}
}

// ButtonTemplateMessage is a fluent helper method for creating a SendRequest containing text
// and buttons to request input from the user.
func ButtonTemplateMessage(text string, buttons ...*Button) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: ButtonPayload{
					TemplateType: "button",
					Text:         text,
					Buttons:      buttons,
				},
			},
		},
	}
}

// To is a fluent helper method for setting the Recipient of a SendRequest. It is a mutator
// and returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) To(userId string) *SendRequest {
	sr.Recipient = Principal{Id: userId}

	return sr
}

/*
SendRequest is the top level structure for representing any type of message to send.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference#request
*/
type SendRequest struct {
	Recipient Principal `json:"recipient" binding:"required"`
	Message   Message   `json:"message" binding:"required"`
}

// Message can represent either a text message, or a message with an attachment. Either
// Text or Attachment mut be set, but not both.
type Message struct {
	Text       string      `json:"text,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

// Attachment is used to build a message with attached media, or a structured message.
type Attachment struct {
	Type    string      `json:"type" binding:"required"`
	Payload interface{} `json:"payload" binding:"required"`
}

/*
MediaPayload is used to hold the URL of media attached to a message.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/image-attachment
*/
type MediaPayload struct {
	Url string `json:"url" binding:"required"`
}

/*
ButtonPayload is used to build a structured message using the button template.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/button-template
*/
type ButtonPayload struct {
	TemplateType string    `json:"template_type" binding:"required"`
	Text         string    `json:"text" binding:"required"`
	Buttons      []*Button `json:"buttons" binding:"required"`
}

// Button represents a single button in a structured message using the button template.
type Button struct {
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Url     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

/*
SendResponse is returned when sending a SendRequest.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference#response
*/
type SendResponse struct {
	RecipientId string     `json:"recipient_id" binding:"required"`
	MessageId   string     `json:"message_id" binding:"required"`
	Error       *SendError `json:"error"`
}

/*
SendError indicates an error returned from Facebook.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference#errors
*/
type SendError struct {
	Message   string `json:"message" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Code      int    `json:"code" binding:"required"`
	ErrorData string `json:"error_data" binding:"required"`
	FBTraceId string `json:"fbtrace_id" binding:"required"`
}

/*------------------------------------------------------
Webhook
------------------------------------------------------*/

type Callback struct {
	Object  string   `json:"object" binding:"required"`
	Entries []*Entry `json:"entry" binding:"required"`
}

type Entry struct {
	PageId    string            `json:"id" binding:"required"`
	Time      int               `json:"time" binding:"required"`
	Messaging []*MessagingEntry `json:"messaging"`
}

type MessagingEntry struct {
	Sender    Principal        `json:"sender" binding:"required"`
	Recipient Principal        `json:"recipient" binding:"required"`
	Timestamp int              `json:"timestamp"`
	Message   *CallbackMessage `json:"message"`
	Delivery  *Delivery        `json:"delivery"`
	Postback  *Postback        `json:"postback"`
	OptIn     *OptIn           `json:"optin"`
}

/*
Message Received

Messages can have either "text" or "attachments".
*/

type CallbackMessage struct {
	MessageId   string                `json:"mid" binding:"required"`
	Sequence    int                   `json:"seq" binding:"required"`
	Text        string                `json:"text"`
	Attachments []*CallbackAttachment `json:"attachments"`
}

type CallbackAttachment struct {
	Type    string  `json:"type" binding:"required"`
	Payload Payload `json:"payload" binding:"required"`
}

type Payload struct {
	Url string `json:"url" binding:"required"`
}

/*
Message Delivered
*/

type Delivery struct {
	MessageIds []string `json:"mids"`
	Watermark  int      `json:"watermark" binding:"required"`
	Sequence   int      `json:"seq" bindging:"required"`
}

/*
Postback
*/

type Postback struct {
	Payload string `json:"payload" binding:"required"`
}

/*
Authentication
*/

type OptIn struct {
	Ref string `json:"ref" binding:"required"`
}

/*------------------------------------------------------
Common
------------------------------------------------------*/

type Principal struct {
	Id string `json:"id" binding:"required"`
}

/*------------------------------------------------------
User Profile
------------------------------------------------------*/

/*
UserProfile represents additional information about the user.

See https://developers.facebook.com/docs/messenger-platform/user-profile
*/
type UserProfile struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ProfilePhotoUrl string `json:"profile_pic"`
	Locale          string `json:"locale"`
	Timezone        int    `json:"timezone"`
	Gender          string `json:"gender"`
}
