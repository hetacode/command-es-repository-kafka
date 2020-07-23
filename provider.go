package cerk

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	goeh "github.com/hetacode/go-eh"
)

// Provider interface
type Provider interface {
	FetchAllEvents(batch int) (<-chan []goeh.Event, error)
	SendEvents(events []goeh.Event) error
	Close()
}

// KafkaProvider implemented provider for kafka
type KafkaProvider struct {
	topic        string
	servers      string
	groupName    string
	consumer     *kafka.Producer
	eventsMapper *goeh.EventsMapper
}

// FetchAllEvents get all events from all partitions from specified topic
func (p *KafkaProvider) FetchAllEvents(batch int) (<-chan []goeh.Event, error) {
	rand.Seed(time.Now().UTC().UnixNano())

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":    p.servers,
		"group.id":             fmt.Sprintf("%s-replay-%d", p.groupName, rand.Int63n(10000)),
		"auto.offset.reset":    "earliest",
		"enable.partition.eof": true,
	})

	if err != nil {
		return nil, err
	}

	metadata, err := c.GetMetadata(&p.topic, false, 2000)
	if err != nil {
		return nil, err
	}

	partitionsMap := make(map[int32]bool)

	for _, partition := range metadata.Topics[p.topic].Partitions {
		partitionsMap[partition.ID] = false
	}

	c.SubscribeTopics([]string{p.topic}, nil)

	eventsChan := make(chan []goeh.Event)

	go func() {
		defer c.Close()
		defer close(eventsChan)

		run := true
		currentMessageNo := 0
		events := make([]goeh.Event, 0)
		for run == true {
			ev := c.Poll(0)
			switch e := ev.(type) {
			case *kafka.Message:
				event, err := p.eventsMapper.Resolve(string(e.Value))
				if err != nil {
					log.Printf("cannot resolve event: %s| Err: %s", string(e.Value), err)
					continue
				}

				events = append(events, event)
				currentMessageNo++

				if currentMessageNo >= batch {
					eventsChan <- events
					events = make([]goeh.Event, 0)
					currentMessageNo = 0
				}

			case kafka.PartitionEOF:
				fmt.Printf("%% Reached %v\n", e)
				partitionsMap[e.Partition] = true
				for _, e := range partitionsMap {
					run = !e
					if !e {
						break
					}
				}
			case kafka.Error:
				fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
				run = false
			}
		}

		if len(events) > 0 {
			eventsChan <- events
		}
	}()

	return eventsChan, nil
}

// SendEvents put messages on kafka topic
func (p *KafkaProvider) SendEvents(events []goeh.Event) error {
	for _, e := range events {
		e.SavePayload(e)
		message := kafka.Message{
			Key:            []byte(e.GetID()),
			TopicPartition: kafka.TopicPartition{Topic: &p.topic, Partition: kafka.PartitionAny},
			Value:          []byte(e.GetPayload()),
		}

		if err := p.consumer.Produce(&message, nil); err != nil {
			return err
		}
	}

	return nil
}

// Close just close producer instance
func (p *KafkaProvider) Close() {
	if p.consumer != nil {
		p.consumer.Close()
	}
}

// NewKafkaProvider create new instance of provider
func NewKafkaProvider(topic string, groupName string, servers string, eventsMapper *goeh.EventsMapper) Provider {
	provider := new(KafkaProvider)
	provider.servers = servers
	provider.topic = topic
	provider.groupName = groupName
	provider.eventsMapper = eventsMapper

	pr, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": provider.servers})

	if err != nil {
		panic(err)
	}

	provider.consumer = pr

	return provider
}
