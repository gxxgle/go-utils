package json

import (
	"github.com/json-iterator/go"
)

// global variable
var (
	JSON = jsoniter.ConfigCompatibleWithStandardLibrary
)

// UseNumber solve very big int64 digits loss.
func UseNumber() {
	JSON = jsoniter.Config{
		UseNumber:              true,
		EscapeHTML:             true,
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
	}.Froze()
}

// Marshal returns the JSON encoding of v.
func Marshal(v interface{}) ([]byte, error) {
	return JSON.Marshal(v)
}

// MustMarshal must returns the JSON encoding of v.
func MustMarshal(v interface{}) []byte {
	data, _ := JSON.Marshal(v)
	return data
}

// MarshalToString returns the JSON encoding to string of v.
func MarshalToString(v interface{}) (string, error) {
	return JSON.MarshalToString(v)
}

// MustMarshalToString must returns the JSON encoding to string of v.
func MustMarshalToString(v interface{}) string {
	str, _ := JSON.MarshalToString(v)
	return str
}

// Unmarshal parses the JSON-encoded data and stores the result
// in the value pointed to by v.
func Unmarshal(data []byte, v interface{}) error {
	return JSON.Unmarshal(data, v)
}

// UnmarshalFromString unmarshal string to v.
func UnmarshalFromString(str string, v interface{}) error {
	return JSON.UnmarshalFromString(str, v)
}

// Valid check JSON data.
func Valid(data []byte) bool {
	return JSON.Valid(data)
}

// ValidFromString check JSON string.
func ValidFromString(str string) bool {
	return Valid([]byte(str))
}

// Get get value from JSON data by path.
func Get(data []byte, path ...interface{}) jsoniter.Any {
	return JSON.Get(data, path...)
}

// GetFromString get value from JSON string by path.
func GetFromString(str string, path ...interface{}) jsoniter.Any {
	return Get([]byte(str), path...)
}
