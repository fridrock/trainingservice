package consumers

import (
	"context"
	"database/sql"
	"strings"
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
	data := []struct {
		testName       string
		message        string
		resultExpected string
		errMessage     string
	}{
		{
			"Negative case wrong input",
			`{"user_id":1`,
			wrongInput,
			"error while validation, received: %s",
		},
		{
			"Positive case",
			`{"user_id":1,"name":"Back"}`,
			"SUCCESS: id:1",
			"wrong result adding exgroup, received: %s",
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			clientProducer.PublishMessage(context.Background(),
				"sport_bot",
				"trainings.exgroup.create",
				d.message,
			)
			received := string((<-clientConsumer.LastMessageCh).Body)
			if received != d.resultExpected {
				t.Errorf(d.errMessage, received)
			}
		})
	}
}

func TestDeleteByName(t *testing.T) {
	data := []struct {
		testName       string
		message        string
		resultExpected string
		errorMessage   string
	}{
		{"Negative case wrong input",
			`{"user_id:2}`,
			wrongInput,
			"error with wrong input, received: %s"},
		{"Negative case no such ex group",
			`{"user_id":2, "name":"Unexisting"}`,
			notDeleted,
			"error with deleting unexisting ex group, received: %s"},
		{"Positive case",
			`{"user_id":2, "name":"Back"}`,
			success,
			"error with successful deletion of ex group, received: %s"},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			clientProducer.PublishMessage(
				context.Background(),
				"sport_bot",
				"trainings.exgroup.delete",
				d.message,
			)
			received := <-clientConsumer.LastMessageCh
			if string(received.Body) != d.resultExpected {
				t.Errorf(d.errorMessage, string(received.Body))
			}
		})
	}
}

func TestFindByName(t *testing.T) {
	data := []struct {
		testName       string
		message        string
		resultExpected string
		errMessage     string
	}{
		{
			"Negative case wrong input",
			`{"user_id":2}`,
			wrongInput,
			"error with wrong input, received: %s",
		},
		{
			"Negative case no such ex group",
			`{"user_id":2, "name":"Unexisting"}`,
			"ERROR: " + sql.ErrNoRows.Error(),
			"error with finding unexisting ex group, received: %s",
		},
		{
			"Positive case found",
			`{"user_id":2, "name":"Back"}`,
			successFinding,
			"error with successful finding of ex group, received: %s",
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			clientProducer.PublishMessage(
				context.Background(),
				"sport_bot",
				"trainings.exgroup.find",
				d.message,
			)
			received := string((<-clientConsumer.LastMessageCh).Body)
			if d.testName != "Positive case found" && received != d.resultExpected {
				t.Errorf(d.errMessage, received)
			}
			if d.testName == "Positive case found" {
				parts := strings.Split(received, ":")
				if parts[0] != success {
					t.Errorf("error getting right result:%v", parts[0])
				}
			}
		})
	}
}

func TestUpdateByName(t *testing.T) {
	data := []struct {
		testName       string
		message        string
		resultExpected string
		errMessage     string
	}{
		{
			"Negative case wrong input",
			`{"user_id":2}`,
			wrongInput,
			"error with wrong input, received: %s",
		},
		{
			"Negative case no such ex group",
			`{"user_id":2,"name":"Unexisting","newname":"NewBack"}`,
			"ERROR: " + stores.NotUpdated.Error(),
			"error with updating unexisting ex group, received: %s",
		},
		{
			"Positive case updated",
			`{"user_id":2, "name":"Back", "newname":"NewBack"}`,
			success,
			"error with successful updating of ex group, received: %s",
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			clientProducer.PublishMessage(
				context.Background(),
				"sport_bot",
				"trainings.exgroup.update",
				d.message,
			)
			received := string((<-clientConsumer.LastMessageCh).Body)
			if received != d.resultExpected {
				t.Errorf(d.errMessage, received)
			}
		})
	}
}

func TestFindByUser(t *testing.T) {
	data := []struct {
		testName       string
		message        string
		resultExpected string
		errMessage     string
	}{
		{
			"Negative case wrong input",
			`{"userid":2}`,
			wrongInput,
			"error with wrong input, received: %s",
		},
		{
			"Negative case no such user",
			`{"user_id":0}`,
			"ERROR: " + sql.ErrNoRows.Error(),
			"error with updating unexisting user, received: %s",
		},
		{
			"Positive case got groups",
			`{"user_id":2}`,
			success,
			"error with successful getting ex groups, received: %s",
		},
	}

	for _, d := range data {
		t.Run(d.testName, func(t *testing.T) {
			clientProducer.PublishMessage(
				context.Background(),
				"sport_bot",
				"trainings.exgroup.findByUser",
				d.message,
			)
			received := string((<-clientConsumer.LastMessageCh).Body)
			if received != d.resultExpected {
				t.Errorf(d.errMessage, received)
			}
		})
	}
}
