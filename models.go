package fbmessenger

import (
	"strings"
)

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

/*
ImageMessage is a fluent helper method for creating a SendRequest containing a message with
an image attached using the URL of the image.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/image-attachment
*/
func ImageMessage(url string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "image",
				Payload: ResourcePayload{
					URL: url,
				},
			},
		},
	}
}

/*
ImageDataMessage is a fluent helper method for creating a SendRequest containing a message
with an image attached by uploading the bytes of the image.

	imageBytes, _ := ioutil.ReadFile("./cool-pic.png")
	request := ImageDataMessage(imageBytes, "image/png").To("USER_ID")

Detecting the content type of a file dynamically, or converting to one of the image formats
supported by Facebook is the responsibility of the user.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/image-attachment
*/
func ImageDataMessage(data []byte, contentType string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "image",
				Payload: DataPayload{
					Data:        data,
					ContentType: contentType,
					FileName:    strings.Replace(contentType, "/", ".", -1), //Should yield file names like "image.png"
				},
			},
		},
	}
}

/*
ButtonTemplateMessage is a fluent helper method for creating a SendRequest containing text
and buttons to request input from the user.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/button-template
*/
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

/*
GenericTemplateMessage is a fluent helper method for creating a SendRequest containing
a carousel of elements, each composed of an image attachment, short description and
buttons to request input from the user.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/generic-template
*/
func GenericTemplateMessage(elements ...*GenericPayloadElement) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: GenericPayload{
					TemplateType: "generic",
					Elements:     elements,
				},
			},
		},
	}
}

// URLButton is a fluent helper method for creating a button with type "web_url" for
// use in a message with a button template or generic template attachment.
func URLButton(title, url string) *Button {
	return &Button{
		Type:  "web_url",
		Title: title,
		URL:   url,
	}
}

// PostbackButton is a fluent helper method for creating a button with type "payload" for
// use in a message with a button template or generic template attachment.
func PostbackButton(title, payload string) *Button {
	return &Button{
		Type:    "postback",
		Title:   title,
		Payload: payload,
	}
}

// To is a fluent helper method for setting Recipient. It is a mutator
// and returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) To(userId string) *SendRequest {
	sr.Recipient = Recipient{Id: userId}

	return sr
}

// ToPhoneNumber is a fluent helper method for setting Recipient. It
// is a mutator and returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) ToPhoneNumber(phoneNumber string) *SendRequest {
	sr.Recipient = Recipient{PhoneNumber: phoneNumber}
	return sr
}

// Regular is a fluent helper method for setting NotificationType. It is a mutator and
// returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) Regular() *SendRequest {
	sr.NotificationType = "REGULAR"

	return sr
}

// SilentPush is a fluent helper method for setting NotificationType. It is a mutator and
// returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) SilentPush() *SendRequest {
	sr.NotificationType = "SILENT_PUSH"

	return sr
}

// NoPush is a fluent helper method for setting NotificationType. It is a mutator and
// returns the same SendRequest on which it is called to support method chaining.
func (sr *SendRequest) NoPush() *SendRequest {
	sr.NotificationType = "NO_PUSH"

	return sr
}

/*
SendRequest is the top level structure for representing any type of message to send.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference#request
*/
type SendRequest struct {
	Recipient        Recipient `json:"recipient" binding:"required"`
	Message          Message   `json:"message" binding:"required"`
	NotificationType string    `json:"notification_type,omitempty"`
}

// Recipient identifies the user to send to. Either Id or PhoneNumber must be set, but not both.
type Recipient struct {
	Id          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// Message can represent either a text message, or a message with an attachment. Either
// Text or Attachment must be set, but not both.
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
ResourcePayload is used to hold the URL of a resource (image, file, etc.) to attach to a message.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/image-attachment
*/
type ResourcePayload struct {
	URL string `json:"url" binding:"required"`
}

/*
DataPayload is used to hold the bytes of a resource (image, file, etc.) to upload and attach
to a message. All fields are required. FileName will only be visible to the recipient when
type of the attachment is "file".

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/image-attachment
*/
type DataPayload struct {
	Data        []byte `json:"-"`
	ContentType string `json:"-"`
	FileName    string `json:"-"`
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
	URL     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

/*
GenericPayload is used to build a structured message using the generic template.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/generic-template
*/
type GenericPayload struct {
	TemplateType string                   `json:"template_type" binding:"required"`
	Elements     []*GenericPayloadElement `json:"elements" binding:"required"`
}

// GenericPayloadElement represents one item in the carousel of a generic template message.
type GenericPayloadElement struct {
	Title    string    `json:"title" binding:"required"`
	ImageURL string    `json:"image_url": binding:"required"`
	Subtitle string    `json:"subtitle" binding:"required"`
	Buttons  []*Button `json:"buttons" binding:"required"`
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

/*
Callback is the top level structure that represents a callback received by your
webhook endpoint.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference#format
*/
type Callback struct {
	Object  string   `json:"object" binding:"required"`
	Entries []*Entry `json:"entry" binding:"required"`
}

// Entry is part of the common format of callbacks.
type Entry struct {
	PageId    string            `json:"id" binding:"required"`
	Time      int               `json:"time" binding:"required"`
	Messaging []*MessagingEntry `json:"messaging"`
}

/*
MessagingEntry is an individual interaction a user has with a page.
The Sender and Recipient fields are common to all types of callbacks and the
other fields only apply to specific types of callbacks.
*/
type MessagingEntry struct {
	Sender    Principal        `json:"sender" binding:"required"`
	Recipient Principal        `json:"recipient" binding:"required"`
	Timestamp int              `json:"timestamp"`
	Message   *CallbackMessage `json:"message"`
	Delivery  *Delivery        `json:"delivery"`
	Postback  *Postback        `json:"postback"`
	OptIn     *OptIn           `json:"optin"`
}

// Principal holds the Id of a sender or recipient.
type Principal struct {
	Id string `json:"id" binding:"required"`
}

/*
CallbackMessage represents a message a user has sent to your page.
Either the Text or Attachments field will be set, but not both.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference/message-received
*/
type CallbackMessage struct {
	MessageId   string                `json:"mid" binding:"required"`
	Sequence    int                   `json:"seq" binding:"required"`
	Text        string                `json:"text"`
	Attachments []*CallbackAttachment `json:"attachments"`
}

// CallbackAttachment holds the type and payload of an attachment sent by a user.
type CallbackAttachment struct {
	Type    string                    `json:"type" binding:"required"`
	Payload CallbackAttachmentPayload `json:"payload" binding:"required"`
}

// CallbackAttachmentPayload holds the URL of an attachment sent by the user.
type CallbackAttachmentPayload struct {
	URL string `json:"url" binding:"required"`
}

/*
Delivery holds information about which of the messages that you've sent have been delivered.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference/message-delivered
*/
type Delivery struct {
	MessageIds []string `json:"mids"`
	Watermark  int      `json:"watermark" binding:"required"`
	Sequence   int      `json:"seq" bindging:"required"`
}

/*
Postback holds the data defined for buttons the user taps.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference/postback-received
*/
type Postback struct {
	Payload string `json:"payload" binding:"required"`
}

/*
OptIn holds the data defined for the Send-to-Messenger plugin.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference/authentication
*/
type OptIn struct {
	Ref string `json:"ref" binding:"required"`
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
	ProfilePhotoURL string `json:"profile_pic"`
	Locale          string `json:"locale"`
	Timezone        int    `json:"timezone"`
	Gender          string `json:"gender"`
}
