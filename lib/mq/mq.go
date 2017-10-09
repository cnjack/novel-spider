package mq

import (
	"encoding/json"
	"github.com/nsqio/go-nsq"
	"math/rand"
	"sync"
	"log"
)

type Broker struct {
	Addrs []string
	sync.Mutex
	running bool
	p       []*nsq.Producer
	c       []*subscriber
}

type subscriber struct {
	topic string
	c     *nsq.Consumer
	// handler so we can resubcribe
	h       nsq.HandlerFunc
	channel string
	// concurrency
	n int
}
type Handler func(Publication) error

type Message struct {
	Header map[string]string
	Body   []byte
}

// Publication is given to a subscription handler for processing
type Publication interface {
	Topic() string
	Message() *Message
	Ack() error
}

func (p *publication) Topic() string {
	return p.topic
}

func (p *publication) Message() *Message {
	return p.m
}

func (p *publication) Ack() error {
	p.nm.Finish()
	return nil
}

type publication struct {
	topic string
	m     *Message
	nm    *nsq.Message
}

func (n *Broker) Connect() error {
	n.Lock()
	defer n.Unlock()
	if n.running {
		return nil
	}
	config := nsq.NewConfig()
	var producers []*nsq.Producer
	for _, addr := range n.Addrs {

		p, err := nsq.NewProducer(addr, config)
		if err != nil {
			return err
		}
		producers = append(producers, p)
	}
	n.p = producers
	n.running = true
	return nil
}
func (n *Broker) Disconnect() error {
	n.Lock()
	defer n.Unlock()

	if !n.running {
		return nil
	}

	// stop the producers
	for _, p := range n.p {
		p.Stop()
	}

	// stop the consumers
	for _, c := range n.c {
		c.c.Stop()

		// disconnect from all nsq brokers
		for _, addr := range n.Addrs {
			c.c.DisconnectFromNSQD(addr)
		}
	}

	n.p = nil
	n.running = false
	return nil
}

func (n *Broker) Publish(topic string, message *Message) error {
	p := n.p[rand.Int()%len(n.p)]
	b, err := json.Marshal(message)
	if err != nil {
		return err
	}
	return p.Publish(topic, b)
}

func (n *Broker) Subscribe(topic, channel string, handler Handler) (*subscriber, error) {
	config := nsq.NewConfig()
	c, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Println("NewConsumer err: ",err)
		return nil, err
	}

	h := nsq.HandlerFunc(func(nm *nsq.Message) error {
		nm.DisableAutoResponse()
		var m *Message
		if err := json.Unmarshal(nm.Body, &m); err != nil {
			return err
		}
		return handler(&publication{
			topic: topic,
			m:     m,
			nm:    nm,
		})
	})
	c.AddHandler(h)
	if err := c.ConnectToNSQDs(n.Addrs); err != nil {
		log.Println("ConnectToNSQDs err: ",err)
		return nil, err
	}
	sub := &subscriber{
		topic:   topic,
		c:       c,
		h:       h,
		channel: channel,
	}
	n.c = append(n.c, sub)
	return sub, nil
}

func (s *subscriber) Topic() string {
	return s.topic
}

func (s *subscriber) Unsubscribe() error {
	s.c.Stop()
	return nil
}
