package testutil

import (
	"strings"
	"testing"
	jwt_key "via/internal/jwt"
	jwt_key_mock "via/internal/jwt/mock"
	"via/internal/log"
	mock_log "via/internal/log/mock"
)

func InjectNoOpLogger() {
	log.Set(new(mock_log.MockNoOpLogger))
}

func InjectMockJWTKey() {
	jwt_key.Reset()
	jwt_key.Init(jwt_key.JWTConfig{
		PrivateKey: jwt_key_mock.GetPrivateKey(), PublicKey: jwt_key_mock.GetPublicKey()})
}

func WithTestSetup(t *testing.T, setupFunc func(t *testing.T), cleaupFunc func(), testFunc func(t *testing.T)) {
	setupFunc(t)
	t.Cleanup(cleaupFunc)
	testFunc(t)
}

func SplitSSEEvents(event string) []string {
	events := strings.Split(event, "\n\n")
	cleanEvents := []string{}
	for _, e := range events {
		cleanEvents = append(cleanEvents, strings.TrimPrefix(e, "data:"))
	}
	return cleanEvents
}
