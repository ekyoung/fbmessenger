package fbmessenger

type MessageEntryHandler func(cb *MessagingEntry) error

type CallbackDispatcher struct {
	MessageHandler        MessageEntryHandler
	DeliveryHandler       MessageEntryHandler
	PostbackHandler       MessageEntryHandler
	AuthenticationHandler MessageEntryHandler
}

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
