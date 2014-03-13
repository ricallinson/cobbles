package cobbles

import (
	"github.com/go-yaml/yaml"
)

// Return the given interface as a YAML byte slice.
func toYaml(i interface{}) []byte {
	data, err := yaml.Marshal(i)
	if err != nil {
		panic(err)
	}
	return data
}

// Reads the given YAML byte slice into the given interface.
func fromYaml(b []byte, i interface{}) {
	err := yaml.Unmarshal(b, i)
	if err != nil {
		panic(err)
	}
}

// As it says.
func reverseStringSlice(in []string) []string {
    size := len(in)
    out := make([]string, size)
    size--
    for i, v := range in {
        out[size - i] = v
    }
    return out
}