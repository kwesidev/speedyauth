package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

// JSONB Interface
type JSONB map[string]interface{}

// Value Marshal
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan Unmarshal
func (j *JSONB) Scan(value interface{}) error {
	source, ok := value.([]byte)
	if !ok {
		return errors.New("Type assertion .([]byte) failed.")
	}

	var i interface{}
	err := json.Unmarshal(source, &i)
	if err != nil {
		return err
	}

	*j, ok = i.(map[string]interface{})
	if !ok {
		return errors.New("Type assertion .(map[string]interface{}) failed.")
	}
	return nil
}
