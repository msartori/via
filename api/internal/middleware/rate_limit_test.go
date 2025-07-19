package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	biz_language "via/internal/biz/language"
	"via/internal/i18n"
	"via/internal/ratelimit"
	"via/internal/testutil"

	mock_ds "via/internal/ds/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRateLimitMiddleware(t *testing.T) {
	testutil.InjectNoOpLogger()

	type testCase struct {
		name         string
		path         string
		expectedCode int
		expectBody   string
	}
	mockDS := new(mock_ds.MockDS)
	mockDS.On("Incr", mock.Anything, mock.Anything).Return(0, errors.New("ds error")).Once()
	mockDS.On("Incr", mock.Anything, mock.Anything).Return(2, nil).Once()
	mockDS.On("Incr", mock.Anything, mock.Anything).Return(6, nil).Once()
	rateLimit := RateLimitMiddleware{
		RateLimiter: ratelimit.New("rateLimiter", 5, 1*time.Second, mockDS),
		KeyGetter:   func(r *http.Request) string { return "key" },
	}
	tests := []testCase{
		{
			name:         "No rate limiter for path",
			path:         "/no-limit",
			expectedCode: http.StatusOK,
			expectBody:   "handler ok\n",
		},
		{
			name:         "Rate limiter returns error",
			path:         "/error",
			expectedCode: http.StatusInternalServerError,
			expectBody:   i18n.GetWithLang(biz_language.DEFAULT, i18n.MsgInternalServerError),
		},
		{
			name:         "Allowed",
			path:         "/allowed",
			expectedCode: http.StatusOK,
			expectBody:   "handler ok\n",
		},
		{
			name:         "Rate limited",
			path:         "/rate-limited",
			expectedCode: http.StatusTooManyRequests,
			expectBody:   i18n.GetWithLang(biz_language.DEFAULT, i18n.MsgTooManyRequestsError),
		},
	}
	limiters := map[string]RateLimitMiddleware{}
	for i, tt := range tests {
		if i > 0 { //first is not rate limited
			limiters[tt.path] = rateLimit
		}
	}

	rateLimitMiddleware := NewRateLimitMiddleware(limiters)

	handler := rateLimitMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("handler ok\n"))
	}))

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)
			assert.Equal(t, tt.expectedCode, rec.Code)
			assert.Contains(t, rec.Body.String(), tt.expectBody)
		})
	}
}
