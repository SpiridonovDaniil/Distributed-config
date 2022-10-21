package domain

import "encoding/json"

type Config struct {
	Service string          `json:"service"`
	Data    json.RawMessage `json:"data"`
}
