package mq

// import (
// 	"sync"

// 	"github.com/Shopify/sarama"
// 	"github.com/bsm/sarama-cluster"
// 	"github.com/gxxgle/go-utils/log"
// )

// type KafkaConfig struct {
// 	Servers       []string `json:"servers"`
// 	ConsumerGroup string   `json:"consumer_group"`
// }

// type kafka struct {
// 	cf      *KafkaConfig
// 	cfg     *cluster.Config
// 	stopped bool
// 	exit    chan bool
// 	wg      sync.WaitGroup

// 	consumers []*cluster.Consumer
// }

// func NewKafka(cf *KafkaConfig) (Client, error) {
// 	cfg := cluster.NewConfig()
// 	cfg.Net.SASL.Enable = false
// 	cfg.Net.TLS.Enable = false
// 	cfg.Producer.Return.Successes = true
// 	cfg.Consumer.Return.Errors = true
// 	cfg.Consumer.Offsets.Initial = sarama.OffsetOldest
// 	cfg.Version = sarama.V0_10_0_0

// 	if err := cfg.Validate(); err != nil {
// 		return nil, err
// 	}

// 	return &kafka{
// 		cf:   cf,
// 		cfg:  cfg,
// 		exit: make(chan bool, 0),
// 	}, nil
// }

// func (k *kafka) Publish(topic string, key string, body []byte) error {
// 	producer, err := sarama.NewSyncProducer(k.cf.Servers, &k.cfg.Config)
// 	if err != nil {
// 		return err
// 	}

// 	msg := &sarama.ProducerMessage{
// 		Topic: topic,
// 		Key:   sarama.StringEncoder(key),
// 		Value: sarama.ByteEncoder(body),
// 	}

// 	_, _, err = producer.SendMessage(msg)
// 	return err
// }

// func (k *kafka) Subscribe(topic string, handler func([]byte) error) error {
// 	consumer, err := cluster.NewConsumer(k.cf.Servers, k.cf.ConsumerGroup, []string{topic}, k.cfg)
// 	if err != nil {
// 		return err
// 	}

// 	k.consumers = append(k.consumers, consumer)
// 	k.wg.Add(1)

// 	go func() {
// 		for !k.stopped {
// 			select {
// 			case msg, ok := <-consumer.Messages():
// 				if !ok || msg == nil {
// 					continue
// 				}

// 				if err := handler(msg.Value); err != nil {
// 					log.Errorw("mq kafka subscribe handler message error",
// 						"queue", string(msg.Key), "body", string(msg.Value),
// 						"err", err)
// 					continue
// 				}

// 				consumer.MarkOffset(msg, "")

// 			case err, ok := <-consumer.Errors():
// 				if !ok || err == nil {
// 					continue
// 				}

// 				log.Errorw("mq kafka subscribe error", "err", err)

// 			case <-k.exit:
// 			}
// 		}

// 		k.wg.Done()
// 		log.Infow("mq kafka subscribe stopped")
// 	}()

// 	return nil
// }

// func (k *kafka) Close() {
// 	k.stopped = true
// 	close(k.exit)
// 	k.wg.Wait()

// 	for _, c := range k.consumers {
// 		c.Close()
// 	}
// }
