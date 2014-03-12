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
			dims := b.flattenDimensions(b.dimensions)
			// fmt.Printf("%s", toYaml(dims))
			AssertEqual("en", dims["lang/en"])
			AssertEqual("en_CA", dims["lang/en/en_CA"])
			AssertEqual("fr", dims["lang/fr"])
			AssertEqual("fr_CA", dims["lang/fr/fr_FR/fr_CA"])
		})
	})

	Report(t)
}
