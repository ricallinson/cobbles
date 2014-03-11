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

func TestRouter(t *testing.T) {

	Describe("", func() {
		It("", func() {
			New("./fixtures")
			// fmt.Printf("%v", string(b.Debug()))
			AssertEqual(true, false)
		})
	})

	Report(t)
}
