package secret

import (
	"os"
	"strings"
	util_logger "via/internal/util/logger"
)

func ReadSecret(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		util_logger.Get().Fatal(err, "msg", "unable to get secret", "secret", path)
	}
	return strings.TrimSpace(string(data))
}
