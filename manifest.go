package webpacking

import (
	"io/ioutil"
	"encoding/json"
	"github.com/pkg/errors"
)

var manifestReadFile = ioutil.ReadFile

type FakeAssetReader struct {
	Contents map[string]string
}

func (f FakeAssetReader) ReadFile(filename string) ([]byte, error) {
	return json.Marshal(f.Contents)
}

func (f FakeAssetReader) Restore() {
	manifestReadFile = ioutil.ReadFile
}

func InstallFakeAssetReader() *FakeAssetReader {
	fake := FakeAssetReader{}
	fake.Contents = make(map[string]string)
	fake.Contents["testing.js"] =  "testing-12345.js"
	fake.Contents["testing.css"] =  "testing-12345.css"

	manifestReadFile = fake.ReadFile
	return &fake
}


// Read a webpack-manifest-plugin format manifest file stored at path
// It returns a map with string key/values for the contents of the map
func ReadManifest(config *Config) (map[string]string, error) {

	data, err := manifestReadFile(config.ManifestPath + "/manifest.json")

	if err != nil {
		return nil, errors.Wrap(err, "webpacking: Error when loading manifest from file")
	}

	response := make(map[string]string)
	err = json.Unmarshal(data, &response)
	return response, err
}
