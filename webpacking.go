package webpacking

import (
	"io"
	"os"
	"fmt"
"github.com/davecgh/go-spew/spew"
	"os/exec"
	"strings"
	"html/template"

)

// configuration for webpacking, is passed to the Init
type Config struct {
	// host for webpack-dev-server, defaults to localhost
	DevHost string
	// port for webpack-dev-server, defaults to 8080
	DevPort string
	// FsPath filesystem path to public webpack dir
	ManifestPath string
	// WebPath http path to public webpack dir
	WebPath string
	// path to node_modules where webpack is installed
	// defaults to './node_modules'
	NodeModulesPath string
	// destination to copy stdout to.  Defaults to os.Stdout
	Stdout *os.File
	// destination to copy stderr to.  Defaults to os.Stderr
	Stderr *os.File
	// Verbose - show more info
	Verbose bool
	// IsDev - true to use webpack-dev-server
	// false to use filesystem and manifest.json
	IsDev bool
}


type WebPacking struct {
	config *Config
	Assets map[string]string
	Process *exec.Cmd
}

var execCommand = exec.Command

func (wp *WebPacking) Run() error {
	if !wp.config.IsDev {
		manifest, err := ReadManifest(wp.config)
		wp.Assets = manifest
		return err
	}

	wp.Process = execCommand(
		fmt.Sprintf("%s/.bin/webpack-dev-server", wp.config.NodeModulesPath),
		"--port", wp.config.DevPort,
		"--host", wp.config.DevHost,
	)

	stdoutIn, _ := wp.Process.StdoutPipe()
	stderrIn, _ := wp.Process.StderrPipe()

	// Start command
	if err := wp.Process.Start(); err != nil {
		return err
	}

	go io.Copy(wp.config.Stderr, stderrIn)
	go io.Copy(wp.config.Stdout, stdoutIn)

	return nil
}


// finds an hashed asset filename from the manifest file and returns it
func (wp *WebPacking) GetAsset(asset string) (string, error) {
	if wp.Assets == nil || wp.config.IsDev {
		manifest, err := ReadManifest(wp.config)
		if err != nil {
			return "", err
		}
		wp.Assets = manifest
	}
	found, ok := wp.Assets[asset]
	if ok {
		return found, nil
	}
	spew.Dump(wp.Assets)
	return "", fmt.Errorf("asset %s was not found", asset)
}


func (wp *WebPacking) AssetHelper() func(string) (template.HTML, error) {
	return func(key string) (template.HTML, error) {
		asset, err := wp.GetAsset(key)
		if err != nil {
			return "", err
		}

		parts := strings.Split(asset, ".")
		kind := parts[len(parts)-1]

		return template.HTML(AssetTag(kind, asset)), nil
	}
}


// initialize and return a WebPacking instance
func New(config *Config) (*WebPacking, error) {
	if config.NodeModulesPath == "" {
		config.NodeModulesPath = "./node_modules"
	}
	if config.DevHost == "" {
		config.DevHost = "localhost"
	}
	if config.DevPort == "" {
		config.DevPort = "8080"
	}
	if config.ManifestPath == "" {
		config.ManifestPath = "./public/assets"
	}
	if config.Stdout == nil {
		config.Stdout = os.Stdout
	}
	if config.Stderr == nil {
		config.Stderr = os.Stderr
	}
	packer := WebPacking {
		config: config,
	}
	return &packer, nil
}
