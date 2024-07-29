package handlers

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	grpclient "todowithgin/grpcClient"
	"todowithgin/mock"

	"github.com/gin-gonic/gin"
)

func TestGetTasksHandler(t *testing.T){
	mockService := &mock.MockService{}
	grpcService := &grpclient.GrpcService{}
	handler := NewTaskHandler(mockService, grpcService)
	router := gin.Default()
	router.GET("/tasks", handler.GetTasksHandler)
	httpServer := httptest.NewServer(router)
	cases := map[string]struct {
		err    mock.MockStatusCode
		status int
	}{
		"Internal error": {
			err:    mock.ErrInternalServer,
			status: http.StatusInternalServerError,
		},
		"successful ": {
			err:    mock.Ok,
			status: http.StatusOK,
		},
	}
	for key, v := range cases {
		t.Run(key, func(t *testing.T) {
			if v.err != mock.Ok {
				mockService.MockErr = v.err
			} else {
				mockService.MockErr = mock.Ok
			}
			client := http.Client{} //create a instance of http.Client struct

			requestURL := httpServer.URL + "/tasks"
			req, err := http.NewRequest(http.MethodGet, requestURL, nil)
			if err != nil {
				t.Fatalf("request creation is failed %v", req)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("request failed %v", err)
			}
			defer res.Body.Close()
			if res.StatusCode != v.status {
				t.Errorf("expected status %d got status %d ", v.status, res.StatusCode)
			}
			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			fmt.Println(string(body))
		})
	}

}
