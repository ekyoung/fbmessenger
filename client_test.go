package fbmessenger_test

import (
	. "github.com/ekyoung/fbmessenger"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"

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

			Expect(response.RecipientId).To(Equal(userId))
			Expect(err).To(BeNil())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})
	})
})
