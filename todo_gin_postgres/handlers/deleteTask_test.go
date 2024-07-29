package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"todowithgin/mock"

	"github.com/gin-gonic/gin"
)

func TestDeleteTaskHandler(t *testing.T) {
	mockdbService := &mock.MockService{}
	grpcService := &mock.MockGrpcService{}
	handler := NewTaskHandler(mockdbService, grpcService)
	router := gin.Default()
	router.DELETE("/tasks/:id",handler.PrivateRouteHandler ,handler.DeleteTaskHandler)
	httpServer := httptest.NewServer(router)
	defer httpServer.Close() 

	testCases := map[string]struct {
		id         string
		authHeader string
		err        mock.MockStatusCode
		grpcErr    mock.MockStatusCode
		status     int
		resp       string
	}{
		"valid id": {
			id:         "1",
			authHeader: "Bearer valid_token",
			err:        mock.Ok,
			grpcErr:    mock.Ok,
			status:     http.StatusOK,
			resp:       `{"message":"task is deleted"}`,
		},
		"invalid id": {
			id:         "abc",
			authHeader: "Bearer valid_token",
			err:        mock.ErrInvalidId,
			grpcErr:    mock.Ok,
			status:     http.StatusBadRequest,
			resp:       `{"message":"Invalid id param"}`,
		},
		"not found": {
			id:         "123",
			authHeader: "Bearer valid_token",
			err:        mock.ErrNotFound,
			grpcErr:    mock.Ok,
			status:     http.StatusOK,
			resp:       `{"message":"task is deleted"}`,
		},
		"internal server error": {
			id:         "12",
			authHeader: "Bearer valid_token",
			err:        mock.ErrInternalServer,
			grpcErr:    mock.Ok,
			status:     http.StatusInternalServerError,
			resp:       `{"message":"internal server error"}`,
		},
		"invalid token": {
			id:         "1",
			authHeader: "Bearer invalid_token",
			err:        mock.Ok,
			grpcErr:    mock.ErrInvalidToken,
			status:     http.StatusUnauthorized,
			resp:       `{"message":"Invalid token"}`,
		},
		"missing auth header": {
			id:         "1",
			authHeader: "",
			grpcErr:    mock.Ok,
			err:        mock.Ok,
			status:     http.StatusUnauthorized,
			resp:       `{"message":"missing authorization header"}`,
		},
		
	}

	for key, v := range testCases {
		t.Run(key, func(t *testing.T) {
			mockdbService.MockErr = v.err
			grpcService.GrpcErr = v.grpcErr

			client := http.Client{}
			requestURL := httpServer.URL + "/tasks/" + v.id
			req, err := http.NewRequest(http.MethodDelete, requestURL, nil)
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			if v.authHeader != "" {
				req.Header.Set("Authorization", v.authHeader)
			}

			res, err := client.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer res.Body.Close()

			if res.StatusCode != v.status {
				t.Errorf("expected status %d, got %d", v.status, res.StatusCode)
			}

			body, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatalf("not able to read body: %s", err)
			}
			actualResp := strings.TrimSpace(string(body))
			expectedResp := strings.TrimSpace(v.resp)
			if actualResp != expectedResp {
				t.Errorf("expected response body %q, got %q", expectedResp, actualResp)
			}
		})
	}
}
