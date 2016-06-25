package fbmessenger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Getter interface {
	Get(url string) (*http.Response, error)
}

type UserProfileApi struct {
	PageAccessToken string
	Getter          Getter
}

func NewUserProfileApi(pageAccessToken string) *UserProfileApi {
	return &UserProfileApi{
		PageAccessToken: pageAccessToken,
		Getter:          &http.Client{},
	}
}

func (api *UserProfileApi) Get(userId string) (*UserProfile, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v2.6/%v?fields=first_name,last_name,profile_pic,locale,timezone,gender&access_token=%v", userId, api.PageAccessToken)
	resp, err := api.Getter.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	userProfile := &UserProfile{}
	err = json.Unmarshal(body, userProfile)
	if err != nil {
		return nil, err
	}

	return userProfile, nil
}
