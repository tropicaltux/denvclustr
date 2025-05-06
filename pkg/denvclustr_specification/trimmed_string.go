package denvclustr

import (
	"encoding/json"
	"strings"
)

// TrimmedString is a string that is automatically trimmed of whitespace when unmarshaled.
type TrimmedString string

// UnmarshalJSON implements json.Unmarshaler, trimming whitespace from string values.
func (s *TrimmedString) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	*s = TrimmedString(strings.TrimSpace(str))
	return nil
}
