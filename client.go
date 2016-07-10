package fbmessenger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net/http"
)

const apiURL = "https://graph.facebook.com/v2.6/"

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type Client struct {
	URL      string
	httpDoer httpDoer
}

func (c *Client) Send(sendRequest *SendRequest, pageAccessToken string) (*SendResponse, error) {
	return c.SendWithContext(context.Background(), sendRequest, pageAccessToken)
}

func (c *Client) SendWithContext(ctx context.Context, sendRequest *SendRequest, pageAccessToken string) (*SendResponse, error) {
	requestBytes, err := json.Marshal(sendRequest)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", c.buildURL("me/messages?access_token="+pageAccessToken), bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	response := &SendResponse{}
	err = c.doRequest(ctx, req, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Client) GetUserProfile(userId, pageAccessToken string) (*UserProfile, error) {
	return c.GetUserProfileWithContext(context.Background(), userId, pageAccessToken)
}

func (c *Client) GetUserProfileWithContext(ctx context.Context, userId, pageAccessToken string) (*UserProfile, error) {
	url := c.buildURL(fmt.Sprintf("%v?fields=first_name,last_name,profile_pic,locale,timezone,gender&access_token=%v", userId, pageAccessToken))

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
