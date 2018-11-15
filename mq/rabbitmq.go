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

func newRabbitmqPublisher(cli *rabbitmq, exchange string) *rabbitmqPublisher {
	out := &rabbitmqPublisher{
		cli:      cli,
		puber:    cony.NewPublisher(exchange, ""),
		exchange: exchange,
		msgs:     make(chan *Message, 0),
	}

	out.run()
	out.cli.Publish(out.puber)
	return out
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

func (c *rabbitmq) NewPublisher(exchange string) (Publisher, error) {
	return newRabbitmqPublisher(c, exchange), nil
}

func (c *rabbitmq) NewSubscriber(queue string) (Subscriber, error) {
	return newRabbitmqSubscriber(c, queue), nil
}

func (c *rabbitmq) Purge(queue string) error {
	purge := func(der cony.Declarer) error {
		ch, ok := der.(*amqp.Channel)
		if !ok || ch == nil {
			return nil
		}

		_, err := ch.QueuePurge(queue, true)
		return err
	}

	c.Declare([]cony.Declaration{purge})
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
					log.Errorw("mq rabbitmq client error", "err", err)
				}

			case <-c.exit:
			}
		}

		c.wg.Done()
	}()
}

func (p *rabbitmqPublisher) send(msg *Message) {
	err := p.puber.PublishWithRoutingKey(amqp.Publishing{Body: msg.Body}, msg.Key)
	if err != nil {
		log.Errorw("mq rabbitmq publisher send message error", "exchange",
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
		log.Infow("mq rabbitmq publisher stopped", "exchange", p.exchange)
	}()
}

func (p *rabbitmqPublisher) Publish(key string, body []byte) error {
	if p.cli.stopped {
		return cony.ErrPublisherDead
	}

	p.msgs <- &Message{Key: key, Body: body}
	return nil
}

func (s *rabbitmqSubscriber) Subscribe(handler func([]byte) error) {
	s.cli.wg.Add(1)

	go func() {
		for !s.cli.stopped {
			select {
			case msg := <-s.coner.Deliveries():
				if err := handler(msg.Body); err != nil {
					log.Errorw("mq rabbitmq subscriber handler message error",
						"queue", s.queue, "body", string(msg.Body), "err", err)
					continue
				}
				msg.Ack(false)

			case err := <-s.coner.Errors():
				if err != nil {
					log.Errorw("mq rabbitmq subscriber error", "queue", s.queue, "err", err)
				}

			case <-s.cli.exit:
			}
		}

		s.cli.wg.Done()
		log.Infow("mq rabbitmq subscriber stopped", "queue", s.queue)
	}()
}
