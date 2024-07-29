package handlers

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todowithgin/mock"

	"github.com/gin-gonic/gin"
)

func TestUpdateTaskHandler(t *testing.T) {
	mockService := &mock.MockService{}
	grpcService := &mock.MockGrpcService{}
	handler := NewTaskHandler(mockService, grpcService)
	router := gin.Default()
	router.PUT("/tasks/:id", handler.PrivateRouteHandler, handler.UpdateTaskHandler)
	httpServer := httptest.NewServer(router)

	defer httpServer.Close()

	testCases := map[string]struct {
		id         string
		authHeader string
		err        mock.MockStatusCode
		grpcErr    mock.MockStatusCode
		status     int
		resp       string
		inputData  string
	}{
		"valid input": {
			id:         "1",
			authHeader: "Bearer valid_token",
			err:        mock.Ok,
			grpcErr:    mock.Ok,
			status:     http.StatusOK,
			resp:       `{"message":"Task is updated"}`,
			inputData: `{
				"title":       "task 1",
				"description": "this is task 1",
				"isCompleted": false
			}`,
		},
		"invalid input": {
			id:         "12",
			authHeader: "Bearer valid_token",
			err:        mock.ErrInvalidInput,
			grpcErr:    mock.Ok,
			status:     http.StatusBadRequest,
			resp:       `{"message":"Invalid input"}`,
			inputData:  ``,
		},
		"validate input": {
			id:         "12",
			authHeader: "Bearer valid_token",
			err:        mock.ErrInvalidInput,
			grpcErr:    mock.Ok,
			status:     http.StatusBadRequest,
			resp:       `[{"message":"Title is too short atleast 5 character needed"}]`,
			inputData: `{
				"title":       "task",
				"description": "this is task 1",
				"isCompleted":   false	
			}`,
		},
		"unique violation": {
			id:         "12",
			authHeader: "Bearer valid_token",
			err:        mock.ErrConflict,
			grpcErr:    mock.Ok,
			status:     http.StatusConflict,
			resp:       `{"message":"title already exists"}`,
			inputData: `{
				"title":       "unique task",
				"description": "this is task",
				"isCompleted": false
			}`,
		},
		"internal server error": {
			id:         "2",
			authHeader: "Bearer valid_token",
			err:        mock.ErrInternalServer,
			grpcErr:    mock.Ok,
			status:     http.StatusInternalServerError,
			resp:       `{"message":"internal server error"}`,
			inputData: `{
				"title":       "task 1",
				"description": "create task",
				"isCompleted": false
			}`,
		},
		"invalid id": {
			id:         "abc",
			authHeader: "Bearer valid_token",
			err:        mock.ErrInvalidId,
			grpcErr:    mock.Ok,
			status:     http.StatusBadRequest,
			resp:       `{"message":"Invalid id param"}`,
			inputData: `{
				"title":       "task 1",
				"description": "create task",
				"isCompleted": false
			}`,
		},
		"invalid token": {
			id:         "1",
			authHeader: "Bearer invalid token",
			err:        mock.Ok,
			grpcErr:    mock.ErrInvalidToken,
			status:     http.StatusUnauthorized,
			resp:       `{"message":"Invalid token"}`,
			inputData: `{
				"title":       "task 1",
				"description": "create task",
				"isCompleted": false
			}`,
		},
		"missing auth header": {
			id:         "1",
			authHeader: "",
			grpcErr:    mock.Ok,
			err:        mock.Ok,
			status:     http.StatusUnauthorized,
			resp:       `{"message":"missing authorization header"}`,
			inputData: `{
				"title":       "task 1",
				"description": "create task",
				"isCompleted": false
			}`,
		},
		"not found": {
			id:         "1",
			authHeader: "Bearer valid_token",
			grpcErr:    mock.Ok,
			err:        mock.ErrNotFound,
			status:     http.StatusUnauthorized,
			resp:       `{"message":"unauthorized"}`,
			inputData: `{
				"title":       "task 1",
				"description": "create task",
				"isCompleted": false
			}`,
		},
		"unauth": {
			id:         "54",
			authHeader: "Bearer valid token",
			grpcErr:    mock.ErrUnauthorized,
			err:        mock.ErrUnauthorized,
			status:     http.StatusUnauthorized,
			resp:       `{"message":"user not authorized"}`,
			inputData: `{
				"title":       "task 1",
				"description": "create task",
				"isCompleted":   false
			}`,
		},
	}

	for key, v := range testCases {
		t.Run(key, func(t *testing.T) {
			mockService.MockErr = v.err
			grpcService.GrpcErr = v.grpcErr

			requestURL := httpServer.URL + "/tasks/" + v.id

			client := http.Client{}
			requestBody := []byte(v.inputData)
			req, err := http.NewRequest(http.MethodPut, requestURL, bytes.NewBuffer(requestBody))
			if err != nil {
				t.Errorf("failed to make request %s", err)
			}

			req.Header.Set("Content-Type", "application/json")
			if v.authHeader != "" {
				req.Header.Set("Authorization", v.authHeader)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Errorf("failed request %s", err)
			}
			if res.StatusCode != v.status {
				t.Errorf("expected status code %d and got %d", v.status, res.StatusCode)
			}

			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("failed to read response body: %v", err)
			}
			defer res.Body.Close()

			actualResp := strings.TrimSpace(string(body))
			expectedResp := strings.TrimSpace(v.resp)
			if actualResp != expectedResp {
				t.Errorf("expected response body %q, got %q", expectedResp, actualResp)
			}
		})
	}
}
