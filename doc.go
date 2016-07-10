/*
Package fbmessenger is a library for making requests to and handling callbacks from
the Facebook Messenger Platform API.

Key Features

	* Fluent API makes building messages to send easy.
	* Timeoutable, cancellable requests using `context.Context`.
	* Designed for use with one or many subscribed pages.

Quick Start

The primary types in the package are CallbackDispatcher and Client. CallbackDispatcher
is used to handle the callbacks Facebook sends to your webhook endpoint. Client is used to send
messages and to get user profiles.

CallbackDispatcher Usage

	// Unmarshal the json received at your webhook endpoint into an instance of type Callback.

	cb := &fbmessenger.Callback{}
	err := json.Unmarshal(requestBytes, cb)

	// Use type CallbackDispatcher to route each MessageEntry included in the callback to an
	// appropriate handler for the type of entry. Note that due to webhook batching, a
	// handler may be called more than once per callback.

	dispatcher := &fbmessenger.CallbackDispatcher{
		MessageHandler: MessageReceived
	}

	err := dispatcher.Dispatch(cb)

	// Callback handlers should have a signature mathing the MessageEntryHandler type.

	func MessageReceived(cb *fbmessenger.MessagingEntry) error {
		//Do stuff
	}

Client Usage

	// Create a `Client` to make requests to the messenger API.

	client := fbmessenger.Client{}

	// There are structs for the different types of messages you can send. The easiest way to
	// create them is with the fluent API.

	request := fbmessenger.TextMessage("Hello, world!").To("USER_ID")

	// Then send your request and handle errors in sending, and errors returned from Facebook.

	response, err := client.Send(request, "YOUR_PAGE_ACCESS_TOKEN")
	if err != nil {
		//Got an error. Request never got to Facebook.
	} else if response.Error != nil {
		//Request got to Facebook. Facebook returned an error.
	} else {
		//Hooray!
	}

	// Get a user's profile using their userId.

	userProfile, err := client.GetUserProfile("USER_ID", "YOUR_PAGE_ACCESS_TOKEN")

*/
package fbmessenger
