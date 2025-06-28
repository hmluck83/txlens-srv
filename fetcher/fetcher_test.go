package fetcher

import (
	"encoding/json"
	"testing"
)

func TestFetcher(t *testing.T) {
	profile, addressLabel, err := FetchTransaction("0x7dd3733b3daa58222376221346f529a8da42f4b8389b639b718c3661e276381d")
	if err != nil {
		t.Fatal(err)
	}
	data, err := json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))
	data, err = json.Marshal(addressLabel)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(string(data))

}
