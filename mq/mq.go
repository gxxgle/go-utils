package mq

// Client is a message queue client.
type Client interface {
	NewPublisher(string) (Publisher, error)
	NewSubscriber(string) (Subscriber, error)
	Purge(string) error
	Close()
}

// Publisher can send message.
type Publisher interface {
	Publish(string, []byte) error
}

// Subscriber can receive message.
type Subscriber interface {
	Subscribe(func([]byte) error)
}

type Message struct {
	Key  string
	Body []byte
}
