package secret

import (
	"context"
	"os"
	"strings"
	"sync"
	"via/internal/log"
)

var (
	instance SecretReader
	mutex    = &sync.RWMutex{}
)

type SecretReader interface {
	Read(path string) string
}

type FileSecretReader struct {
}

func (s *FileSecretReader) Read(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		//will always be used to initialize critic instances
		log.Get().Fatal(context.Background(), err, "msg", "unable to get secret", "secret", path)
	}
	return strings.TrimSpace(string(data))
}

func Get() SecretReader {
	mutex.RLock()
	defer mutex.RUnlock()
	return instance
}

func Set(sr SecretReader) {
	mutex.Lock()
	defer mutex.Unlock()
	instance = sr
}
