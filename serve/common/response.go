package common

import (
	"encoding/json"
)

type Response struct {
	Code    int             `json:"code"`
	ErrCode string          `json:"err_code"`
	Body    json.RawMessage `json:"body"`
}
