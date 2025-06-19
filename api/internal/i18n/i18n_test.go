package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWithLang(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		key      string
		args     []interface{}
		expected string
	}{
		{
			name:     "known key in es",
			lang:     "es",
			key:      MsgGuideRequired,
			expected: "El código de guía es requerido.",
		},
		{
			name:     "known key in en",
			lang:     "en",
			key:      MsgRequestTimeout,
			expected: "Request timeout.",
		},
		{
			name:     "unknown key in known lang",
			lang:     "es",
			key:      "unknown_key",
			expected: "unknown_key",
		},
		{
			name:     "unknown lang",
			lang:     "fr",
			key:      MsgGuideRequired,
			expected: "guide_required",
		},
		{
			name:     "format string with args",
			lang:     "es",
			key:      MsgInternalServerError,
			args:     []any{"extra"},
			expected: "Error interno del servidor.", // still matches, no placeholder
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetWithLang(tt.lang, tt.key, tt.args...)
			if got != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		key      string
		args     []interface{}
		expected string
	}{
		{
			name:     "from header - es",
			lang:     "es",
			key:      MsgGuideInvalid,
			expected: "El código de guía es inválido.",
		},
		{
			name:     "from header - en",
			lang:     "en",
			key:      MsgRequestCanceledByClient,
			expected: "Request canceled by client.",
		},
		{
			name:     "missing header",
			lang:     "", // no Accept-Language
			key:      MsgGuideNotFound,
			expected: "Guía no econtrada.",
		},
		{
			name:     "unknown key",
			lang:     "es",
			key:      "foo",
			expected: "foo",
		},
		{
			name:     "unknown language",
			lang:     "de",
			key:      MsgRequestTimeout,
			expected: "request_timeout",
		},
		{
			name:     "format string with args",
			lang:     "en",
			key:      MsgUnexpectedContextError,
			args:     []interface{}{"ignored"},
			expected: "Unexpected context error.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.lang != "" {
				req.Header.Set("Accept-Language", tt.lang)
			}
			got := Get(req, tt.key, tt.args...)
			if got != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, got)
			}
		})
	}
}
