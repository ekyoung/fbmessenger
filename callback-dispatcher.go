package fbmessenger

// MessageEntryHandler functions are for handling individual interactions with a user.
type MessageEntryHandler func(cb *MessagingEntry) error

/*
CallbackDispatcher routes each MessagingEntry included in a callback to an appropriate
handler for the type of entry. Note that due to webhook batching, a handler may be called
more than once per callback.
*/
type CallbackDispatcher struct {
	MessageHandler        MessageEntryHandler
	DeliveryHandler       MessageEntryHandler
	PostbackHandler       MessageEntryHandler
	AuthenticationHandler MessageEntryHandler
}

/*
Dispatch routes each MessagingEntry included in the callback to an appropriate
handler for the type of entry.
*/
func (dispatcher *CallbackDispatcher) Dispatch(cb *Callback) error {
	for _, entry := range cb.Entries {
		for _, messagingEntry := range entry.Messaging {
			if messagingEntry.Message != nil {
				if dispatcher.MessageHandler != nil {
					dispatcher.MessageHandler(messagingEntry)
				}
			} else if messagingEntry.Delivery != nil {
				if dispatcher.DeliveryHandler != nil {
					dispatcher.DeliveryHandler(messagingEntry)
				}
			} else if messagingEntry.Postback != nil {
				if dispatcher.PostbackHandler != nil {
					dispatcher.PostbackHandler(messagingEntry)
				}
			} else if messagingEntry.OptIn != nil {
				if dispatcher.AuthenticationHandler != nil {
					dispatcher.AuthenticationHandler(messagingEntry)
				}
			}
		}
	}

	return nil
}
