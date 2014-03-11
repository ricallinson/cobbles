package cobbles

import (
	"launchpad.net/goyaml"
)

// Return the given interface as a YAML byte slice.
func toYaml(i interface{}) []byte {
	data, err := goyaml.Marshal(i)
	if err != nil {
		panic(err)
	}
	return data
}

// Reads the given YAML byte slice into the given interface.
func fromYaml(yaml []byte, i interface{}) {
	err := goyaml.Unmarshal(yaml, i)
	if err != nil {
		panic(err)
	}
}