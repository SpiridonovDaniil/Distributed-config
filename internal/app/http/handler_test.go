package http

import (
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
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_http.NewMockservice(c)
			testCase.mockBehavior(create, testCase.expectedError)

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
		inputBody              string
		mockBehavior           mockBehavior
		expectedTestStatusCode int
		expectedError          error
		expectedResponse       string
	}{
		{
			name:      "get HTTP status 200",
			inputBody: `{"service":"test","data":[{"key1":"test1"},{"key2":"test2"}]}`,
			mockBehavior: func(s *mock_http.Mockservice, expectedError error) {
				s.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(&domain.Request{})).Return(expectedError)
			},
			expectedTestStatusCode: 200,
			expectedResponse:       "",
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_http.NewMockservice(c)
			testCase.mockBehavior(create, testCase.expectedError)

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
