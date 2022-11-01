package http

import (
	"context"
	mock_http "github.com/SpiridonovDaniil/Distributed-config/internal/app/http/mocks"
	"github.com/gofiber/fiber/v2"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestHandler_createHandler(t *testing.T) {
	type mockBehavior func(s *mock_http.Mockservice, ctx context.Context)
	testTable := []struct {
		name string
		inputBody string
		inputRequest context.Context
		mockBehavior mockBehavior
		expectedTestStatusCode int
		expectRequest error
	} {
		{
			name: "OK",
			inputBody: `{"service":"test", "data":"{"key1":"test1"}, {"key2":"test2"}"}`,
			mockBehavior: func(s *mock_http.Mockservice, ctx context.Context) {
				s.EXPECT().Create(ctx, ctx).Return(nil)
			},
			expectedTestStatusCode: 200,
			expectRequest: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			c := gomock.NewController(t)
			defer c.Finish()

			create := mock_http.NewMockservice(c)
			testCase.mockBehavior(create, testCase.inputRequest)

			f := fiber.New()

		})
	}
}