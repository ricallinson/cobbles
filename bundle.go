package cobbles

import (
	"io/ioutil"
	"path"
	// "fmt"
	"strings"
)

const (
	CONTEXT_SETTER    = "="
	CONTEXT_SEPARATOR = ","
	DEFAULT           = "*"
	SEPARATOR         = "/"
)

type Bundle struct {
	dimensions     []map[string]interface{}
	settings       map[string][]byte
	dimensionIndex map[string]int
	dimensionPaths map[string]map[string]string
}

// Used in the tumble function.
type combi struct {
	current int
	total   int
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

// String to context map.
func (this *Bundle) stringToContext(contextString string) map[string]string {
	context := map[string]string{}
	parts := strings.Split(contextString, CONTEXT_SEPARATOR)
	for _, part := range parts {
		bits := strings.Split(part, CONTEXT_SETTER)
		if len(bits) == 2 {
			context[bits[0]] = bits[1]
		}
	}
	return context
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
	context := this.stringToContext(name[:strings.Index(name, ".")])
	key := this.makeLookupPath(context)
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
func (this *Bundle) flattenDimension(prefix string, dimension map[interface{}]interface{}, b ...map[string]string) map[string]string {

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
			this.flattenDimension(newPrefix, nextDimension.(map[interface{}]interface{}), build)
		}
	}

	return build
}

// Returns a slice of ordered lookup strings for the bundles dimensions.
func (this *Bundle) makeOrderedLookupList(context map[string]string) map[string][]string {

	list := map[string][]string{}

	for _, set := range this.dimensions {
		for dimensionName, _ := range set {
			// For each dimension value see if have a match.
			for lookupPath, match := range this.dimensionPaths[dimensionName] {
				// If there is a match add it to the list.
				if _, ok := context[dimensionName]; ok && match == context[dimensionName] {
					// Reverse the path.
					slice := strings.Split(DEFAULT+SEPARATOR+lookupPath, SEPARATOR)
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

// Tumbles over a slice of combi objects.
func (this *Bundle) tumble(combination []combi, pos int) bool {
	// If the position is not found return.
	if pos < 0 {
		return false
	}
	// Move along to the next item.
	combination[pos].current += 1
	// If the next item is not found move to the prev position.
	if combination[pos].current > combination[pos].total {
		combination[pos].current = 0
		return this.tumble(combination, pos-1)
	}
	return true
}

// Returns an ordered slice of lookup paths for the given context.
func (this *Bundle) makeLookupPaths(context map[string]string) []string {

	paths := []string{}
	values := this.makeOrderedLookupList(context)
	combination := make([]combi, len(this.dimensionIndex))
	startPos := len(combination) - 1

	for dimensionName, pos := range this.dimensionIndex {
		combination[pos] = combi{
			current: 0,
			total:   len(values[dimensionName]) - 1,
		}
	}

	for {
		path := []string{}
		for dimensionName, pos := range this.dimensionIndex {
			path = append(path, values[dimensionName][combination[pos].current])
		}
		paths = append(paths, strings.Join(path, SEPARATOR))
		if this.tumble(combination, startPos) == false {
			return reverseStringSlice(paths)
		}
	}

	// If got here something went very wrong.
	return []string{}
}

// Takes the given context and returns its lookup path.
func (this *Bundle) makeLookupPath(context map[string]string) string {

	lookup := map[string]string{}
	path := []string{}
	lookupList := this.makeOrderedLookupList(context)

	for dimensionName, _ := range this.dimensionIndex {
		if match, ok := context[dimensionName]; ok {
			for _, lookupName := range lookupList[dimensionName] {
				if match == lookupName {
					lookup[dimensionName] = lookupName
				}
			}
		}
		if _, ok := lookup[dimensionName]; ok == false {
			lookup[dimensionName] = DEFAULT
		}
	}

	for _, item := range lookup {
		path = append(path, item)
	}

	return strings.Join(path, SEPARATOR)
}

// Reads the configuration for the given context into the given configuration interface.
func (this *Bundle) Read(config interface{}, c ...string) {
	context := map[string]string{}
	if len(c) == 1 {
		context = this.stringToContext(c[0])
	}
	lookupPaths := this.makeLookupPaths(context)
	for _, path := range lookupPaths {
		if yaml, ok := this.settings[path]; ok {
			fromYaml(yaml, config)
		}
	}
}
