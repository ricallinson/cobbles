package cobbles

import (
	// "fmt"
	. "github.com/ricallinson/simplebdd"
	"testing"
)

type SimpleConfig struct {
	Settings  map[string]string
	Title_key string
	Data_url  string
	Logo      string
	Links     map[string]string
}

func TestBundle(t *testing.T) {

	Describe("flattenDimensions", func() {
		It("should return a flattened dimensions map", func() {
			b := New("./fixtures")
			index, dims := b.flattenDimensions(b.dimensions)
			// fmt.Printf("%s", toYaml(b.dimensions))
			AssertEqual(index["environment"], 0)
			AssertEqual(index["lang"], 1)
			AssertEqual(index["region"], 2)
			AssertEqual(index["flavor"], 3)
			// fmt.Printf("%s", toYaml(dims))
			AssertEqual(dims["lang"]["en"], "en")
			AssertEqual(dims["lang"]["en/en_CA"], "en_CA")
			AssertEqual(dims["lang"]["fr"], "fr")
			AssertEqual(dims["lang"]["fr/fr_FR/fr_CA"], "fr_CA")
		})
	})

	Describe("makeOrderedLookupList", func() {
		It("should return an ordered lookup list", func() {
			b := New("./fixtures")
			list := b.makeOrderedLookupList("lang=fr_CA,region=ir,environment=staging")
			// fmt.Printf("%s", toYaml(list))
			AssertEqual(list["environment"][0], "staging")
			AssertEqual(list["environment"][1], "*")
			AssertEqual(list["flavor"][0], "*")
			AssertEqual(list["lang"][0], "fr_CA")
			AssertEqual(list["lang"][1], "fr_FR")
			AssertEqual(list["lang"][2], "fr")
			AssertEqual(list["lang"][3], "*")
			AssertEqual(list["region"][0], "ir")
			AssertEqual(list["region"][1], "gb")
			AssertEqual(list["region"][2], "europe")
			AssertEqual(list["region"][3], "*")
		})
	})

	Describe("makeLookupPath", func() {
		It("should return a", func() {
			b := &Bundle{}
			p := b.makeLookupPath("lang=fr_CA,region=ir,environment=staging")
			AssertEqual(p, "fr_CA")
		})
	})

	Report(t)
}
