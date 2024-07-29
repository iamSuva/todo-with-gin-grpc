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

func TestLoginHandler(t *testing.T) {
	mockService := &mock.MockService{}
	mockgrpcService := &mock.MockGrpcService{}
	handler := NewTaskHandler(mockService, mockgrpcService)
	router := gin.Default()
	router.POST("/login", handler.LoginHandler)
	httpServer := httptest.NewServer(router)
	testcases := map[string]struct {
		err       mock.MockStatusCode
		status    int
		res       string
		inputData string
	}{
		"valid input": {
			err:    mock.Ok,
			status: http.StatusOK,
			res:    `{"token":"h34h34hjg3t33bjnfjgj45kjgj333hjfkj4jnj54nn"}`,
			inputData: `{
				"username":"suvadip",
				"password":"suvadip@632"
			}`,
		},
		"failed to generate token": {
			err: mock.ErrFailedToGenerate,
			status: http.StatusInternalServerError,
			res:    `{"message":"failed to generate token"}`,
			inputData: `{
				"username":"suvadip",
				"password":"suvadip@632"
			}`,
		},
		"invalid input": {
			err: mock.ErrInvalidInput,

			status: http.StatusBadRequest,
			res:    `{"message":"Invalid input"}`,
			inputData: `{
				"username":"suvadip",
			}`,
		},
		"user not found": {
			err: mock.ErrNotFound,

			status: http.StatusNotFound,
			res:    `{"message":"User not found"}`,
			inputData: `{
				"username":"suvdip",
				"password":"Suvadip&632"
			}`,
		},
		"password not match": {
			err:    mock.ErrPassword,
			status: http.StatusUnauthorized,
			res:    `{"message":"password not match"}`,
			inputData: `{
				"username":"suvdip",
				"password":"Suvadip@632"
			}`,
		},
		"internal server": {
			err:    mock.ErrInternalServer,
			status: http.StatusInternalServerError,
			res:    `{"message":"internal server error"}`,
			inputData: `{
				"username":"suvdip",
				"password":"Suvadip@632"
			}`,
		},
	}
	for key, v := range testcases {
		t.Run(key, func(t *testing.T) {
			if v.err != mock.Ok {
				mockService.MockErr = v.err
				mockgrpcService.GrpcErr = v.err
			} else {
				mockService.MockErr = mock.Ok
			}
			requestURL := httpServer.URL + "/login"
			client := http.Client{}
			requestBody := []byte(v.inputData)
			req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestBody))
			//creates a *bytes.Buffer containing the request body bytes.
			//Buffer is used to handle the byte slice (requestBody) as an io.Reader
			if err != nil {
				t.Errorf("failed to make request %s", err)
			}
			req.Header.Set("Content-Type", "application/json")
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
			expectedResp := strings.TrimSpace(v.res)
			if actualResp != expectedResp {
				t.Errorf("expected response body %q, got %q", expectedResp, actualResp)
			}
		})
	}

}
