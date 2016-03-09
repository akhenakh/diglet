// This is a script to generate go resources from other file formats
package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	pathlib "path"
	"regexp"
	"strings"
)

const ResourceDir = "resources"

// Takes the files from args and creates resources for them in resources/{file}.go
// Also creates a resources/version.go
func main() {
	resources := os.Args[1:]
	os.MkdirAll(ResourceDir, 0777)
	for _, path := range resources {
		contents := bytes.NewBuffer(nil)
		file, _ := os.Open(path)
		io.Copy(contents, file)
		name := pathlib.Base(path)
		rsc := writeResource(name, string(contents.Bytes()))
		logf("Created resource %s", rsc)
	}
	version := version()
	writeResource("Version", version)
	logf("Created resource for Version %s", version)
}

func writeResource(name, contents string) (path string) {
	//dump := base64.StdEncoding.EncodeToString(buf.Bytes())
	export := varSlugged(name)
	local := strings.ToLower(export)
	lines := []string{
		"// THIS FILE IS AUTOGENERATED FROM scripts/include.go EDITS WILL BE SQUISHED",
		"// IT IS MEANT FOR HOLDING CONFIGURABLE OPTIONS",
		"package resources",
		"const (",
		sprintf("%s = `%s`", local, contents),
		")",
		sprintf("func %s() string { return %s }", export, local),
	}
	path = pathlib.Join(ResourceDir, sprintf("%s.go", slugged(name)))
	rsc, err := os.Create(path)
	if err != nil {
		panic(err)
	}
	for _, line := range lines {
		rsc.Write([]byte(sprintf("%s\n", line)))
	}
	return
}

func logf(format string, vals ...interface{}) {
	fmt.Printf(sprintf("include.go: %s\n", format), vals...)
}

func sprintf(format string, vals ...interface{}) string {
	return fmt.Sprintf(format, vals...)
}

var slugger = regexp.MustCompile("[^a-zA-Z0-9]+")

func slugged(s string) string {
	return strings.Trim(slugger.ReplaceAllString(strings.ToLower(s), "_"), "_")
}

// Normal slugs don't make good variable names
// This keeps case and makes an Exported variable
// settings.yml -> Settings_yml
func varSlugged(s string) (vslug string) {
	vslug = strings.Trim(slugger.ReplaceAllString(s, "_"), "_")
	vslug = strings.ToUpper(string(vslug[0])) + vslug[1:]
	return
}

func version() string {
	cmd := exec.Command("git", "describe", "--always")
	var ver bytes.Buffer
	cmd.Stdout = &ver
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
	return strings.TrimSpace(ver.String())
	//TODO add a '{+n}' to version if git diff --numstat isn't empty
}
