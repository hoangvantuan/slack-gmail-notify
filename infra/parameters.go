package infra

import (
	"os"
)

var (
	// ENCRYPTKey was use to encode token
	ENCRYPTKey string
)

func setupEncryptKey() {
	ENCRYPTKey = os.Getenv("ENCRYPT_KEY")
}
