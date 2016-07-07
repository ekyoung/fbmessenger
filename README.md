# Package `fbmessenger`

[![CircleCI](https://circleci.com/gh/ekyoung/fbmessenger.svg?style=svg)](https://circleci.com/gh/ekyoung/fbmessenger)

Go (golang) package for writing bots on the [Facebook Messenger Platform](https://developers.facebook.com/docs/messenger-platform).

## Installation

Expect breaking changes:

```bash
go get https://github.com/ekyoung/fbmessenger
```

Stable installation via [gopkg.in](http://labix.org/gopkg.in) coming soon.

## Quick Start

### Webhook Reference

Unmarshal the json received at your webhook endpoint into an instance of type `Callback`.

```go
cb := &fbmessenger.Callback{}
err := json.Unmarshal(requestBytes, cb)
```

Use type `CallbackDispatcher` to route each `MessageEntry` included in the callback to an appropriate
handler for the type of entry. Note that due to webhook batching, a handler may be called more than
once per callback.

```go
dispatcher := &fbmessenger.CallbackDispatcher{
	MessageHandler: MessageReceived
}

err := dispatcher.Dispatch(cb)
```

Callback handlers should have a signature mathing the `MessageEntryHandler` type.

```go
func MessageReceived(cb *fbmessenger.MessagingEntry) error {
	//Do stuff
}
```

### Send API Reference

Type `Sender` handles sending to the messenger API. Create one using the page access
token for the page you want to send as.

```go
sender := fbmessenger.NewSender("YOUR_PAGE_ACCESS_TOKEN")
```

There are structs for the different types of messages you can send. The easiest way to create them
is with the fluent API.

```go
request := fbmessenger.TextMessage("Hello, world!").To("USER_ID")
```

Then send your request and handle errors in sending, and errors returned from Facebook.

```go
response, err := sender.Send(request)
if err != nil {
	//Got an error. Request never got to Facebook.
} else if response.Error != nil {
	//Request got to Facebook. Facebook returned an error.
} else {
	//Hooray!
}
```

### User Profile Reference

Use type `UserProfileGetter` to get a user's profile. Create one using the page
access token for the page the user is conversing with.

```go
userProfileGetter := fbmessenger.NewUserProfileGetter("YOUR_PAGE_ACCESS_TOKEN")
```

Then request a user's profile using their userId.

```go
userProfile, err := userProfileGetter.Get(userId)
```
