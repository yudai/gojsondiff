package tests

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func LoadFixture(t *testing.T, file string) interface{} {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		t.Errorf("Fixture file `%s` not found.", file)
	}
	var result interface{}
	err = json.Unmarshal(content, &result)
	if err != nil {
		t.Errorf("Unmarshaling JSON of `%s` failed: %s", file, err)
	}
	return result
}
