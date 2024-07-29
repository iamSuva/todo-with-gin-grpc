package handlers

import (
	grpclient "todowithgin/grpcClient"
	"todowithgin/service"
)

type TaskHandler struct {
	Service     service.TaskService
	GrpcService grpclient.GrpcInterface
}

func NewTaskHandler(t service.TaskService, g grpclient.GrpcInterface) *TaskHandler {
	return &TaskHandler{Service: t, GrpcService: g}
}
