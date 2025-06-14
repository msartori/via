package guide_web_provider

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"testing"
	http_client "via/internal/client/http"
	mock_http "via/internal/client/http/mock"
	"via/internal/log"
	mock_log "via/internal/log/mock"
	"via/internal/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockParser struct {
	err   error
	guide model.Guide
}

func (m *mockParser) Parse(data []byte, v any) error {
	ptr, ok := v.(*model.Guide)
	if ok {
		*ptr = m.guide
	}
	return m.err
}

func TestGetGuide(t *testing.T) {
	mockLog := new(mock_log.MockNoOpLogger)
	// logger mock setup
	log.Set(mockLog)

	tests := []struct {
		name         string
		mockResp     *http.Response
		mockErr      error
		parserErr    error
		expectErr    bool
		expectGuide  model.Guide
		parserGuide  model.Guide
		responseBody string
		statusCode   int
	}{
		{
			name:      "http client error",
			mockErr:   errors.New("http error"),
			expectErr: true,
		},
		{
			name:         "non-200 status code",
			statusCode:   http.StatusBadRequest,
			responseBody: "error",
			expectErr:    true,
			mockErr:      nil,
			mockResp: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       io.NopCloser(bytes.NewBufferString("error")),
			},
		},
		{
			name: "body read error",
			mockResp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       &errorReadCloser{},
			},
			expectErr: true,
		},
		{
			name:         "parser returns generic error",
			statusCode:   http.StatusOK,
			responseBody: "html...",
			parserErr:    errors.New("parse error"),
			mockResp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("html...")),
			},
			expectErr: true,
		},
		{
			name:         "parser returns ErrNoResultRow",
			statusCode:   http.StatusOK,
			responseBody: "html...",
			parserErr:    ErrNoResultRow,
			parserGuide:  model.Guide{ID: "123"},
			mockResp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("html...")),
			},
			expectGuide: model.Guide{ID: "123"},
		},
		{
			name:         "success",
			statusCode:   http.StatusOK,
			responseBody: "<html>...</html>",
			parserGuide:  model.Guide{ID: "999"},
			mockResp: &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(bytes.NewBufferString("<html>...</html>")),
			},
			expectGuide: model.Guide{ID: "999"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRequester := new(mock_http.MockHttpRequester)
			mockRequester.On("Do", mock.Anything).Return(tt.mockResp, tt.mockErr)

			client := &http_client.HttpClient{
				Requester: mockRequester,
				BaseURL:   "http://fake-url.com",
			}

			parser := &mockParser{
				err:   tt.parserErr,
				guide: tt.parserGuide,
			}

			provider := &GuideWebProvider{
				client:      client,
				guideParser: parser,
			}

			guide, err := provider.GetGuide(context.Background(), "123456789012")

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectGuide, guide)
			}
		})
	}
}

// errorReadCloser simula error al leer el cuerpo
type errorReadCloser struct{}

func (e *errorReadCloser) Read(p []byte) (int, error) {
	return 0, errors.New("read error")
}
func (e *errorReadCloser) Close() error {
	return nil
}
