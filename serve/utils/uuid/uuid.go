package uuid

import (
	"crypto/rand"
	"log"

	"github.com/bangwork/import-tools/serve/utils/base58"
)

func genUUIDv4() []byte {
	var u = make([]byte, 16)
	_, err := rand.Read(u[:])
	if err != nil {
		log.Printf("gen uuid error: %v", err)
	}
	// Set version
	u[6] = (u[6] & 0x0F) | (4 << 4)
	// Set variant bits
	u[8] = (u[8] | 0x40) & 0x7F
	return u
}

// 生成 Base58 表示的 RFC4122 V4 版本的 UUID，长度为 22
func V4Compressed() string {
	return base58.Encode(genUUIDv4())
}
