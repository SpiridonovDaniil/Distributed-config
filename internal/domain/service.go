package domain

import "encoding/json"

type Request struct {
	Service string          `json:"service"`
	Data    json.RawMessage `json:"data"`
}

type Config struct {
	Config  json.RawMessage `json:"config"`
	Version int             `json:"version"`
}
