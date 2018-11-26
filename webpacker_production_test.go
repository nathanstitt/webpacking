package webpacking

import (
	"fmt"
	"bytes"
	"strings"
	"testing"
	"io/ioutil"
)


type FakeReadFiler struct {
	Str string
}

func (f FakeReadFiler) ReadFile(filename string) ([]byte, error) {
	buf := bytes.NewBufferString(f.Str)
	return ioutil.ReadAll(buf)
}

func TestProductionMode(t *testing.T) {
	jsHashedAsset := "foo-12345.js"
	cssHashedAsset := "foo-12345.css"
	fake := FakeReadFiler{
		Str: fmt.Sprintf(`{ "foo.js": "%s", "foo.css": "%s" }`, jsHashedAsset, cssHashedAsset),
	}
	defer func() {
		manifestReadFile = ioutil.ReadFile
	}()
	wp, err := New(&Config{
		IsDev: false,
	})
	err = wp.Run()
	if err == nil {
		t.Error("run success but it should have failed")
	}

	manifestReadFile = fake.ReadFile
	err = wp.Run()
	if err != nil {
		t.Errorf("run failed: %v", err)
	}

	asset, err := wp.GetAsset("foo.js")
	if err != nil {
		t.Errorf("asset find failed: %v", err)
	}
	if !strings.EqualFold(asset, jsHashedAsset) {
		t.Errorf("asset didn't return correct string, should have been %s but was %s", asset,
			jsHashedAsset)
	}

	helper := wp.AssetHelper()
	jsTag, err := helper("foo.js")
	if err != nil {
		t.Errorf("asset helper find failed: %v", err)
	}

	jsExpectedTag := `<script type="text/javascript" src="foo-12345.js"></script>`
	if !strings.EqualFold(
		string(jsTag),
		jsExpectedTag,
	) {
		t.Errorf(
			"asset didn't return correct string, should have been %s but was %s",
			jsExpectedTag, jsTag,
		)
	}

	cssTag, err := helper("foo.css")
	if err != nil {
		t.Errorf("asset helper find failed: %v", err)
	}
	cssExpectedTag := `<link type="text/css" rel="stylesheet" href="foo-12345.css"></link>`
	if !strings.EqualFold(
		string(cssTag),
		cssExpectedTag,
	) {
		t.Errorf(
			"asset didn't return correct string, should have been %s but was %s",
			cssExpectedTag, cssTag,
		)
	}


}
