package main

import (
	"flag"
	"fmt"

	cerk "github.com/hetacode/command-es-repository-kafka"
	goeh "github.com/hetacode/go-eh"
)

func fetch(server string, topic string, group string) {
	eventsMapper := new(goeh.EventsMapper)

	provider := cerk.NewKafkaProvider(topic, group, server, eventsMapper)
	eventsBatch, err := provider.FetchAllEvents(10)

	if err != nil {
		panic(err)
	}

	for events := range eventsBatch {
		for _, event := range events {
			fmt.Printf("%s: %s\n", event.GetID(), event.GetPayload())
		}
	}
}

func fill(server string, topic string, group string) {
	eventsMapper := new(goeh.EventsMapper)

	events := make([]goeh.Event, 0)
	provider := cerk.NewKafkaProvider(topic, group, server, eventsMapper)

	i := 0
	for i < 100 {
		event := &goeh.EventData{
			Payload: fmt.Sprintf(`{"type": "GenericEvent", "createTime":"2009-11-10T23:00:00Z", "version":1, "id": "%d"}`, i),
		}
		events = append(events, event)

		i++
	}
	if err := provider.SendEvents(events); err != nil {
		panic(err)
	}
}

func main() {
	action := flag.String("action", "fill", "fill | fetch events")
	topic := flag.String("topic", "test", "a kafka topic")
	server := flag.String("server", "localhost:9092", "a kafka server")
	group := flag.String("group", "test-group", "a kafka consumer group name")

	flag.Parse()

	switch *action {
	case "fetch":
		fetch(*server, *topic, *group)
	case "fill":
		fill(*server, *topic, *group)
	}

}
