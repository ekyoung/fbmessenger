package fbmessenger_test

import (
	. "github.com/ekyoung/fbmessenger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
)

var _ = Describe("Callback Models", func() {
	Describe("Message Model", func() {
		It("should unmarshal a callback with a text message", func() {
			var cb Callback
			loadCallback("text-message.json", &cb)

			Expect(cb.Entries[0].Messaging[0].Message.Text).To(Equal("hello, world!"))
		})

		It("should unmarshal a callback with a message with an attachment", func() {
			var cb Callback
			loadCallback("message-with-attachment.json", &cb)

			Expect(cb.Entries[0].Messaging[0].Message.Attachments).ToNot(BeNil())
			attachment := cb.Entries[0].Messaging[0].Message.Attachments[0]
			Expect(attachment.Type).To(Equal("image"))
			Expect(attachment.Payload.URL).To(Equal("IMAGE_URL"))
		})
	})

	Describe("Delivery Model", func() {
		It("should unmarshal a delivery callback", func() {
			var cb Callback
			loadCallback("delivery.json", &cb)
			Expect(len(cb.Entries[0].Messaging[0].Delivery.MessageIds)).To(Equal(1))
			Expect(cb.Entries[0].Messaging[0].Delivery.MessageIds[0]).To(Equal("mid.1458668856218:ed81099e15d3f4f233"))
		})
	})

	Describe("Postback Model", func() {
		It("should unmarshal a postback callback", func() {
			var cb Callback
			loadCallback("postback.json", &cb)
			Expect(cb.Entries[0].Messaging[0].Postback.Payload).To(Equal("USER_DEFINED_PAYLOAD"))
		})
	})

	Describe("Authentication Model", func() {
		It("should unmarshal an authentication callback", func() {
			var cb Callback
			loadCallback("authentication.json", &cb)
			Expect(cb.Entries[0].Messaging[0].OptIn.Ref).To(Equal("PASS_THROUGH_PARAM"))
		})
	})
})

var _ = Describe("Send API Models", func() {
	It("should marshal a send request with a text message", func() {
		sendRequest := TextMessage("Hello, world!").To("USER_ID")

		expectCorrectMarshaling(sendRequest, "text-message.json")
	})

	It("should marshal a send request with an image attached using the URL of the image", func() {
		sendRequest := ImageMessage("IMAGE_URL").To("USER_ID")

		expectCorrectMarshaling(sendRequest, "message-with-image-attachment.json")
	})

	It("should marshal an a message with an image attached by uploading the image", func() {
		imageBytes, err := ioutil.ReadFile("./sample-send-api-data/fb-logo.png")
		if err != nil {
			Fail(fmt.Sprintf("Error reading image file: %v", err))
		}

		sendRequest := ImageDataMessage(imageBytes, "image/png")

		expectCorrectMarshaling(sendRequest.Message, "message-with-image-upload-attachment.json")
	})

	It("should marshal a send request with a button attachment", func() {
		sendRequest := ButtonTemplateMessage("What do you want to do next?",
			URLButton("Show Website", "https://petersapparel.parseapp.com"),
			PostbackButton("Start Chatting", "USER_DEFINED_PAYLOAD")).
			To("USER_ID")

		expectCorrectMarshaling(sendRequest, "message-with-button-attachment.json")
	})

	It("should marshal a send request with a generic attachment", func() {
		viewWebsite := URLButton("View Website", "https://petersapparel.parseapp.com/view_item?item_id=100")

		startChatting := PostbackButton("Start Chatting", "USER_DEFINED_PAYLOAD")

		welcome := &GenericPayloadElement{
			Title:    "Welcome to Peter's Hats",
			ImageURL: "http://petersapparel.parseapp.com/img/item100-thumb.png",
			Subtitle: "We've got the right hat for everyone.",
			Buttons:  []*Button{viewWebsite, startChatting},
		}

		sendRequest := GenericTemplateMessage(welcome).To("USER_ID")

		expectCorrectMarshaling(sendRequest, "message-with-generic-template-attachment.json")
	})

	It("should marshal a send request to a phone number", func() {
		sendRequest := TextMessage("Hello, world!").ToPhoneNumber("+1(212)555-2368")

		expectCorrectMarshaling(sendRequest, "text-message-to-phone-number.json")
	})

	It("should marshal a send request with a REGULAR notification type", func() {
		sendRequest := TextMessage("Hello, world!").To("USER_ID").Regular()

		expectCorrectMarshaling(sendRequest, "text-message-regular.json")
	})

	It("should marshal a send request with a SILENT_PUSH notification type", func() {
		sendRequest := TextMessage("Hello, world!").To("USER_ID").SilentPush()

		expectCorrectMarshaling(sendRequest, "text-message-silent-push.json")
	})

	It("should marshal a send request with a NO_PUSH notification type", func() {
		sendRequest := TextMessage("Hello, world!").To("USER_ID").NoPush()

		expectCorrectMarshaling(sendRequest, "text-message-no-push.json")
	})

	It("should unmarshal a successful response", func() {
		var response SendResponse
		loadSendResponse("successful-response.json", &response)

		Expect(response.Error).To(BeNil())
		Expect(response.RecipientId).To(Equal("1008372609250235"))
		Expect(response.MessageId).To(Equal("mid.1456970487936:c34767dfe57ee6e339"))
	})

	It("should unmarshal an error response", func() {
		var response SendResponse
		loadSendResponse("error-response.json", &response)

		Expect(response.Error.Message).To(Equal("Invalid parameter"))
	})
})

func loadCallback(fileName string, cb *Callback) {
	fileBytes, err := ioutil.ReadFile("./sample-callback-data/" + fileName)
	if err != nil {
		Fail(fmt.Sprintf("Error reading file \"%v\": %v", fileName, err))
	}

	err = json.Unmarshal(fileBytes, cb)
	if err != nil {
		Fail(fmt.Sprintf("File contents is not a valid callback: %v", err))
	}
}

func loadSendRequestString(fileName string) string {
	fileBytes, err := ioutil.ReadFile("./sample-send-api-data/" + fileName)
	if err != nil {
		Fail(fmt.Sprintf("Error reading file \"%v\": %v", fileName, err))
	}

	return makeOneLine(string(fileBytes))
}

func expectCorrectMarshaling(v interface{}, fileName string) {
	sendBytes, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		Fail(fmt.Sprintf("Error marshaling value: %v", err))
	}

	sendString := makeOneLine(string(sendBytes))

	Expect(sendString).To(Equal(loadSendRequestString(fileName)))
}

func makeOneLine(s string) string {
	s = strings.Replace(s, "\r\n", "", -1)
	s = strings.Replace(s, "\n", "", -1)
	return s
}

func loadSendResponse(fileName string, response *SendResponse) {
	fileBytes, err := ioutil.ReadFile("./sample-send-api-data/" + fileName)
	if err != nil {
		Fail(fmt.Sprintf("Error reading file \"%v\": %v", fileName, err))
	}

	err = json.Unmarshal(fileBytes, response)
	if err != nil {
		Fail(fmt.Sprintf("File contents is not a valid send response: %v", err))
	}
}
