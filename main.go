package main

import (
	"log"

	rs "github.com/fridrock/rabbitsimplier"
	"github.com/fridrock/trainingservice/api/routers"
)

// defining answer event
type ExGroupCreatedEvent struct {
	Event string `json:"event"`
	Text  string `json:"text"`
}

func setupConfigurer() *rs.RConfigurer {
	brc := &rs.RConfigurer{}
	err := brc.Configure(rs.Config{
		Username: "guest",
		Password: "guest",
		Host:     "localhost",
	})
	if err != nil {
		log.Fatal("error creating connection to rabbitmq")
	}
	return brc
}

func main() {
	//setting up configurer
	brc := setupConfigurer()
	defer brc.Stop()
	//setting up routers
	exgroupRouter := routers.NewExGroupRouter(brc)
	exgroupRouter.Setup()
	defer exgroupRouter.Stop()
	trainingsRouter := routers.NewTrainingRouter(brc)
	trainingsRouter.Setup()
	defer trainingsRouter.Stop()
	//infinite work of service
	var forever chan struct{}
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
