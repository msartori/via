package secret

import (
	"context"
	"os"
	"strings"
	"via/internal/log"
)

func ReadSecret(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Get().Fatal(context.Background(), err, "msg", "unable to get secret", "secret", path)
	}
	return strings.TrimSpace(string(data))
}
