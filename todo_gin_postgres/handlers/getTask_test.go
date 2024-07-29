package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	grpclient "todowithgin/grpcClient"
	"todowithgin/mock"
	"github.com/gin-gonic/gin"
)

func TestGetTaskHandler(t *testing.T) {
	mockService := &mock.MockService{}
	grpcService := &grpclient.GrpcService{}
	handler := NewTaskHandler(mockService, grpcService)

	router := gin.Default()
	router.GET("/tasks/:id", handler.GetTaskHandler)

	httpServer := httptest.NewServer(router)
	baseUrl := httpServer.URL
	fmt.Println(baseUrl)
	testCases := map[string]struct {
		id     string
		err    mock.MockStatusCode
		status int
		res    string
	}{
		"valid id": {
			id:     "1",
			err:    mock.Ok,
			status: http.StatusOK,
			res:    `{"data":{"id":1,"title":"Task 1","description":"this is task 1","isCompleted":false,"createdAt_utc":"%s","updatedAt_utc":"%s"}}`,
		},
		"invalid id": {
			id:     "abc",
			err:    mock.ErrInvalidId,
			status: http.StatusBadRequest,
			res:    `{"message":"Invalid id param"}`,
		},
		"not found": {
			id:     "100",
			err:    mock.ErrNotFound,
			status: http.StatusNotFound,
			res:    `{"message":"no data found"}`,
		},
		"internal server": {
			id:     "12",
			err:    mock.ErrInternalServer,
			status: http.StatusInternalServerError,
			res:    `{"message":"internal server error"}`,
		},
	}

	for key, v := range testCases {
		t.Run(key, func(t *testing.T) {
			if v.err != mock.Ok {
				mockService.MockErr = v.err
			} else {
				mockService.MockErr = mock.Ok
			}
			url := httpServer.URL + "/tasks/" + v.id
			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				t.Errorf("failed to make request %s ", err)
				return
			}
			client := http.Client{}
			res, err := client.Do(req)
			if err != nil {
				t.Errorf("failed request %s ", err)
				return
			}
			if res.StatusCode != v.status {
				t.Errorf("expected status %d, got %d", v.status, res.StatusCode)
			}
			defer res.Body.Close()

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			result := string(body)
			expectedOutput := v.res
			if v.err == mock.Ok {
				staticTime := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
				expectedOutput = fmt.Sprintf(v.res, staticTime.Format(time.RFC3339), staticTime.Format(time.RFC3339))
			}

			actualOutput := strings.TrimSpace(result)
			if expectedOutput != actualOutput {
				t.Errorf("expected: %s and got : %s", expectedOutput, actualOutput)
				return
			}

		})
	}

}
