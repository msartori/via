package biz_config

import (
	"testing"
	biz_language "via/internal/biz/language"

	"github.com/stretchr/testify/assert"
)

func TestGetPaymentDescription(t *testing.T) {
	tests := []struct {
		name     string
		lang     string
		payment  string
		expected string
	}{
		{
			name:     "Known language and payment (ES, P)",
			lang:     biz_language.ES,
			payment:  PAID_SHIPPING,
			expected: "Origen",
		},
		{
			name:     "Known language and payment (ES, D)",
			lang:     biz_language.ES,
			payment:  PAID_ON_DESTINATION,
			expected: "Destino",
		},
		{
			name:     "Unknown payment with known language",
			lang:     biz_language.ES,
			payment:  "X",
			expected: "Unknown payment: X",
		},
		{
			name:     "Unknown language, fallback to DEFAULT missing payment",
			lang:     "FR",
			payment:  "X",
			expected: "Unknown payment: X",
		},
		{
			name:     "Unknown language, fallback to DEFAULT known payment",
			lang:     "FR",
			payment:  "X",
			expected: "Unknown payment: X", // No fallback because DEFAULT is missing in paymentDescription
		},
		{
			name:     "Unknown language, fallback to DEFAULT known payment",
			lang:     "FR",
			payment:  "P",
			expected: "Origen",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPaymentDescription(tt.lang, tt.payment)
			assert.Equal(t, tt.expected, got)
		})
	}
}
