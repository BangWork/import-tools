package utils

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"unicode/utf8"

	"github.com/bangwork/import-tools/serve/utils/uuid"
)

func UUID() string {
	return uuid.V4Compressed()[0:8]
}

func OutputJSON(input interface{}) []byte {
	j, err := json.Marshal(input)
	if err != nil {
		log.Printf("json marshal error: %+v", err)
		return []byte{}
	}
	return j
}

func TruncateString(s string, maxRuneCount int) string {
	if utf8.RuneCountInString(s) <= maxRuneCount {
		return s
	}
	runes := []rune(s)[:maxRuneCount]
	return string(runes)
}

func Base64Encode(input string) string {
	i := []byte(input)
	return base64.StdEncoding.EncodeToString(i)
}

func Base64Decode(input string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(input)
	return string(decoded), err
}
