# Package `fbmessenger`

[![CircleCI](https://circleci.com/gh/ekyoung/fbmessenger.svg?style=svg)](https://circleci.com/gh/ekyoung/fbmessenger)

Go (golang) package for writing bots on the [Facebook Messenger Platform](https://developers.facebook.com/docs/messenger-platform).

## Key Features

* Fluent API makes building messages easy.
* Timeoutable, cancellable requests using `context.Context`.
* Designed for use with one or many subscribed pages.

## Installation

```bash
go get gopkg.in/ekyoung/fbmessenger.v1
```

## Quick Start

The primary types in the package are `CallbackDispatcher` and `Client`. `CallbackDispatcher`
is used to handle the callbacks Facebook sends to your webhook endpoint. `Client` is used to send
messages and to get user profiles.

### CallbackDispatcher

Unmarshal the json received at your webhook endpoint into an instance of type `Callback`.

```go
cb := &fbmessenger.Callback{}
err := json.Unmarshal(requestBytes, cb)
```

Use type `CallbackDispatcher` to route each `MessagingEntry` included in the callback to an appropriate
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

### Client

Create a `Client` to make requests to the messenger API.

```go
client := fbmessenger.Client{}
```

There are structs for the different types of messages you can send. The easiest way to create them
is with the fluent API.

```go
request := fbmessenger.TextMessage("Hello, world!").To("USER_ID")
```

Then send your request and handle errors in sending, and errors returned from Facebook.

```go
response, err := client.Send(request, "YOUR_PAGE_ACCESS_TOKEN")
if err != nil {
	//Got an error. Request never got to Facebook.
} else if response.Error != nil {
	//Request got to Facebook. Facebook returned an error.
} else {
	//Hooray!
}
```

Get a user's profile using their userId.

```go
userProfile, err := client.GetUserProfile("USER_ID", "YOUR_PAGE_ACCESS_TOKEN")
```

For more control over requests (timeouts, etc.) use the `*WithContext` version of the
above methods.

```go
ctx, _ := context.WithTimeout(context.Background(), 500*time.Millisecond)
response, err := client.SendWithContext(ctx, request, "YOUR_PAGE_ACCESS_TOKEN")
userProfile, err := userProfileGetter.GetUserProfileWithContext(ctx, "USER_ID", "YOUR_PAGE_ACCESS_TOKEN")
```

## Inspiration

Some ideas where pulled from [Go Client Library Best Practices](https://medium.com/@cep21/go-client-library-best-practices-83d877d604ca) by Jack Lindamood.