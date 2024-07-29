package mock

import (
	"fmt"
	"time"
	grpclient "todowithgin/grpcClient"
	"todowithgin/models"
	"todowithgin/service"
)

type MockStatusCode int
type MockService struct {
	MockErr MockStatusCode
}

type MockGrpcService struct {
	GrpcErr MockStatusCode
}

const (
	ErrInternalServer MockStatusCode = iota
	ErrNotFound
	ErrBadRequest
	ErrConflict
	ErrInvalidId
	ErrInvalidInput
	ErrPassword
	ErrFailedToGenerate
	ErrInvalidToken
	ErrUnauthorized
	Ok
)

func (m *MockService) GetTask(id int) (models.Task, error) {
	if m.MockErr == ErrNotFound {
		return models.Task{}, service.ErrorNoRows
	} else if m.MockErr == ErrInternalServer {
		return models.Task{}, service.ErrorInternalserver
	}
	staticTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
	return models.Task{
		Id:            1,
		Title:         "Task 1",
		Description:   "this is task 1",
		IsCompleted:   false,
		CreatedAt_UTC: staticTime,
		UpdatedAt_UTC: staticTime,
	}, nil
}
func (m *MockService) GetTasks() ([]models.Task, error) {
	if m.MockErr == ErrInternalServer {
		return nil, service.ErrorInternalserver
	}
	return []models.Task{}, nil
}
func (m *MockService) CreateTask(models.Task) error {
	if m.MockErr == ErrInternalServer {
		return service.ErrorInternalserver
	} else if m.MockErr == ErrConflict {
		return service.ErrorUniqueTitle
	} else if m.MockErr == ErrUnauthorized {
		return service.ErrorUnauthorized
	} else if m.MockErr == ErrInvalidId {
		return service.ErrorInvalidUserId
	}
	return nil
}
func (m *MockService) UpdateTask(task models.Task, id int, userid int) error {
	if m.MockErr == ErrConflict {
		return service.ErrorUniqueTitle
	} else if m.MockErr == ErrInternalServer {
		return service.ErrorInternalserver
	} else if m.MockErr == ErrNotFound {
		return service.ErrorNoRows
	}
	return nil
}
func (m *MockService) DeleteTask(id int, userid int) error {
	fmt.Println("mock delete")
	if m.MockErr == ErrInternalServer {
		return service.ErrorInternalserver
	}
	if m.MockErr == ErrNotFound {
		return service.ErrorNoRows
	}
	if m.MockErr == ErrUnauthorized {
		return service.ErrorUnauthorized
	}
	return nil

}

func (m *MockService) SignUpUser(models.User) (int, error) {
	fmt.Println("sign up mock")
	if m.MockErr == ErrConflict {
		return 0, service.ErrorUniqueUsername
	}
	if m.MockErr == ErrInternalServer {
		return 0, service.ErrorInternalserver
	}

	return 1, nil
}

func (m *MockService) LoginUser(models.User) (models.User, error) {
	fmt.Print("login mock")
	if m.MockErr == ErrInternalServer {
		return models.User{}, service.ErrorInternalserver
	} else if m.MockErr == ErrPassword {
		return models.User{}, service.ErrorPasswordNotMatch
	} else if m.MockErr == ErrNotFound {
		return models.User{}, service.ErrorNoRows
	}

	return models.User{
		UserId:   1,
		Username: "suvadip",
		Password: "password",
	}, nil

}

//grpc mock

func (m *MockGrpcService) GetTokenHandler(username string, id int) (string, error) {
	fmt.Print("mock grpc called")
	token := "h34h34hjg3t33bjnfjgj45kjgj333hjfkj4jnj54nn"

	if m.GrpcErr == ErrFailedToGenerate {
		return "", grpclient.ErrorFailedTogenerateToken
	}
	return token, nil
}

func (m *MockGrpcService) VeriFyTokenHandler(token string) (*grpclient.VerifiedDecodedToken, error) {
	fmt.Println("verify grpc mock is called ")

	if m.GrpcErr == ErrInvalidToken {
		return nil, grpclient.ErrorInvalidToken
	}
	if m.GrpcErr == ErrUnauthorized {
		return &grpclient.VerifiedDecodedToken{
			UserId:   -5,
			Username: "",
		}, nil
	}
	return &grpclient.VerifiedDecodedToken{
		Username: "suvadip",
		UserId:   1,
	}, nil

}
