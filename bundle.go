package cobbles

import (
	"io/ioutil"
	"path"
    // "fmt"
    "strings"
    "sort"
)

type Bundle struct {
	dimensions []interface{}
	settings   map[string][]byte
}

type Context struct {
    Settings map[string]string
}

// Loads the filepath and creates a new Bundle for the given configuration interface.
func New(dirpath string) *Bundle {

	this := &Bundle{
        settings: map[string][]byte{},
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
	settings, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
    name := path.Base(filepath)
    key := this.makeLookupPath(name[:strings.Index(name, ".")])
    if _, ok := this.settings[key]; ok {
        panic("Settings group " + key + " has already loaded.")
    }
    this.settings[key] = settings
}

func (this *Bundle) flattenDimensions() {

}

func (this *Bundle) makeLookupPath(context string) string {
    if context == "master" {
        return "master"
    }
    parts := strings.Split(context, ",")
    sort.Strings(parts)
    return strings.Join(parts, ",")
}

// Reads the configuration for the given context into the given configuration interface.
func (this *Bundle) Read(context string, config interface{}) {

}

// Returns the bundle as a YAML byte slice.
func (this *Bundle) Debug() []byte {
	type dump struct {
		Dimensions []interface{}
		Settings   map[string]interface{}
	}
    d := dump{this.dimensions, map[string]interface{}{}}
    for k, v := range this.settings {
        var i interface{}
        fromYaml(v, &i)
        d.Settings[k] = i
    }
	return toYaml(d)
}
