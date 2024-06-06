package consumers

import (
	"context"
	"database/sql"
	"log/slog"
	"testing"

	rs "github.com/fridrock/rabbitsimplier"
	"github.com/fridrock/trainingservice/test"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

var (
	rmqContainer   *rabbitmq.RabbitMQContainer
	clientConsumer *test.AllMessagesConsumer
	clientProducer *rs.RProducer
	exGroupRouter  *ExGroupRouter
)

const (
	wrongInput     = "ERROR: wrong input"
	notDeleted     = "ERROR: no rows deleted"
	success        = "SUCCESS"
	successFinding = `SUCCESS: {"id":1,"user_id":2,"name":"Back"}`
)

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
	clientProducer.PublishMessage(context.Background(),
		"sport_bot",
		"trainings.exgroup.create",
		`{"user_id":1,"name":"Back"}`)
	d := <-clientConsumer.LastMessageCh
	if string(d.Body) != `SUCCESS: id:1` {
		t.Error("wrong result adding exgroup")
	}
	slog.Info(string(d.Body))
}

func TestDeleteByName(t *testing.T) {
	//negative case
	clientProducer.PublishMessage(
		context.Background(),
		"sport_bot",
		"trainings.exgroup.delete",
		`{"user_id":2}`,
	)
	d := <-clientConsumer.LastMessageCh
	if string(d.Body) != wrongInput {
		t.Errorf("error with wrong input, received: %s", string(d.Body))
	}

	//negative case, searching for unexisting ex group
	clientProducer.PublishMessage(
		context.Background(),
		"sport_bot",
		"trainings.exgroup.delete",
		`{"user_id":2, "name":"Unexisting"}`,
	)
	d = <-clientConsumer.LastMessageCh
	slog.Info(string(d.Body))
	if string(d.Body) != notDeleted {
		t.Errorf("error with deleting unexisting ex group, received: %s", string(d.Body))
	}
	//positive case
	clientProducer.PublishMessage(
		context.Background(),
		"sport_bot",
		"trainings.exgroup.delete",
		`{"user_id":2, "name":"Back"}`,
	)
	d = <-clientConsumer.LastMessageCh
	if string(d.Body) != success {
		t.Errorf("error with successful deletion of ex group, received: %s", string(d.Body))
	}
}

func TestFindByName(t *testing.T) {
	//negative case: wrong input
	clientProducer.PublishMessage(
		context.Background(),
		"sport_bot",
		"trainings.exgroup.find",
		`{"user_id":2}`,
	)
	d := <-clientConsumer.LastMessageCh
	slog.Info(string(d.Body))
	if string(d.Body) != wrongInput {
		t.Errorf("error with wrong input, received: %s", string(d.Body))
	}
	//negative case: unexisting
	clientProducer.PublishMessage(
		context.Background(),
		"sport_bot",
		"trainings.exgroup.find",
		`{"user_id":2, "name":"Unexisting"}`,
	)
	d = <-clientConsumer.LastMessageCh
	slog.Info(string(d.Body))
	if string(d.Body) != "ERROR: "+sql.ErrNoRows.Error() {
		t.Errorf("error with finding unexisting ex group, received: %s", string(d.Body))
	}
	//positive case
	clientProducer.PublishMessage(
		context.Background(),
		"sport_bot",
		"trainings.exgroup.find",
		`{"user_id":2, "name":"Back"}`,
	)
	d = <-clientConsumer.LastMessageCh
	slog.Info(string(d.Body))
	if string(d.Body) != successFinding {
		t.Errorf("error with successful finding of ex group, received: %s", string(d.Body))
	}
}
