package fbmessenger

import (
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"net/http"
)

type httpDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type UserProfileGetter interface {
	Get(userId string) (*UserProfile, error)
	GetWithContext(ctx context.Context, userId string) (*UserProfile, error)
}

type userProfileGetter struct {
	pageAccessToken string
	httpDoer        httpDoer
}

func NewUserProfileGetter(pageAccessToken string) UserProfileGetter {
	return &userProfileGetter{
		pageAccessToken: pageAccessToken,
		httpDoer:        &http.Client{},
	}
}

func (upg *userProfileGetter) Get(userId string) (*UserProfile, error) {
	return upg.GetWithContext(context.Background(), userId)
}

func (upg *userProfileGetter) GetWithContext(ctx context.Context, userId string) (*UserProfile, error) {
	url := fmt.Sprintf("https://graph.facebook.com/v2.6/%v?fields=first_name,last_name,profile_pic,locale,timezone,gender&access_token=%v", userId, upg.pageAccessToken)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Cancel = ctx.Done()

	resp, err := upg.httpDoer.Do(req)
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
