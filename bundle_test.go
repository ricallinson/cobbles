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

	Describe("Bundle.stringToContext()", func() {
		It("should return a map", func() {
			b := &Bundle{}
			c := b.stringToContext("lang=fr_CA,region=ir,environment=staging")
			AssertEqual(c["lang"], "fr_CA")
			AssertEqual(c["region"], "ir")
			AssertEqual(c["environment"], "staging")
		})
	})

	Describe("Bundle.flattenDimensions()", func() {
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

	Describe("Bundle.makeOrderedLookupList()", func() {
		It("should return an ordered lookup list", func() {
			b := New("./fixtures")
			list := b.makeOrderedLookupList(map[string]string{"lang": "fr_CA", "region": "ir", "environment": "staging"})
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

	Describe("Bundle.makeLookupPath()", func() {
		It("should return the lookup path for the given context", func() {
			b := New("./fixtures")
			p := b.makeLookupPath(map[string]string{"lang": "fr_CA", "region": "ir", "environment": "staging"})
			AssertEqual(p, "staging/fr_CA/ir/*")
			// fmt.Printf("%s\n", b.Debug())
		})
	})

	Report(t)
}
