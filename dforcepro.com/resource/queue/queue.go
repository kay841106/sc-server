package queue

import (
	"dforcepro.com/resource"
	machinery "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/tasks"
)

var TaskServer *machinery.Server

func GetTaskServer(di *resource.Di) *machinery.Server {
	if TaskServer != nil {
		return TaskServer
	}
	taskServer, err := machinery.NewServer(&di.Rabbitmq)
	if err != nil {
		panic(err)
	}
	TaskServer = taskServer
	return TaskServer
}

func SendTask(taskName string, args []tasks.Arg) (*tasks.TaskState, error) {
	signature := &tasks.Signature{
		Name: taskName,
		Args: args,
	}
	asyncResult, err := TaskServer.SendTask(signature)
	if err != nil {
		return nil, err
	}
	return asyncResult.GetState(), nil
}
