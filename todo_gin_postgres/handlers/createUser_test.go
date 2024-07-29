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

func TestCreateUserHandler(t *testing.T) {
	mockService := &mock.MockService{}
	mockgrpcService := &mock.MockGrpcService{}
	handler := NewTaskHandler(mockService, mockgrpcService)
	router := gin.Default()
	router.POST("/signup", handler.CreateUserHandler)
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
			res:    `{"message":"user created","userid":1}`,
			inputData: `{
				"username":"suvadip",
				"password":"Suvadip@632"
			}`,
		},
		"unique user": {
			err:    mock.ErrConflict,
			status: http.StatusConflict,
			res:    `{"message":"username already exists"}`,
			inputData: `{
				"username":"suvadip",
				"password":"Suvadip@632"
			}`,
		},
		
		"validate input": {
			err:    mock.ErrInvalidInput,
			status: http.StatusBadRequest,
			res: `[{"message":"Username is too short atleast 5 character needed"}]`,
			inputData: `
				{
				"username":"some",
				"password":"Suvadip@632"
				}`,
		},
		"invalid input": {
			err:    mock.ErrInvalidInput,
			status: http.StatusBadRequest,
			res: `{"message":"Invalid input"}`,
			inputData: `{`,
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

			} else {
				mockService.MockErr = mock.Ok
			}
			requestURL := httpServer.URL + "/signup"
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
