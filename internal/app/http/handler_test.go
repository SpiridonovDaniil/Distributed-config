package http

import (
	"encoding/json"
	mock_http "github.com/SpiridonovDaniil/Distributed-config/internal/app/http/mocks"
	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestHandler_createHandler(t *testing.T) {
	type mockBehavior func(s *mock_http.Mockservice, expectedError error)
	testTable := []struct {
		name                   string
		inputBody              string
		mockBehavior           mockBehavior
		expectedTestStatusCode int
		expectedError          error
		expectedResponse       string
	}{
		{
			name:      "get HTTP status 201",
			inputBody: `{"service":"test","data":[{"key1":"test1"},{"key2":"test2"}]}`,
			mockBehavior: func(s *mock_http.Mockservice, expectedError error) {
				s.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&domain.Request{})).Return(expectedError)
			},
			expectedTestStatusCode: 201,
			expectedResponse:       "",
		},
		{
			name:      "get bad request or internal server error",
			inputBody: `{"service:"test","data":[{"key1":"test1"},{"key2":"test2"}]}`,
			expectedTestStatusCode: 500,
			expectedResponse:       "[createHandler] failed to parse request, error: invalid character 't' after object key",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_http.NewMockservice(c)
			if testCase.name == "get HTTP status 201" {
				testCase.mockBehavior(create, testCase.expectedError)
			}

			f := NewServer(create)
			req, err := http.NewRequest("POST", "/config", strings.NewReader(testCase.inputBody))
			req.Header.Add("content-Type", "application/json")
			assert.NoError(t, err)

			resp, err := f.Test(req)
			assert.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.Equal(t, testCase.expectedTestStatusCode, resp.StatusCode)
			assert.Equal(t, testCase.expectedResponse, string(body))
		})
	}
}

func TestHandler_rollBackHandler(t *testing.T) {
	type mockBehavior func(s *mock_http.Mockservice, expectedError error)
	testTable := []struct {
		name                   string
		request string
		mockBehavior           mockBehavior
		expectedTestStatusCode int
		expectedError          error
		expectedResponse       string
	}{
		{
			name:      "get HTTP status 200",
			request: "/config/rollback?service=test&version=1",
			mockBehavior: func(s *mock_http.Mockservice, expectedError error) {
				s.EXPECT().RollBack(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedError)
			},
			expectedTestStatusCode: 200,
			expectedResponse:       "",
		},
		{
			name:      "get bad request 400",
			request: "/config/rollback?service=test&version=t",
			mockBehavior: func(s *mock_http.Mockservice, expectedError error) {
				s.EXPECT().RollBack(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedError)
			},
			expectedTestStatusCode: 500, 										//todo StatusCode 400
			expectedResponse:       "[RollBackHandler] version must be positive integer",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			rollBack := mock_http.NewMockservice(c)
			if testCase.name == "get HTTP status 200" {
				testCase.mockBehavior(rollBack, testCase.expectedError)
			}

			f := NewServer(rollBack)
			req, err := http.NewRequest("POST", testCase.request, strings.NewReader(""))
			req.Header.Add("content-Type", "application/json")
			assert.NoError(t, err)

			resp, err := f.Test(req)
			assert.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.Equal(t, testCase.expectedTestStatusCode, resp.StatusCode)
			assert.Equal(t, testCase.expectedResponse, string(body))
		})
	}
}

func TestHandler_getHandler(t *testing.T) {
	type mockBehavior func(s *mock_http.Mockservice, jsonMessage json.RawMessage, expectedError error)
	testTable := []struct {
		name                   string
		jsonMessage	json.RawMessage
		mockBehavior           mockBehavior
		expectedTestStatusCode int
		expectedError          error
		expectedResponse       string
	}{
		{
			name:      "get HTTP status 201",
			jsonMessage: []byte("{\"key1\":\"test1\",\"key2\":\"test2\"}"),
			mockBehavior: func(s *mock_http.Mockservice, jsonMessage json.RawMessage, expectedError error) {
				s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(jsonMessage, expectedError)
			},
			expectedTestStatusCode: 200,
			expectedResponse:       "",
		},
		{
			name:      "get bad request or internal server error",
			expectedTestStatusCode: 500,
			expectedResponse:       "[createHandler] failed to parse request, error: invalid character 't' after object key",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			get := mock_http.NewMockservice(c)
			if testCase.name == "get HTTP status 200" {
				testCase.mockBehavior(get, testCase.jsonMessage, testCase.expectedError)
			}

			f := NewServer(get)
			req, err := http.NewRequest("GET", "/config?service=test", strings.NewReader(""))
			req.Header.Add("content-Type", "application/json")
			assert.NoError(t, err)

			resp, err := f.Test(req)
			assert.NoError(t, err)

			body, err := io.ReadAll(resp.Body)
			assert.Equal(t, testCase.expectedTestStatusCode, resp.StatusCode)
			assert.Equal(t, testCase.expectedResponse, string(body))
		})
	}
}