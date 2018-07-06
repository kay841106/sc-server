package queue

import (
	"fmt"
	"testing"

	"dforcepro.com/resource"
	"github.com/RichardKnop/machinery/v1/config"
	"github.com/stretchr/testify/assert"
)

func WorkerFuc() (bool, error) {
	return false, nil
}
func GetDi() *resource.Di {
	var cnf = config.Config{
		Broker:        "amqp://rd:1234@dforcepro-db:5672/",
		DefaultQueue:  "ytz_tasks",
		ResultBackend: "redis://127.0.0.1:6379",
		AMQP: &config.AMQPConfig{
			Exchange:     "machinery_exchange",
			ExchangeType: "direct",
			BindingKey:   "ytz_tasks",
		},
	}
	di := resource.Di{Rabbitmq: cnf}
	return &di
}

func Test_GetTaskServer(t *testing.T) {
	server1 := GetTaskServer(GetDi())
	err := server1.RegisterTask("aaa", WorkerFuc)
	if err != nil {
		fmt.Println(err.Error())
	}
	server2 := TaskServer

	fmt.Println(TaskServer.GetRegisteredTaskNames())
	assert.Equal(t, server1, server2, "should be equal.")
}
