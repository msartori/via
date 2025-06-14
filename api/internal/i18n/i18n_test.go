package i18n

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetWithLang(t *testing.T) {
	t.Run("should return spanish message without args", func(t *testing.T) {
		got := GetWithLang("es", MsgGuideRequired)
		want := "El código de guía es requerido."
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("should return spanish message with args", func(t *testing.T) {
		got := GetWithLang("es", MsgOtherBranch, "Rosario")
		want := "La guía solicitada corresponde a la sucursal Rosario."
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("should return english message", func(t *testing.T) {
		got := GetWithLang("en", MsgRequestTimeout)
		want := "Request timeout."
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("should fallback to key if lang unknown", func(t *testing.T) {
		got := GetWithLang("fr", MsgGuideRequired)
		if got != MsgGuideRequired {
			t.Errorf("expected fallback key %q, got %q", MsgGuideRequired, got)
		}
	})

	t.Run("should fallback to key if key unknown", func(t *testing.T) {
		got := GetWithLang("es", "nonexistent_key")
		if got != "nonexistent_key" {
			t.Errorf("expected fallback to key, got %q", got)
		}
	})
}

func TestGet(t *testing.T) {
	t.Run("should return message based on Accept-Language header", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Language", "es")

		got := Get(req, MsgDelivered)
		want := "La guía solicitada ya se ha entregado."
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("should fallback to es if header missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)

		got := Get(req, MsgWithdrawAvailable)
		want := "La guía solicitada está disponible para su retiro."
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("should format with args", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Language", "es")

		got := Get(req, MsgOtherBranch, "Bahía Blanca")
		want := "La guía solicitada corresponde a la sucursal Bahía Blanca."
		if got != want {
			t.Errorf("expected %q, got %q", want, got)
		}
	})

	t.Run("should fallback to key if not found", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Accept-Language", "es")

		got := Get(req, "nonexistent_key")
		if got != "nonexistent_key" {
			t.Errorf("expected fallback key, got %q", got)
		}
	})
}
