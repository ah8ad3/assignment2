package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestServer_PostData(t *testing.T) {
	stopC := make(chan bool)
	defer close(stopC)
	server := newServer(stopC)

	// Test cases
	tests := []struct {
		name           string
		input          Input
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid Input",
			input: Input{
				UniqueID: 1,
				UserID:   1,
			},
			expectedStatus: 201,
			expectedError:  "",
		},
		{
			name: "Missing UniqueID",
			input: Input{
				UserID: 1,
			},
			expectedStatus: 400,
			expectedError:  "unique_id and user_id must be present and valid",
		},
		{
			name: "Missing UserID",
			input: Input{
				UniqueID: 1,
			},
			expectedStatus: 400,
			expectedError:  "unique_id and user_id must be present and valid",
		},
		{
			name: "Missing UserID",
			input: Input{
				UniqueID: 0,
				UserID:   0,
			},
			expectedStatus: 400,
			expectedError:  "unique_id and user_id must be present and valid",
		},
		{
			name: "User Does Not Exist",
			input: Input{
				UniqueID: 2,
				UserID:   3,
			},
			expectedStatus: 403,
			expectedError:  "user does not have any quota",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := createRequestBody(test.input)

			// Create a response recorder.
			rr := httptest.NewRecorder()

			// Call the handler.
			server.PostData(rr, req)

			// Check the status code.
			assert.Equal(t, test.expectedStatus, rr.Code)

			// Check the response body.
			var response map[string]any
			err := json.Unmarshal(rr.Body.Bytes(), &response)
			assert.NoError(t, err)

			// Check the error message.
			if test.expectedError != "" {
				assert.Equal(t, test.expectedError, response["text"])
			}
		})
	}
}

// Helper function to create a request body.
func createRequestBody(input Input) *http.Request {
	body, err := json.Marshal(input)
	if err != nil {
		log.Fatal(err)
	}
	return httptest.NewRequest(http.MethodPost, "/kv",
		bytes.NewBuffer(body))
}
