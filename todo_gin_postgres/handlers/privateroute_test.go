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

func TestPrivaterouteHandler(t *testing.T) {

	mockgrpcService := &mock.MockGrpcService{}
	handler := NewTaskHandler(nil, mockgrpcService)

	router := gin.Default()
	router.GET("/test", handler.PrivateRouteHandler, func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"message": "success",
		})
	})
	httpServer := httptest.NewServer(router)

	testCases := map[string]struct {
		authHeader string
		grpcErr    mock.MockStatusCode
		status     int
		resp       string
	}{
		"valid token": {
			authHeader: "Bearer valid_token",
			grpcErr:    mock.Ok,
			status:     http.StatusOK,
			resp:       `{"message":"success"}`,
		},
		"missing token": {
			authHeader: "",
			grpcErr:    mock.Ok,
			status:     http.StatusUnauthorized,
			resp:       `{"message":"missing authorization header"}`,
		},
		"invalid token": {
			authHeader: "Bearer invalid token",
			grpcErr:    mock.ErrInvalidToken,
			status:     http.StatusUnauthorized,
			resp:       `{"message":"Invalid token"}`,
		},
	}

	for key, v := range testCases {
		t.Run(key, func(t *testing.T) {
			if v.grpcErr != mock.Ok {
				mockgrpcService.GrpcErr = v.grpcErr
			} else {
				mockgrpcService.GrpcErr = mock.Ok
			}
			requestURL := httpServer.URL + "/test"

			req, err := http.NewRequest(http.MethodGet, requestURL, nil)
			
			if err != nil {
				t.Fatalf("failed to create request: %v", err)
			}

			if v.authHeader != "" {
				req.Header.Set("Authorization", v.authHeader)
			}

			client := http.Client{}
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
				t.Fatalf("unable to read body: %v", err)
			}
			actualResp := strings.TrimSpace(string(body))
			expectedResp := strings.TrimSpace(v.resp)
			if actualResp != expectedResp {
				t.Errorf("expected response body %s, got %s", expectedResp, actualResp)
			}

		})

	}

}
