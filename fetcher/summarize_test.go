package fetcher

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func Test_SummarizerAll(t *testing.T) {
	direntry, err := os.ReadDir("test")
	if err != nil {
		t.Fatal(err)
	}

	for _, entry := range direntry {
		addressLabelPath := filepath.Join("test", entry.Name(), "address-label.json")
		profilePath := filepath.Join("test", entry.Name(), "profile.json")

		var profileObj Profile
		var addressLabelObj AddressLabels

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

func Test_FetchAndSummury(t *testing.T) {
	profile, addressLabel, err := FetchTransaction("0x7dd3733b3daa58222376221346f529a8da42f4b8389b639b718c3661e276381d")
	if err != nil {
		t.Fatal(err)
	}

	obj := Summarizer(*profile, *addressLabel)
	jsonString, err := json.Marshal(obj)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(jsonString))
}
