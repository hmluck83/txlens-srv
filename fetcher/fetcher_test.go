package fetcher

import (
	"encoding/json"
	"os"
	"path/filepath"
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
	t.Log("---- PROFILE -----")
	t.Log(string(data))
	data, err = json.Marshal(addressLabel)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("---- Address Label -----")

	t.Log(string(data))
}

func TestCreateSample(t *testing.T) {
	caseName := "ChaingenowDeposit_1"
	caseTx := "0x81a2341ca06e1b72ea35cd812380e2cf8d312ce79c1827a76c18623d31638237"
	if err := os.Mkdir(filepath.Join("test", caseName), os.ModePerm); err != nil {
		t.Fatal(err)
	}

	profile, addressLabel, err := FetchTransaction(caseTx)
	if err != nil {
		t.Fatal(err)
	}
	t.Log((profile))

	profilestring, err := json.Marshal(profile)
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile(filepath.Join("test", caseName, "profile.json"), profilestring, 0644)

	addressObj, err := json.Marshal(addressLabel)
	if err != nil {
		t.Fatal(err)
	}

	os.WriteFile(filepath.Join("test", caseName, "address-label.json"), addressObj, 0644)
}
