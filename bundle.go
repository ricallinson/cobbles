package cobbles

import (
	"io/ioutil"
	"path"
	// "fmt"
	// "sort"
	"strings"
    // "reflect"
)

const(
    CONTEXT_SETTER = "="
    CONTEXT_SEPARATOR = ","
    DEFAULT = "*"
    SEPARATOR = "/"
)

type Bundle struct {
	dimensions     []map[string]interface{}
	settings       map[string][]byte
    dimensionIndex map[string]int
    dimensionPaths map[string]map[string]string
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
	this.dimensions = []map[string]interface{}{}
	fromYaml(b, &this.dimensions)
    this.dimensionIndex, this.dimensionPaths = this.flattenDimensions(this.dimensions)
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

// Takes the unmarshalled YAML "dimensions" slice and returns it as a flattened map.
func (this *Bundle) flattenDimensions(dimensions []map[string]interface{}) (map[string]int, map[string]map[string]string) {

    index := map[string]int{}
	build := make(map[string]map[string]string)

	for i, set := range dimensions {
		for name, dimension := range set {
            index[name] = i
			build[name] = this.flattenDimension("", dimension.(map[interface{}]interface{}))
		}
	}

	return index, build
}

// Takes an unmarshalled YAML "dimensions" map and adds it as a flattened string to the given "build" map.
func (this *Bundle) flattenDimension(prefix string, dimension map[interface{}]interface{}, b... map[string]string) map[string]string {

    var build map[string]string

    if len(b) == 0 {
        build = make(map[string]string)
    } else {
        build = b[0]
    }

    for k, nextDimension := range dimension {
        key := k.(string)
        newPrefix := key
        if prefix != "" {
            newPrefix = prefix + SEPARATOR + key
        }
        build[newPrefix] = key
        // If nextDimension is a value and is a map.
        if nextDimension != nil {
            this.flattenDimension(newPrefix, nextDimension.(map[interface{}]interface{}), build);
        }
    }

    return build
}

// Takes the given context string and returns it as an ordered lookup path.
func (this *Bundle) makeLookupPath(context string) string {
	return context
}

// Returns a slice of ordered lookup strings for the given context.
func (this *Bundle) getLookupPaths(context string) []string {
    this.makeOrderedLookupList(context)
    return []string{}
}

// Returns a slice of ordered lookup strings for the bundles dimensions.
func (this *Bundle) makeOrderedLookupList(context string) map[string][]string {

    list := map[string][]string{}

    for _, set := range this.dimensions {
        for dimensionName, _ := range set {
            // For each dimension value see if have a match.
            for lookupPath, value := range this.dimensionPaths[dimensionName] {
                match := dimensionName + CONTEXT_SETTER + value
                // If there is a match add it to the list.
                if strings.Index(context, match) > -1 {
                    // Reverse the path.
                    slice := strings.Split(DEFAULT + SEPARATOR + lookupPath, SEPARATOR)
                    // sort.Sort(sort.Reverse(sort.StringSlice(slice)))
                    list[dimensionName] = reverseStringSlice(slice)
                }
            }
            // If no match was found use the default.
            if _, ok := list[dimensionName]; ok == false {
                list[dimensionName] = []string{DEFAULT}
            }
        } 
    }

    return list
}

// Replaces any substitutions found in the final configuration.
func (this *Bundle) applySubstitutions(config interface{}) {
    // ...
}

// Reads the configuration for the given context into the given configuration interface.
func (this *Bundle) Read(context string, config interface{}) {
    lookupPaths := this.getLookupPaths(context)
    for _, path := range lookupPaths {
        if yaml, ok := this.settings[path]; ok {
            fromYaml(yaml, config)
        }
    }
    this.applySubstitutions(config)
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
