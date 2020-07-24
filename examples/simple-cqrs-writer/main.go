package main

import (
	"io/ioutil"
	"log"
	"net/http"

	cerk "github.com/hetacode/command-es-repository-kafka"
	examplecommandhandlers "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/command_handlers"
	examplecommands "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/commands"
	examplerepository "github.com/hetacode/command-es-repository-kafka/examples/simple-cqrs-writer/repository"
	goeh "github.com/hetacode/go-eh"
)

type MainContainer struct {
	EventsMapper         *goeh.EventsMapper
	EventsHandlerManager *goeh.EventsHandlerManager
}

func (h *MainContainer) handler(res http.ResponseWriter, req *http.Request) {
	if req.Method == "POST" {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			log.Panic(err)
			res.WriteHeader(400)
		}
		log.Printf("Received commmand: %s", string(body))

		event, err := h.EventsMapper.Resolve(string(body))
		if err != nil {
			log.Panic(err)
			res.WriteHeader(400)
		}
		err = h.EventsHandlerManager.Execute(event)
		if err != nil {
			log.Panic(err)
			res.WriteHeader(400)
		}
	}
}

func main() {
	eventsMapper := new(goeh.EventsMapper)
	eventsMapper.Register(new(examplecommands.CreateUserCommand))
	eventsMapper.Register(new(examplecommands.UpdateUserCommand))

	eventManager := new(goeh.EventsHandlerManager)

	provider := cerk.NewKafkaProvider("example-topic", "example-app-group", "192.168.1.151:9092", eventsMapper)
	defer provider.Close()

	repository := new(examplerepository.UsersRepository)
	repository.MemoryRepository = cerk.NewMemoryRepository()
	if err := repository.InitProvider(provider, repository); err != nil {
		panic(err)
	}

	eventManager.Register(new(examplecommands.CreateUserCommand), &examplecommandhandlers.CreateUserCommandHandler{Repository: repository})
	eventManager.Register(new(examplecommands.UpdateUserCommand), &examplecommandhandlers.UpdateUserCommandHandler{Repository: repository})

	h := &MainContainer{
		EventsMapper:         eventsMapper,
		EventsHandlerManager: eventManager,
	}

	http.HandleFunc("/", h.handler)
	http.ListenAndServe(":4000", nil)
}
