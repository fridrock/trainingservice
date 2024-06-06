package consumers

import (
	"context"
	"fmt"
	"testing"

	rs "github.com/fridrock/rabbitsimplier"
	"github.com/fridrock/trainingservice/db/stores"
	"github.com/fridrock/trainingservice/test"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

var (
	rmqContainer   *rabbitmq.RabbitMQContainer
	clientConsumer *test.AllMessagesConsumer
	clientProducer *rs.RProducer
	exGroupRouter  *ExGroupRouter
)

type EGSStub struct{}

func (egss EGSStub) Save(group stores.ExGroup) (int64, error) {
	return 1, nil
}
func (egss EGSStub) FindById(id int64) (stores.ExGroup, error) {
	var group stores.ExGroup
	return group, nil
}
func (egss EGSStub) FindByName(userId int64, name string) (stores.ExGroup, error) {
	var group stores.ExGroup
	return group, nil
}
func (egss EGSStub) DeleteById(id int64) error {
	return nil
}
func (egss EGSStub) DeleteByName(userId int64, name string) error {
	return nil
}
func (egss EGSStub) Update(group stores.ExGroup) error {
	return nil
}
func (egss EGSStub) UpdateByName(userId int64, name string, updated stores.ExGroup) error {
	return nil
}

func (egss EGSStub) FindByUserId(userId int64) ([]stores.ExGroup, error) {
	return nil, nil
}

func TestMain(m *testing.M) {
	//setting up
	rmqContainer = test.GetRmqContainer()
	clientProducer = test.GetClientProducer()
	clientConsumer = test.GetClientConsumer()
	exGroupRouter = &ExGroupRouter{}
	exGroupRouter.CreateConsumer(test.GetClientConfigurer())
	exGroupRouter.CreateProducer(test.GetClientConfigurer())
	egs := EGSStub{}
	exGroupRouter.SetEGS(egs)
	exGroupRouter.Setup()
	//running tests
	m.Run()
	//tearing down
	exGroupRouter.Stop()
	test.Stop()
}

func TestAddExGroup(t *testing.T) {
	clientProducer.PublishMessage(context.Background(), "sport_bot", "trainings.exgroup.create", `{"user_id":1,"name":"Back"}`)
	d := <-clientConsumer.LastMessageCh
	fmt.Println(string(d.Body))

}
