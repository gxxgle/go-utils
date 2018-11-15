package mq

// Client is a message queue client.
type Client interface {
	Publish(exchange string, key string, body []byte) error
	Subscribe(queue string, handler func([]byte) error) error
	Close()
}

type Message struct {
	Key  string
	Body []byte
}
