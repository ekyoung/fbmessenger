package fbmessenger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Sender interface {
	Send(sendRequest *SendRequest) (*SendResponse, error)
}

type sender struct {
	url string
}

func NewSender(pageAccessToken string) Sender {
	return &sender{
		url: "https://graph.facebook.com/v2.6/me/messages?access_token=" + pageAccessToken,
	}
}

func (s *sender) Send(sendRequest *SendRequest) (*SendResponse, error) {
	requestBytes, err := json.Marshal(sendRequest)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", s.url, bytes.NewBuffer(requestBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	responseBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	response := &SendResponse{}
	err = json.Unmarshal(responseBytes, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
