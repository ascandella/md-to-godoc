// Copyright 2016 Aiden Scandella
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"bufio"
	"bytes"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"

	"github.com/sectioneight/md-to-godoc/render"

	"github.com/russross/blackfriday"
)

var (
	inFile      = flag.String("input", "README.md", "Path to markdown file to parse")
	outFile     = flag.String("output", "doc.go", "Path to write file to")
	stdout      = flag.Bool("stdout", false, "Write to STDOUT instead of a file")
	stdin       = flag.Bool("stdin", false, "Read from STDIN instead of a file")
	pkgName     = flag.String("pkg", "", "Package name. If empty, infer from directory of input")
	license     = flag.Bool("license", true, "Add license header from file")
	licenseFile = flag.String("licenseFile", "LICENSE.txt", "File to read license header from")
	badges      = flag.Bool("badges", false, "Enable output for badges (links with images)")

	goListCmd = []string{"list", "-f", "{{.Name}}"}
)

func init() {
	flag.Parse()
}

func main() {
	input, err := ioutil.ReadAll(reader())
	if err != nil {
		log.Fatal("Could not read input file: ", err)
	}

	renderer := render.Godoc(packageName(), *badges)
	output := blackfriday.Markdown(input, renderer, blackfriday.Options{
		Extensions: render.GodocExtensions,
	})

	w := writer()
	defer w.Close()

	if *license {
		if _, err := os.Stat(*licenseFile); err == nil {
			writelicense(w, *licenseFile)
		}
	}
	w.Write(output)
}

func writelicense(w io.Writer, path string) {
	licenseLines := readlicense(*licenseFile)
	for _, line := range licenseLines {
		w.Write([]byte("//"))
		if len(line) > 0 {
			w.Write([]byte(" "))
			w.Write(line)
		}
		w.Write([]byte("\n"))
	}
	w.Write([]byte("\n"))
}

func readlicense(path string) [][]byte {
	fb, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	var lines [][]byte
	bs := bufio.NewScanner(bytes.NewBuffer(fb))
	for bs.Scan() {
		lines = append(lines, bs.Bytes())
	}
	return lines
}

func reader() io.Reader {
	if *stdin {
		return os.Stdin
	}

	f, err := os.Open(*inFile)
	if err != nil {
		panic(err)
	}

	return f
}

func writer() io.WriteCloser {
	if *stdout {
		return os.Stdout
	}

	// Assume they want doc.go to go into the same directory as the input file,
	// Unless they manually set the output.
	inBase := filepath.Dir(*inFile)
	if inBase != "." && *outFile == "doc.go" {
		*outFile = path.Join(inBase, *outFile)
	}

	f, err := os.Create(*outFile)
	if err != nil {
		panic(err)
	}
	return f
}

func packageName() string {
	if *pkgName != "" {
		return *pkgName
	}
	dir := filepath.Dir(*inFile)
	if !filepath.IsAbs(dir) && dir != "." {
		dir = "./" + dir
	}
	cmd := exec.Command("go", append(goListCmd, dir)...)
	output, err := cmd.Output()
	if err != nil {
		log.Printf("Unable to run: %v", cmd.Args)
		log.Println(output)
		panic(err)
	}
	return string(output[:len(output)-1])
}
