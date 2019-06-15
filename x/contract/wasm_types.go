package contract

import (
	"encoding/json"
	"errors"
)

type SendResponse struct {
	Error  string `json:"error"`
	Result string `json:"result"`
}

func ParseResponse(raw string) (*SendResponse, error) {
	var out SendResponse
	err := json.Unmarshal([]byte(raw), &out)
	if err != nil {
		return nil, err
	}
	if out.Error != "" {
		return nil, errors.New(out.Error)
	}
	return &out, nil
}
