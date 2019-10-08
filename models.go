package fbmessenger

import (
	"encoding/json"
	"strings"
)

/*------------------------------------------------------
Send API
------------------------------------------------------*/

// TextMessage is a fluent helper method for creating a SendRequest containing a text message.
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

See https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments
*/
func ImageMessage(url string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "image",
				Payload: ResourcePayload{
					URL:      url,
					Reusable: true,
				},
			},
		},
	}
}

/*
VideoMessage is a fluent helper method for creating a SendRequest containing a message with
an video attached using the URL of the video.

See https://developers.facebook.com/docs/messenger-platform/send-messages#sending_attachments
*/
func VideoMessage(url string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "video",
				Payload: ResourcePayload{
					URL:      url,
					Reusable: true,
				},
			},
		},
	}
}

/*
SavedImageMessage is a fluent helper method for creating a SendRequest containing a message with
an image attached using the identifier of the image asset.

See https://developers.facebook.com/docs/messenger-platform/send-messages#attachment_reuse
*/
func SavedImageMessage(id string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "image",
				Payload: SavedAssetPayload{
					AttachmentID: id,
				},
			},
		},
	}
}

/*
SavedVideoMessage is a fluent helper method for creating a SendRequest containing a message with
an video attached using the identifier of the video asset.

See https://developers.facebook.com/docs/messenger-platform/send-messages#attachment_reuse
*/
func SavedVideoMessage(id string) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "video",
				Payload: SavedAssetPayload{
					AttachmentID: id,
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
WebviewButtonTemplateMessage is a fluent helper method for creating a SendRequest containing text
and webview buttons to request input from the user.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/button-template
*/
func WebviewButtonTemplateMessage(text string, buttons ...*WebviewButton) *SendRequest {
	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type: "template",
				Payload: WebviewButtonPayload{
					TemplateType:   "button",
					Text:           text,
					WebviewButtons: buttons,
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
func GenericTemplateMessage(elements ...*GenericElement) *SendRequest {
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

/*
ReceiptTemplateMessage is a fluent helper method for creating a SendRequest containing
a detailed order confirmation.

https://developers.facebook.com/docs/messenger-platform/send-api-reference/receipt-template
*/
func ReceiptTemplateMessage(header *ReceiptHeader, summary *ReceiptSummary, elements ...*ReceiptElement) *SendRequest {
	payload := &ReceiptPayload{
		TemplateType:  "receipt",
		RecipientName: header.RecipientName,
		OrderNumber:   header.OrderNumber,
		Currency:      header.Currency,
		PaymentMethod: header.PaymentMethod,
		OrderURL:      header.OrderURL,
		Timestamp:     header.Timestamp,
		Elements:      elements,
		Summary:       summary,
	}

	return &SendRequest{
		Message: Message{
			Attachment: &Attachment{
				Type:    "template",
				Payload: payload,
			},
		},
	}
}

// ReceiptHeader holds just the top level fields for a ReceiptPayload. For use with
// the ReceiptTemplateMessage fluent helper method.
type ReceiptHeader struct {
	RecipientName string
	OrderNumber   string
	Currency      string
	PaymentMethod string
	OrderURL      string
	Timestamp     string
}

/*
WithReceiptAddress is a fluent helper method for setting the Address of a ReceiptPayload
for the message. It is a mutator and returns the same SendRequest on which it is
called to support method chaining.
*/
func (sr *SendRequest) WithReceiptAddress(address *Address) *SendRequest {
	receipt, _ := sr.Message.Attachment.Payload.(*ReceiptPayload)

	receipt.Address = address

	return sr
}

/*
WithReceiptAdjustments is a fluent helper method for setting the Adjustments of a
ReceiptPayload for the message. It is a mutator and returns the same SendRequest
on which it is called to support method chaining.
*/
func (sr *SendRequest) WithReceiptAdjustments(adjustments ...*ReceiptAdjustment) *SendRequest {
	receipt, _ := sr.Message.Attachment.Payload.(*ReceiptPayload)

	receipt.Adjustments = adjustments

	return sr
}

// URLAction is a fluent helper method for creating an action with type "web_url" for
// use in a message with generic template attachment.
func URLAction(url string) *Action {
	return &Action{
		Type: "web_url",
		URL:  url,
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

// WebviewURLButton is a fluent helper method for creating a webview button with type "web_url" for
// use in a message with a button template or generic template attachment.
func WebviewURLButton(title, url string) *WebviewButton {
	return &WebviewButton{
		Type:        "web_url",
		Title:       title,
		URL:         url,
		HeightRatio: "compact",
		Extension:   true,
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
func (sr *SendRequest) To(userID string) *SendRequest {
	sr.Recipient = Recipient{ID: userID}

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

// TextReply is a fluent helper method for creating a QuickReply with content type "text".
func TextReply(title, payload string) *QuickReply {
	return &QuickReply{
		ContentType: "text",
		Title:       title,
		Payload:     payload,
	}
}

// TextReplyWithImage is a fluent helper method for creating a QuickReply with content type
// "text" and an attached image.
func TextReplyWithImage(title, payload, imageURL string) *QuickReply {
	return &QuickReply{
		ContentType: "text",
		Title:       title,
		Payload:     payload,
		ImageURL:    imageURL,
	}
}

// LocationReply is a fluent helper method for creating a QuickReply with content type "location".
func LocationReply() *QuickReply {
	return &QuickReply{
		ContentType: "location",
	}
}

// WithQuickReplies is a fluent helper method for setting the quick replies to
// a message. It is not additive, it replaces any existing quick replies.
func (sr *SendRequest) WithQuickReplies(replies ...*QuickReply) *SendRequest {
	sr.Message.QuickReplies = replies

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
	ID          string `json:"id,omitempty"`
	PhoneNumber string `json:"phone_number,omitempty"`
}

// Message can represent either a text message, or a message with an attachment. Either
// Text or Attachment must be set, but not both.
type Message struct {
	Text         string        `json:"text,omitempty"`
	Attachment   *Attachment   `json:"attachment,omitempty"`
	QuickReplies []*QuickReply `json:"quick_replies,omitempty"`
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
	URL      string `json:"url" binding:"required"`
	Reusable bool   `json:"is_reusable" binding:"required"`
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
SavedAssetPayload is used to hold the identifier of an already saved asset (image, file, etc.) to attach to a message.

See https://developers.facebook.com/docs/messenger-platform/send-messages#attachment_reuse
*/
type SavedAssetPayload struct {
	AttachmentID string `json:"attachment_id" binding:"required"`
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

/*
WebviewButtonPayload is used to build a structured message using the button template.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/button-template
*/
type WebviewButtonPayload struct {
	TemplateType   string           `json:"template_type" binding:"required"`
	Text           string           `json:"text" binding:"required"`
	WebviewButtons []*WebviewButton `json:"buttons" binding:"required"`
}

// Button represents a single button in a structured message using the button template.
type Button struct {
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	URL     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

// WebviewButton represents a single webview button in a structured message using the button template.
type WebviewButton struct {
	Type        string `json:"type" binding:"required"`
	Title       string `json:"title" binding:"required"`
	URL         string `json:"url,omitempty"`
	Payload     string `json:"payload,omitempty"`
	HeightRatio string `json:"webview_height_ratio,omitempty"`
	Extension   bool   `json:"messenger_extensions,omitempty"`
}

// Action represents a default action for element in a structured message using the generic template.
type Action struct {
	Type    string `json:"type" binding:"required"`
	URL     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

/*
GenericPayload is used to build a structured message using the generic template.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/generic-template
*/
type GenericPayload struct {
	TemplateType string            `json:"template_type" binding:"required"`
	Elements     []*GenericElement `json:"elements" binding:"required"`
}

// GenericElement represents one item in the carousel of a generic template message.
type GenericElement struct {
	Title         string    `json:"title" binding:"required"`
	ImageURL      string    `json:"image_url" binding:"required"`
	Subtitle      string    `json:"subtitle" binding:"required"`
	DefaultAction *Action   `json:"default_action,omitempty"`
	Buttons       []*Button `json:"buttons" binding:"required"`
}

/*
ReceiptPayload is used to build a structured message using the receipt template.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference/receipt-template
*/
type ReceiptPayload struct {
	TemplateType  string               `json:"template_type" binding:"required"`
	RecipientName string               `json:"recipient_name" binding:"required"`
	OrderNumber   string               `json:"order_number" binding:"required"`
	Currency      string               `json:"currency" binding:"required"`
	PaymentMethod string               `json:"payment_method" binding:"required"`
	Timestamp     string               `json:"timestamp,omitempty"`
	OrderURL      string               `json:"order_url,omitempty"`
	Elements      []*ReceiptElement    `json:"elements" binding:"required"`
	Address       *Address             `json:"address,omitempty"`
	Summary       *ReceiptSummary      `json:"summary" binding:"required"`
	Adjustments   []*ReceiptAdjustment `json:"adjustments,omitempty"`
}

// ReceiptElement represents a line item for one purchased item (not tax or shipping)
// on a receipt.
type ReceiptElement struct {
	Title    string      `json:"title" binding:"required"`
	Subtitle string      `json:"subtitle,omitempty"`
	Quantity int         `json:"quantity,omitempty"`
	Price    json.Number `json:"price" binding:"required"`
	Currency string      `json:"currency,omitempty"`
	ImageURL string      `json:"image_url,omitempty"`
}

// Address represents a physical mailing address
type Address struct {
	Street1    string `json:"street_1" binding:"required"`
	Street2    string `json:"street_2,omitempty"`
	City       string `json:"city" binding:"required"`
	PostalCode string `json:"postal_code" binding:"required"`
	State      string `json:"state" binding:"required"`
	Country    string `json:"country" binding:"required"`
}

// ReceiptSummary represents the line items for totals and additional costs
// (tax and shipping) on a receipt.
type ReceiptSummary struct {
	Subtotal     json.Number `json:"subtotal,omitempty"`
	ShippingCost json.Number `json:"shipping_cost,omitempty"`
	TotalTax     json.Number `json:"total_tax,omitempty"`
	TotalCost    json.Number `json:"total_cost" binding:"required"`
}

// ReceiptAdjustment represents discounts applied to a receipt.
type ReceiptAdjustment struct {
	Name   string      `json:"name,omitempty"`
	Amount json.Number `json:"amount"`
}

// QuickReply represents a quick reply to a message.
type QuickReply struct {
	ContentType string `json:"content_type" binding:"required"`
	Title       string `json:"title,omitempty"`
	Payload     string `json:"payload,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
}

/*
SendResponse is returned when sending a SendRequest.

See https://developers.facebook.com/docs/messenger-platform/send-api-reference#response
*/
type SendResponse struct {
	RecipientID  string     `json:"recipient_id" binding:"required"`
	MessageID    string     `json:"message_id" binding:"required"`
	AttachmentID string     `json:"attachment_id,omitempty"`
	Error        *SendError `json:"error"`
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
	FBTraceID string `json:"fbtrace_id" binding:"required"`
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
	PageID    string            `json:"id" binding:"required"`
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
	ID string `json:"id" binding:"required"`
}

/*
CallbackMessage represents a message a user has sent to your page.
Either the Text or Attachments field will be set, but not both.

See https://developers.facebook.com/docs/messenger-platform/webhook-reference/message-received
*/
type CallbackMessage struct {
	MessageID   string                `json:"mid" binding:"required"`
	Sequence    int                   `json:"seq" binding:"required"`
	Text        string                `json:"text"`
	Attachments []*CallbackAttachment `json:"attachments"`
	QuickReply  *CallbackQuickReply   `json:"quick_reply"`
}

// CallbackAttachment holds the type and payload of an attachment sent by a user.
type CallbackAttachment struct {
	Title   string                    `json:"title"`
	URL     string                    `json:"url"`
	Type    string                    `json:"type" binding:"required"`
	Payload CallbackAttachmentPayload `json:"payload" binding:"required"`
}

// CallbackAttachmentPayload holds the URL of a multimedia attachment,
// or the coordinates of a location attachment sent by the user.
type CallbackAttachmentPayload struct {
	URL         string       `json:"url"`
	Coordinates *Coordinates `json:"coordinates"`
}

// Coordinates holds the latitude and longitude of a location.
type Coordinates struct {
	Lat  float64 `json:"lat" binding:"required"`
	Long float64 `json:"long" binding:"required"`
}

// CallbackQuickReply holds the developer defined payload of a quick reply sent by the user.
type CallbackQuickReply struct {
	Payload string `json:"payload" binding:"required"`
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
