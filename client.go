package fbmessenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/textproto"

	"golang.org/x/net/context"
)

const apiURL = "https://graph.facebook.com/v3.3"

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

/*
Client is used to send messages and get user profiles. Use the empty value in most cases.
The URL field can be overridden to allow for writing integration tests that use a different
endpoint (not Facebook).
*/
type Client struct {
	URL      string
	httpDoer httpDoer
}

/*
Send POSTs a request to and returns a response from the Send API. A response from
Facebook indicating an error does not return an error. Be sure to check for errors
in sending, and errors in the response from Facebook.

	response, err := client.Send(request, "YOUR_PAGE_ACCESS_TOKEN")
	if err != nil {
		//Got an error. Request never got to Facebook.
	} else if response.Error != nil {
		//Request got to Facebook. Facebook returned an error.
	} else {
		//Hooray!
	}

*/
func (c *Client) Send(sendRequest *SendRequest, pageAccessToken string) (*SendResponse, error) {
	return c.SendWithContext(context.Background(), sendRequest, pageAccessToken)
}

// SendWithContext is like Send but allows you to timeout or cancel the request using context.Context.
func (c *Client) SendWithContext(ctx context.Context, sendRequest *SendRequest, pageAccessToken string) (*SendResponse, error) {
	var req *http.Request
	var err error

	if isDataMessage(sendRequest) {
		req, err = c.newFormDataRequest(sendRequest, pageAccessToken)
	} else {
		req, err = c.newJSONRequest(sendRequest, pageAccessToken)
	}

	if err != nil {
		return nil, err
	}

	response := &SendResponse{}
	err = c.doRequest(ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func isDataMessage(sendRequest *SendRequest) bool {
	if sendRequest.Message.Attachment == nil {
		return false
	}

	_, ok := sendRequest.Message.Attachment.Payload.(DataPayload)

	return ok
}

func (c *Client) newJSONRequest(sendRequest *SendRequest, pageAccessToken string) (*http.Request, error) {
	requestBytes, err := json.Marshal(sendRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.buildURL("/me/messages?access_token="+pageAccessToken), bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) newFormDataRequest(sendRequest *SendRequest, pageAccessToken string) (*http.Request, error) {
	payload, _ := sendRequest.Message.Attachment.Payload.(DataPayload)

	var reqBuffer bytes.Buffer
	w := multipart.NewWriter(&reqBuffer)

	err := writeFormField(w, "recipient", sendRequest.Recipient)
	if err != nil {
		return nil, err
	}

	err = writeFormField(w, "message", sendRequest.Message)
	if err != nil {
		return nil, err
	}

	if sendRequest.NotificationType != "" {
		err = w.WriteField("notification_type", sendRequest.NotificationType)
		if err != nil {
			return nil, err
		}
	}

	header := make(textproto.MIMEHeader)
	header.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "filedata", payload.FileName))
	header.Set("Content-Type", payload.ContentType)

	fileWriter, err := w.CreatePart(header)
	if err != nil {
		return nil, err
	}

	_, err = fileWriter.Write(payload.Data)
	if err != nil {
		return nil, err
	}

	w.Close()

	req, err := http.NewRequest("POST", c.buildURL("/me/messages?access_token="+pageAccessToken), &reqBuffer)
	//req, err := http.NewRequest("POST", "http://httpbin.org/post", &reqBuffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	return req, nil
}

func writeFormField(w *multipart.Writer, fieldName string, value interface{}) error {
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("error marshaling value for field %v: %v", fieldName, err)
	}

	return w.WriteField(fieldName, string(valueBytes))
}

// GetUserProfile GETs a profile with more information about the user.
func (c *Client) GetUserProfile(userId, pageAccessToken string) (*UserProfile, error) {
	return c.GetUserProfileWithContext(context.Background(), userId, pageAccessToken)
}

// GetUserProfileWithContext is like GetUserProfile but allows you to timeout or cancel the request using context.Context.
func (c *Client) GetUserProfileWithContext(ctx context.Context, userId, pageAccessToken string) (*UserProfile, error) {
	url := c.buildURL(fmt.Sprintf("/%v?fields=first_name,last_name,profile_pic,locale,timezone,gender&access_token=%v", userId, pageAccessToken))

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	userProfile := &UserProfile{}
	err = c.doRequest(ctx, req, userProfile)
	if err != nil {
		return nil, err
	}

	return userProfile, nil
}

func (c *Client) buildURL(path string) string {
	url := c.URL
	if url == "" {
		url = apiURL
	}

	return url + path
}

func (c *Client) doRequest(ctx context.Context, req *http.Request, responseStruct interface{}) error {
	req.Cancel = ctx.Done()

	doer := c.httpDoer
	if doer == nil {
		doer = &http.Client{}
	}

	resp, err := doer.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, responseStruct)
	if err != nil {
		return err
	}

	return nil
}
