package biz_guide_status

import (
	"testing"

	biz_config "via/internal/biz/config"
	biz_language "via/internal/biz/language"

	"github.com/stretchr/testify/assert"
)

func TestGetStatusDescription(t *testing.T) {
	tests := []struct {
		lang     string
		status   string
		expected string
	}{
		{biz_language.ES, INITIAL, "Inicial"},
		{biz_language.ES, "unknown", "Unknown status: unknown"},
		{"fr", INITIAL, "Inicial"}, // fallback to DEFAULT (ES)
		{"fr", "unknown", "Unknown status: unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.status+"_"+tt.lang, func(t *testing.T) {
			result := GetStatusDescription(tt.lang, tt.status)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetNextStatus(t *testing.T) {
	tests := []struct {
		name          string
		currentStatus string
		history       []string
		payment       string
		expected      []string
	}{
		{
			name:          "next status regular",
			currentStatus: INITIAL,
			expected:      []string{ON_HOLD, SUSPENDED, PENDING_RECIPIENT_IDENTIFY},
		},
		{
			name:          "next status previous with history",
			currentStatus: ON_HOLD,
			history:       []string{INITIAL, PENDING_PAYMENT},
			expected:      []string{INITIAL},
		},
		{
			name:          "next status previous with empty history",
			currentStatus: ON_HOLD,
			history:       []string{},
			expected:      []string{},
		},
		{
			name:          "next status with payment",
			currentStatus: RECIPIENT_IDENTIFIED,
			payment:       biz_config.PAID_SHIPPING,
			expected: []string{
				ON_HOLD,
				SUSPENDED,
				PENDING_COUNTER_DELIVERY,
				PENDING_WAREHOUSE_DELIVERY,
			},
		},
		{
			name:          "unknown status returns empty",
			currentStatus: "unknown",
			expected:      []string{},
		},
		{
			name:          "next status previous with history, on hold as previous first",
			currentStatus: ON_HOLD,
			history:       []string{ON_HOLD, INITIAL, PENDING_PAYMENT},
			expected:      []string{INITIAL},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetNextStatus(tt.currentStatus, tt.history, tt.payment)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestGetStatusState(t *testing.T) {
	assert.Equal(t, "error", GetStatusState(SUSPENDED))
	assert.Equal(t, "warn", GetStatusState(ON_HOLD))
	assert.Equal(t, "ok", GetStatusState(INITIAL))
}

func TestIsEnabledToWithdraw(t *testing.T) {
	assert.True(t, IsEnabledToWithdraw(ON_HOLD))
	assert.True(t, IsEnabledToWithdraw(PARTIAL_DELIVERED))
	assert.False(t, IsEnabledToWithdraw(INITIAL))
}

func TestIsInProcess(t *testing.T) {
	statuses := []string{
		INITIAL,
		PENDING_RECIPIENT_IDENTIFY,
		RECIPIENT_IDENTIFIED,
		PENDING_PAYMENT,
		PAID,
		PENDING_COUNTER_DELIVERY,
		PENDING_WAREHOUSE_DELIVERY,
	}
	for _, status := range statuses {
		assert.True(t, IsInProcess(status), status)
	}
	assert.False(t, IsInProcess(DELIVERED))
}

func TestIsDelivered(t *testing.T) {
	assert.True(t, IsDelivered(DELIVERED))
	assert.False(t, IsDelivered(INITIAL))
}

func TestIsAbleToReInit(t *testing.T) {
	assert.True(t, IsAbleToReInit(ON_HOLD))
	assert.False(t, IsAbleToReInit(SUSPENDED))
}

func TestIsValidToCreateForWithdraw(t *testing.T) {
	assert.True(t, IsValidToCreateForWithdraw(PARTIAL_DELIVERED))
	assert.False(t, IsValidToCreateForWithdraw(ON_HOLD))
}

func TestGetMonitorStatus(t *testing.T) {
	result := GetMonitorStatus()
	assert.NotEmpty(t, result)
}

func TestGetOperatorStatus(t *testing.T) {
	result := GetOperatorStatus()
	assert.NotEmpty(t, result)
}
