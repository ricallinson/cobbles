package cobbles

import (
	"io/ioutil"
	"path"
)

type Bundle struct {
	dimensions []interface{}
	settings   []interface{}
}

// Loads the filepath and creates a new Bundle for the given configuration interface.
func New(dirpath string) *Bundle {

	this := &Bundle{}

	files, err := ioutil.ReadDir(dirpath)
	if err != nil {
		panic(err)
	}

	for f := range files {
		file := files[f]
		if file.IsDir() == false {
			filepath := path.Join(dirpath, file.Name())
			if path.Base(filepath) == "dimensions.yaml" {
				this.loadDimensionsFile(filepath)
			} else if path.Ext(filepath) == ".yaml" {
				this.loadSettingsFile(filepath)
			}
		}
	}

	return this
}

// Loads a dimensions YAML file.
func (this *Bundle) loadDimensionsFile(filepath string) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	dimensions := []interface{}{}
	fromYaml(b, &dimensions)
	this.dimensions = append(this.dimensions, dimensions...)
}

// Loads a settings YAML file.
func (this *Bundle) loadSettingsFile(filepath string) {
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	settings := []interface{}{}
	fromYaml(b, &settings)
	this.settings = append(this.settings, settings...)
}

// Reads the configuration for the given context into the given configuration interface.
func (this *Bundle) Read(context string, config interface{}) {

}

// Returns the bundle as a YAML byte slice.
func (this *Bundle) Debug() []byte {
	type dump struct {
		Dimensions []interface{}
		Settings   []interface{}
	}
	return toYaml(dump{this.dimensions, this.settings})
}