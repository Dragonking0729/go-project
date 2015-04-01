//go:generate go run defaults.go default.json defs.go

package main //build !none

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
)

func fatal(str string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, str, v...)
	os.Exit(1)
}

type setting struct {
	Value   int64  `json:"v"`
	Comment string `json:"d"`
}

func main() {
	if len(os.Args) < 3 {
		fatal("usage %s <input> <output>\n", os.Args[0])
	}

	content, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		fatal("error reading file %v\n", err)
	}

	m := make(map[string]setting)
	json.Unmarshal(content, &m)

	filepath := path.Join(os.Getenv("GOPATH"), "src", "github.com", "ethereum", "go-ethereum", "core", os.Args[2])
	output, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm /*0777*/)
	if err != nil {
		fatal("error opening file for writing %v\n", err)
	}

	output.WriteString(`package core

import "math/big"

var (
`)

	for name, setting := range m {
		output.WriteString(fmt.Sprintf("%s=big.NewInt(%d) // %s\n", strings.Title(name), setting.Value, setting.Comment))
	}

	output.WriteString(")\n")
	output.Close()

	cmd := exec.Command("gofmt", "-w", filepath)
	if err := cmd.Run(); err != nil {
		fatal("gofmt failed: %v\n", err)
	}
}
