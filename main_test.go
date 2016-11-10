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
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader_Stdin(t *testing.T) {
	defer overrideBool(stdin, true)()

	assert.Equal(t, os.Stdin, reader())
}

func TestReader_BadFile(t *testing.T) {
	defer overrideString(inFile, "non-existent")()

	assert.Panics(t, func() {
		reader()
	})
}

func TestReader_GoodFile(t *testing.T) {
	r := reader()

	assert.NotEqual(t, os.Stdin, r)
}

func TestWriter_Stdout(t *testing.T) {
	defer overrideBool(stdout, true)()

	assert.Equal(t, os.Stdout, writer())
}

func TestPackageName_Overridden(t *testing.T) {
	defer overrideString(pkgName, "seattle")()

	assert.Equal(t, "seattle", packageName())
}

func TestPackageName_Inferred(t *testing.T) {
	assert.Equal(t, "main", packageName())
}

func TestPackageName_BadDirPanics(t *testing.T) {
	defer overrideString(inFile, "/foobar")()

	assert.Panics(t, func() {
		packageName()
	})
}

func TestPackageName_Relative(t *testing.T) {
	defer overrideString(inFile, "render/README.md")()

	assert.Equal(t, "render", packageName())
}

func TestMain_OK(t *testing.T) {
	assert.NotPanics(t, main)
}

func TestMain_PanicsBadReader(t *testing.T) {
	defer overrideString(inFile, "bad-text-file")()
	assert.Panics(t, main)
}

func TestWriter_CustomFile(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "md-to-godoc")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())
	defer overrideString(outFile, tmpFile.Name())()

	assert.NotPanics(t, func() {
		reader()
	})
}

func TestWritelicense_OK(t *testing.T) {
	defer overrideBool(license, true)()

	require.NotPanics(t, main)
	contents, err := ioutil.ReadFile("doc.go")
	require.NoError(t, err)

	assert.Contains(t, string(contents), "Copyright 2016")
}

func overrideBool(target *bool, val bool) func() {
	old := *target
	*target = val
	return func() {
		*target = old
	}
}

func overrideString(target *string, val string) func() {
	old := *target
	*target = val
	return func() {
		*target = old
	}
}
