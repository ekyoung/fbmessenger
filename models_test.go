package fbmessenger_test

import (
	. "github.com/ekyoung/fbmessenger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"fmt"
	"io/ioutil"
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
			Expect(attachment.Payload.Url).To(Equal("IMAGE_URL"))
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
		sendRequest := NewSendRequest("USER_ID").
			WithTextMessage("Hello, world!")

		expectCorrectSendRequestMarshaling(sendRequest, "text-message.json")
	})

	It("should marshal a send request with an image attachment", func() {
		sendRequest := NewSendRequest("USER_ID").
			WithImageUrl("IMAGE_URL")

		expectCorrectSendRequestMarshaling(sendRequest, "message-with-image-attachment.json")
	})

	It("should marshal a send request with a button attachment", func() {
		showWebsite := &Button{
			Type:  "web_url",
			Url:   "https://petersapparel.parseapp.com",
			Title: "Show Website",
		}

		startChatting := &Button{
			Type:    "postback",
			Title:   "Start Chatting",
			Payload: "USER_DEFINED_PAYLOAD",
		}

		sendRequest := NewSendRequest("USER_ID").
			WithButtonTemplate("What do you want to do next?", showWebsite, startChatting)

		expectCorrectSendRequestMarshaling(sendRequest, "message-with-button-attachment.json")
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

	return string(fileBytes)
}

func expectCorrectSendRequestMarshaling(sendRequest *SendRequest, fileName string) {
	sendBytes, err := json.MarshalIndent(sendRequest, "", "  ")

	if err != nil {
		Fail(fmt.Sprintf("Error marshaling send request: %v", err))
	}

	sendString := string(sendBytes)

	Expect(sendString).To(Equal(loadSendRequestString(fileName)))
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
