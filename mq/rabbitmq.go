package mq

import (
	"sync"

	"github.com/assembla/cony"
	"github.com/gxxgle/go-utils/log"
	"github.com/streadway/amqp"
)

type rabbitmq struct {
	*cony.Client
	stopped bool
	exit    chan bool
	wg      sync.WaitGroup
}

type rabbitmqPublisher struct {
	cli      *rabbitmq
	puber    *cony.Publisher
	exchange string
	msgs     chan *Message
}

type rabbitmqSubscriber struct {
	cli   *rabbitmq
	coner *cony.Consumer
	queue string
}

func NewRabbitMQ(url string) Client {
	out := &rabbitmq{
		Client: cony.NewClient(
			cony.URL(url),
			cony.Backoff(cony.DefaultBackoff),
		),
		exit: make(chan bool, 0),
	}

	out.run()
	return out
}

func (c *rabbitmq) Publish(exchange string, key string, body []byte) error {
	return newRabbitmqPublisher(c, exchange).Publish(key, body)
}

func (c *rabbitmq) Subscribe(queue string, handler func([]byte) error) error {
	newRabbitmqSubscriber(c, queue).Subscribe(handler)
	return nil
}

func (c *rabbitmq) Close() {
	c.stopped = true
	close(c.exit)
	c.Client.Close()
	c.wg.Wait()
}

func (c *rabbitmq) run() {
	c.wg.Add(1)

	go func() {
		for !c.stopped && c.Loop() {
			select {
			case err := <-c.Errors():
				if err != nil {
					log.Errorw("go-utils mq rabbitmq client error", "err", err)
				}

			case <-c.exit:
			}
		}

		c.wg.Done()
	}()
}

func newRabbitmqPublisher(cli *rabbitmq, exchange string) *rabbitmqPublisher {
	out := &rabbitmqPublisher{
		cli:      cli,
		puber:    cony.NewPublisher(exchange, ""),
		exchange: exchange,
		msgs:     make(chan *Message, 0),
	}

	out.run()
	out.cli.Client.Publish(out.puber)
	return out
}

func (p *rabbitmqPublisher) send(msg *Message) {
	err := p.puber.PublishWithRoutingKey(amqp.Publishing{Body: msg.Body}, msg.Key)
	if err != nil {
		log.Errorw("go-utils mq rabbitmq publisher send message error", "exchange",
			p.exchange, "key", msg.Key, "body", string(msg.Body))
	}
}

func (p *rabbitmqPublisher) run() {
	p.cli.wg.Add(1)

	go func() {
		for !p.cli.stopped || len(p.msgs) > 0 {
			select {
			case msg := <-p.msgs:
				p.send(msg)

			case <-p.cli.exit:
			}
		}

		close(p.msgs)
		p.cli.wg.Done()
		log.Infow("go-utils mq rabbitmq publisher stopped", "exchange", p.exchange)
	}()
}

func (p *rabbitmqPublisher) Publish(key string, body []byte) error {
	if p.cli.stopped {
		return cony.ErrPublisherDead
	}

	p.msgs <- &Message{Key: key, Body: body}
	return nil
}

func newRabbitmqSubscriber(cli *rabbitmq, queue string) *rabbitmqSubscriber {
	out := &rabbitmqSubscriber{
		cli:   cli,
		coner: cony.NewConsumer(&cony.Queue{Name: queue}),
		queue: queue,
	}

	out.cli.Consume(out.coner)
	return out
}

func (s *rabbitmqSubscriber) Subscribe(handler func([]byte) error) {
	s.cli.wg.Add(1)

	go func() {
		for !s.cli.stopped {
			select {
			case msg := <-s.coner.Deliveries():
				if err := handler(msg.Body); err != nil {
					log.Errorw("go-utils mq rabbitmq subscriber handler message error",
						"queue", s.queue, "body", string(msg.Body), "err", err)
					continue
				}
				msg.Ack(false)

			case err := <-s.coner.Errors():
				if err != nil {
					log.Errorw("go-utils mq rabbitmq subscriber error", "queue", s.queue, "err", err)
				}

			case <-s.cli.exit:
			}
		}

		s.cli.wg.Done()
		log.Infow("go-utils mq rabbitmq subscriber stopped", "queue", s.queue)
	}()
}
