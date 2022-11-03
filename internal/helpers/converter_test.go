package helpers

import (
	"encoding/json"
	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Converter(t *testing.T) {
	testTable := []struct {
		name          string
		inputBody     *domain.Request
		expectedError error
		expectedBody  []byte
	}{
		{
			name: "successfully",
			inputBody: &domain.Request{
				Service: `{"test"}`,
				Data:    json.RawMessage(`[{"key1":"test1"},{"key2":"test2"}]`),
			},
			expectedError: nil,
			expectedBody:  json.RawMessage(`{"key1":"test1","key2":"test2"}`),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			actual, errActual := Converter(testCase.inputBody)
			assert.Equal(t, testCase.expectedBody, actual)
			assert.Equal(t, testCase.expectedError, errActual)
		})
	}
}
