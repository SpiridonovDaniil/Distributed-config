package helpers

import (
	"encoding/json"
	"fmt"
	"github.com/SpiridonovDaniil/Distributed-config/internal/domain"
)

func Converter(config *domain.Request) ([]byte, error) {
	var metadata []map[string]interface{}
	err := json.Unmarshal(config.Data, &metadata)
	if err != nil {
		err = fmt.Errorf("[converter] failed to unmarshal config.Data, error: %w", err)
		return nil, err
	}

	resultMeta := make(map[string]interface{})
	for _, data := range metadata {
		for key, val := range data {
			resultMeta[key] = val
		}
	}

	rawData, _ := json.Marshal(resultMeta)

	return rawData, nil
}
