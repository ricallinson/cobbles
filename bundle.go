package cobbles

import (
	"io/ioutil"
	"path"
	// "fmt"
	"sort"
	"strings"
    // "reflect"
)

type Bundle struct {
	dimensions     []map[string]interface{}
	settings       map[string][]byte
	dimensionPaths map[string]string
}

type Context struct {
	Settings map[string]string
}

// Loads the given directory and returns a new Bundle.
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

// Loads a dimensions YAML file into the Bundle.
func (this *Bundle) loadDimensionsFile(filepath string) {
    if this.dimensions != nil {
        panic("A dimensions file has already been loaded.")
    }
	b, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	this.dimensions = make([]map[string]interface{}, 4)
	fromYaml(b, &this.dimensions)
}

// Loads a settings YAML file into the Bundle.
func (this *Bundle) loadSettingsFile(filepath string) {
	settings, err := ioutil.ReadFile(filepath)
	if err != nil {
		panic(err)
	}
	name := path.Base(filepath)
	key := this.makeLookupPath(name[:strings.Index(name, ".")])
	if _, ok := this.settings[key]; ok {
		panic("A Settings group with the key [" + key + "] has already loaded.")
	}
	this.settings[key] = settings
}

// Returns a slice of ordered lookup strings.
func (this *Bundle) getLookupPaths(context string) []string {
    return []string{}
}

// Takes the unmarshalled YAML "dimensions" slice and returns it as a flattened map.
func (this *Bundle) flattenDimensions(dimensions []map[string]interface{}) map[string]string {

	build := make(map[string]string)

	for _, set := range dimensions {
		for key, dimension := range set {
			this.flattenDimension(key, dimension.(map[interface{}]interface{}), build)
		}
	}

	return build
}

// Takes an unmarshalled YAML "dimensions" map and adds it as a flattened string to the given "build" map.
func (this *Bundle) flattenDimension(prefix string, dimension map[interface{}]interface{}, build map[string]string) {

    for k, nextDimension := range dimension {
        key := k.(string)
        newPrefix := key
        if prefix != "" {
            newPrefix = prefix + "/" + key
        }
        build[newPrefix] = key
        // If nextDimension is a value and is a map.
        if nextDimension != nil {
            this.flattenDimension(newPrefix, nextDimension.(map[interface{}]interface{}), build);
        }
    }
}

// Takes the given context string and returns it as an ordered lookup path.
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
    lookupPaths := this.getLookupPaths(context)
    for _, path := range lookupPaths {
        if yaml, ok := this.settings[path]; ok {
            fromYaml(yaml, config)
        }
    }
}

// Returns the bundle as a YAML byte slice.
func (this *Bundle) Debug() []byte {
	type dump struct {
		Dimensions []map[string]interface{}
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
