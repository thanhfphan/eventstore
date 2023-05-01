package dtype

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/thanhfphan/eventstore/pkg/errors"
)

// JSON custom type to mapping from jsonb of postgres
type JSON map[string]interface{}

// Scan implements the sql.Scanner interface.
func (m *JSON) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to unmarshal JSONB value: %v", value)
	}

	return json.Unmarshal(bytes, m)
}

// Value implements the driver.Valuer interface.
func (m JSON) Value() (driver.Value, error) {
	b, err := json.Marshal(m)
	if len(b) == 0 {
		return "{}", nil
	}
	return string(b), err
}

// ToJSON convert struct to JSON
func ToJSON(obj interface{}) (JSON, error) {
	if obj == nil {
		return JSON{}, nil
	}

	jsonBytes, err := json.Marshal(obj)
	if err != nil {
		return nil, errors.New("Convert struct to json failed with err=%v", err)
	}

	var result JSON
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, errors.New("Convert json data to JSON failed with err=%v", err)
	}

	return result, nil
}
