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

The struct type `Callback` can be created by unmarshaling the json received at your webhook endpoint.

```go
cb := &fbmessenger.Callback{}
err = json.Unmarshal(requestBytes, cb)
```

The interface type `CallbackDispatcher` examines callbacks and routes each `MessageEntry` to handlers you
register for each type of message. Pass your unmarshalled callback struct to `Dispatch` to feed it
to your handlers. Note that due to webhook batching, a handler may be called more than once per callback.

```go
dispatcher := fbmessenger.NewCallbackDispatcher()
dispatcher.OnMessageReceived(MessageReceived)

err := dispatcher.Dispatch(cb)
if err != nil {
	//Error handling
}
```

Callback handlers should have a signature mathing the `MessageEntryHandler` type.

```go
func MessageReceived(cb *fbmessenger.MessagingEntry) error {
	//Do stuff
}
```

### Send API Reference

The interface type `Sender` handles sending to the messenger API. Create one using the page access
token for the page you want to send as.

```go
sendApi := fbmessenger.NewSendApi("YOUR_PAGE_ACCESS_TOKEN")
```

There are structs for the different types of messages you can send. The easiest way to create them
is with the fluent API.

```go
request := fbmessenger.TextMessage("Hello, world!").To("USER_ID")
```

Then send your request and handle errors in sending, and errors returned from Facebook.

```go
response, err := sendApi.Send(request)
if err != nil {
	//Got an error. Request never got to Facebook.
} else if sendResponse.Error != nil {
	//Request got to Facebook. Facebook returned an error.
} else {
	//Hooray!
}
```