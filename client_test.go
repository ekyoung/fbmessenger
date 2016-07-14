package fbmessenger_test

import (
	. "github.com/ekyoung/fbmessenger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

	"fmt"
	"io/ioutil"
	"mime"
	"net/http"
)

var _ = Describe("Client", func() {
	Describe("Send", func() {
		const (
			pageAccessToken = "SOME_TOKEN"
			userId          = "USER_ID"
		)

		var (
			server *ghttp.Server

			client *Client
		)

		BeforeEach(func() {
			server = ghttp.NewServer()

			client = &Client{
				URL: server.URL(),
			}
		})

		AfterEach(func() {
			server.Close()
		})

		It("should POST json when sending a text message", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/me/messages"),
					ghttp.VerifyHeader(http.Header{
						"Content-Type": []string{"application/json"},
					}),

					ghttp.RespondWithJSONEncoded(200, &SendResponse{
						RecipientId: userId,
						MessageId:   "mid.12345",
					}),
				),
			)

			request := TextMessage("Hello, world!").To("USER_ID")
			response, err := client.Send(request, pageAccessToken)

			if err != nil {
				Fail(fmt.Sprintf("Error returned: %v", err))
			}

			Expect(response.RecipientId).To(Equal(userId))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("should POST json when sending an image message", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/me/messages"),
					ghttp.VerifyHeader(http.Header{
						"Content-Type": []string{"application/json"},
					}),

					ghttp.RespondWithJSONEncoded(200, &SendResponse{
						RecipientId: userId,
						MessageId:   "mid.12345",
					}),
				),
			)

			request := ImageMessage("http://someurl.com/pic.jpg").To("USER_ID")
			response, err := client.Send(request, pageAccessToken)

			if err != nil {
				Fail(fmt.Sprintf("Error returned: %v", err))
			}

			Expect(response.RecipientId).To(Equal(userId))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("should POST form data when sending an image upload message", func() {
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/me/messages"),

					ghttp.RespondWithJSONEncoded(200, &SendResponse{
						RecipientId: userId,
						MessageId:   "mid.12345",
					}),
				),
			)

			imageBytes, err := ioutil.ReadFile("./sample-send-api-data/fb-logo.png")
			if err != nil {
				Fail(fmt.Sprintf("Error reading image file: %v", err))
			}

			request := ImageUploadMessage(imageBytes).To("USER_ID")
			response, err := client.Send(request, pageAccessToken)

			if err != nil {
				Fail(fmt.Sprintf("Error returned: %v", err))
			}

			Expect(response.RecipientId).To(Equal(userId))

			receivedRequests := server.ReceivedRequests()
			Expect(receivedRequests).To(HaveLen(1))

			receivedRequest := receivedRequests[0]

			mediaType, _, err := mime.ParseMediaType(receivedRequest.Header.Get("Content-Type"))
			if err != nil {
				Fail(fmt.Sprintf("Error parsing content type header: %v", err))
			}

			Expect(mediaType).To(Equal("multipart/form-data"))
		})
	})
})
