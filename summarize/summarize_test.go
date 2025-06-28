package summarize

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/hmluck83/txlens-srv/schemas"
)

func Test_SummarizerAll(t *testing.T) {
	direntry, err := os.ReadDir("test")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range direntry {
		addressLabelPath := filepath.Join("test", entry.Name(), "address-label.json")
		profilePath := filepath.Join("test", entry.Name(), "profile.json")

		var profileObj schemas.Profile
		var addressLabelObj schemas.AddressLabels

		addressLabelString, err := os.ReadFile(addressLabelPath)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(addressLabelString, &addressLabelObj)
		if err != nil {
			t.Fatal(err)
		}

		profileString, err := os.ReadFile(profilePath)
		if err != nil {
			t.Fatal(err)
		}

		err = json.Unmarshal(profileString, &profileObj)
		if err != nil {
			t.Fatal(err)
		}

		obj := Summarizer(profileObj, addressLabelObj)
		jsonString, err := json.Marshal(obj)
		if err != nil {
			t.Fatal(err)
		}
		entry.Name()

		os.WriteFile(entry.Name()+".json", jsonString, 0644)

		t.Log(string(jsonString))

	}
}
