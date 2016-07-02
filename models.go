package fbmessenger

/*
Callback

All callbacks will contain the top level info and an
"entry" array containing a "messaging" array.
*/

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

type Principal struct {
	Id string `json:"id" binding:"required"`
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

/*
Send API

For sending messages.
*/
type SendRequest struct {
	Recipient Principal `json:"recipient" binding:"required"`
	Message   Message   `json:"message" binding:"required"`
}

type Message struct {
	Text       string      `json:"text,omitempty"`
	Attachment *Attachment `json:"attachment,omitempty"`
}

type Attachment struct {
	Type    string      `json:"type" binding:"required"`
	Payload interface{} `json:"payload" binding:"required"`
}

type ImagePayload struct {
	Url string `json:"url" binding:"required"`
}

type ButtonPayload struct {
	TemplateType string    `json:"template_type" binding:"required"`
	Text         string    `json:"text" binding:"required"`
	Buttons      []*Button `json:"buttons" binding:"required"`
}

type Button struct {
	Type    string `json:"type" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Url     string `json:"url,omitempty"`
	Payload string `json:"payload,omitempty"`
}

type SendResponse struct {
	RecipientId string     `json:"recipient_id" binding:"required"`
	MessageId   string     `json:"message_id" binding:"required"`
	Error       *SendError `json:"error"`
}

type SendError struct {
	Message   string `json:"message" binding:"required"`
	Type      string `json:"type" binding:"required"`
	Code      int    `json:"code" binding:"required"`
	ErrorData string `json:"error_data" binding:"required"`
	FBTraceId string `json:"fbtrace_id" binding:"required"`
}

func NewSendRequest(userId string) *SendRequest {
	return &SendRequest{
		Recipient: Principal{Id: "USER_ID"},
	}
}

func (sr *SendRequest) WithTextMessage(text string) *SendRequest {
	sr.Message = Message{
		Text: text,
	}

	return sr
}

func (sr *SendRequest) WithImageUrl(url string) *SendRequest {
	sr.Message = Message{
		Attachment: &Attachment{
			Type: "image",
			Payload: ImagePayload{
				Url: url,
			},
		},
	}

	return sr
}

func (sr *SendRequest) WithButtonTemplate(text string, buttons ...*Button) *SendRequest {
	sr.Message = Message{
		Attachment: &Attachment{
			Type: "template",
			Payload: ButtonPayload{
				TemplateType: "button",
				Text:         text,
				Buttons:      buttons,
			},
		},
	}

	return sr
}

/*
User Profile API

For getting info about the user.
*/
type UserProfile struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	ProfilePhotoUrl string `json:"profile_pic"`
	Locale          string `json:"locale"`
	Timezone        int    `json:"timezone"`
	Gender          string `json:"gender"`
}
