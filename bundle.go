package cobbles

import (
	"io/ioutil"
	"path"
)

type Bundle struct {
	dimensions []interface{}
	settings   []interface{}
    lookup     map[string]int
}

type Context struct {
    Settings map[string]string
}

// Loads the filepath and creates a new Bundle for the given configuration interface.
func New(dirpath string) *Bundle {

	this := &Bundle{
        lookup: map[string]int{},
    }

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

    // For each settings group found in the file.
    for i := range settings {
        setting := settings[i]
        context := &Context{}
        Cast(settings, context)
        key := this.makeLookupPath()
        if _, ok := this.lookup[key]; ok {
            // panic("Settings group " + key + " has already loaded.")
        }
        this.settings = append(this.settings, setting)
        this.lookup[key] = len(this.settings)
    }
}

func (this *Bundle) flattenDimensions() {

}

func (this *Bundle) makeLookupPath() string {
    return ""
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
