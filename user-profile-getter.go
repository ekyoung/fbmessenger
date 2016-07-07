package fbmessenger

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type httpGetter interface {
	Get(url string) (*http.Response, error)
}

type UserProfileGetter interface {
	Get(userId string) (*UserProfile, error)
}

type userProfileGetter struct {
	pageAccessToken string
	httpGetter      httpGetter
}

func NewUserProfileGetter(pageAccessToken string) UserProfileGetter {
	return &userProfileGetter{
		pageAccessToken: pageAccessToken,
		httpGetter:      &http.Client{},
	}
}

func (upg *userProfileGetter) Get(userId string) (*UserProfile, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v2.6/%v?fields=first_name,last_name,profile_pic,locale,timezone,gender&access_token=%v", userId, upg.pageAccessToken)
	resp, err := upg.httpGetter.Get(url)
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
