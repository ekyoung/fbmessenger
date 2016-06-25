package fbmessenger

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type SendApi struct {
	PageAccessToken string
}

func NewSendApi(pageAccessToken string) *SendApi {
	return &SendApi{
		PageAccessToken: pageAccessToken,
	}
}

func (api *SendApi) Send(sendRequest *SendRequest) (*SendResponse, error) {
	url := "https://graph.facebook.com/v2.6/me/messages?access_token=" + api.PageAccessToken

	requestBytes, err := json.Marshal(sendRequest)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBytes))
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
