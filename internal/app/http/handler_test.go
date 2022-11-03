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
			name:                   "get bad request or internal server error",
			inputBody:              `{"service:"test","data":[{"key1":"test1"},{"key2":"test2"}]}`,
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
		request                string
		mockBehavior           mockBehavior
		expectedTestStatusCode int
		expectedError          error
		expectedResponse       string
	}{
		{
			name:    "get HTTP status 200",
			request: "/config/rollback?service=test&version=1",
			mockBehavior: func(s *mock_http.Mockservice, expectedError error) {
				s.EXPECT().RollBack(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedError)
			},
			expectedTestStatusCode: 200,
			expectedResponse:       "",
		},
		{
			name:    "get bad request 400",
			request: "/config/rollback?service=test&version=t",
			mockBehavior: func(s *mock_http.Mockservice, expectedError error) {
				s.EXPECT().RollBack(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedError)
			},
			expectedTestStatusCode: 400,
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
		jsonMessage            json.RawMessage
		mockBehavior           mockBehavior
		expectedTestStatusCode int
		expectedError          error
		expectedResponse       string
	}{
		{
			name:        "get HTTP status 200",
			jsonMessage: []byte(`{"key1":"test1","key2":"test2"}`),
			mockBehavior: func(s *mock_http.Mockservice, jsonMessage json.RawMessage, expectedError error) {
				s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(jsonMessage, expectedError)
			},
			expectedTestStatusCode: 200,
			expectedResponse:       `{"key1":"test1","key2":"test2"}`,
		},
		{
			name:        "get internal server error",
			jsonMessage: []byte(`{"key1"":"test1","key2":"test2"}`),
			mockBehavior: func(s *mock_http.Mockservice, jsonMessage json.RawMessage, expectedError error) {
				s.EXPECT().Get(gomock.Any(), gomock.Any()).Return(jsonMessage, expectedError)
			},
			expectedTestStatusCode: 500,
			expectedResponse:       "[getHandler] failed to return JSON answer, error: json: error calling MarshalJSON for type json.RawMessage: invalid character '\"' after object key",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			get := mock_http.NewMockservice(c)
			testCase.mockBehavior(get, testCase.jsonMessage, testCase.expectedError)

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

func TestHandler_getVersionsHandler(t *testing.T) {
	type mockBehavior func(s *mock_http.Mockservice, serviceAnswer []*domain.Config, expectedError error)
	testTable := []struct {
		name                   string
		serviceAnswer          []*domain.Config
		mockBehavior           mockBehavior
		expectedTestStatusCode int
		expectedError          error
		expectedResponse       string
	}{
		{
			name: "get HTTP status 200",
			serviceAnswer: []*domain.Config{
				{
					Config:  json.RawMessage(`{"key1": "valya","key2": "petya"}`),
					Version: 1,
				},
				{
					Config:  json.RawMessage(`{"key1": "valy","key2": "pety"}`),
					Version: 2,
				},
			},
			mockBehavior: func(s *mock_http.Mockservice, serviceAnswer []*domain.Config, expectedError error) {
				s.EXPECT().GetVersions(gomock.Any(), gomock.Any()).Return(serviceAnswer, expectedError)
			},
			expectedTestStatusCode: 200,
			expectedResponse:       `[{"config":{"key1":"valya","key2":"petya"},"version":1},{"config":{"key1":"valy","key2":"pety"},"version":2}]`,
		},
		{
			name: "get internal server error",
			serviceAnswer: []*domain.Config{
				{
					Config:  json.RawMessage(`{"key1": "valya","key2": "petya"}`),
					Version: 1,
				},
				{
					Config:  json.RawMessage(`"\"'{"key1": "valy","key2": "pety"}`),
					Version: 2,
				},
			},
			mockBehavior: func(s *mock_http.Mockservice, serviceAnswer []*domain.Config, expectedError error) {
				s.EXPECT().GetVersions(gomock.Any(), gomock.Any()).Return(serviceAnswer, expectedError)
			},
			expectedTestStatusCode: 500,
			expectedResponse:       "[getVersionsHandler] failed to return JSON answer, error: json: error calling MarshalJSON for type json.RawMessage: invalid character 'k' after top-level value",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			getVersion := mock_http.NewMockservice(c)
			testCase.mockBehavior(getVersion, testCase.serviceAnswer, testCase.expectedError)

			f := NewServer(getVersion)
			req, err := http.NewRequest("GET", "/config/versions?service=test", strings.NewReader(""))
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

func TestHandler_putHandler(t *testing.T) {

	type mockBehavior func(s *mock_http.Mockservice, expectedError error)
	testTable := []struct {
		name                   string
		InputBody              string
		mockBehavior           mockBehavior
		expectedTestStatusCode int
		expectedError          error
		expectedResponse       string
	}{
		{
			name:      "get HTTP status 200",
			InputBody: "{\n \"service\": \"lowef2ef\",\n   \"data\": [\n      {\"key1\": \"valy\"}, \n      {\"key2\": \"pety\"}\n    ]\n}",
			mockBehavior: func(s *mock_http.Mockservice, expectedError error) {
				s.EXPECT().Update(gomock.Any(), gomock.Any()).Return(expectedError)
			},
			expectedTestStatusCode: 200,
			expectedResponse:       "",
		},
		{
			name:      "get bad request error 400",
			InputBody: "{\n \"service\": \"lowef2ef\",\n   \"data\": [\n      {\"\"key1\": \"valy\"}, \n      {\"key2\": \"pety\"}\n    ]\n}",
			mockBehavior: func(s *mock_http.Mockservice, expectedError error) {
			},
			expectedTestStatusCode: 400,
			expectedResponse:       "[putHandler] failed to parse request, error: invalid character 'k' after object key",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			update := mock_http.NewMockservice(c)
			testCase.mockBehavior(update, testCase.expectedError)

			f := NewServer(update)
			req, err := http.NewRequest("PUT", "/config", strings.NewReader(testCase.InputBody))
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
