package webpacking

import (
	"fmt"
	"io"
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
	"io/ioutil"
)

func fakeExecCommand(command string, args...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}

func TestHelperProcess(t *testing.T){
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	// some code here to check arguments perhaps?

	fmt.Printf(strings.Join(os.Args[3:], " "))
	os.Exit(0)
}

func TestDevelopementMode(t *testing.T) {
	jsHashedAsset := "foo-12345.js"
	cssHashedAsset := "foo-12345.css"
	fake := FakeReadFiler{
		Str: fmt.Sprintf(`{ "foo.js": "%s", "foo.css": "%s" }`, jsHashedAsset, cssHashedAsset),
	}
	execCommand = fakeExecCommand
	defer func() {
		manifestReadFile = ioutil.ReadFile
	}()

	stdInR, stdInW, err := os.Pipe()

	manifestReadFile = fake.ReadFile
	config := &Config{
		IsDev: true,
		Stdout: stdInW,
	}
	wp, err := New(config)
	err = wp.Run()


	if err != nil {
		t.Errorf("run failed: %s", err)
	}

	var buf bytes.Buffer

	go io.Copy(&buf, stdInR)

	err = wp.Process.Wait()
	if err != nil {
		t.Errorf("exit stat: %s", err)
	}

	expectedCmd := "./node_modules/.bin/webpack-dev-server --port 8080 --host localhost"
	if !strings.EqualFold(
		buf.String(),
		expectedCmd,
	) {
		t.Errorf(
			"cmd executed wasn't what was expected.  got:\n%s\nbut expected:\n%s",
			buf.String(), expectedCmd,
		)
	}

}
