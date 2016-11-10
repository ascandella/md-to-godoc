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

package render

import (
	"bytes"
	"testing"

	"github.com/russross/blackfriday"
	"github.com/stretchr/testify/assert"
)

func TestGodocCTor(t *testing.T) {
	g := Godoc("mypkg")
	assert.Equal(t, "mypkg", g.(*GodocRenderer).pkg)
}

func TestRender_OK(t *testing.T) {
	g := Godoc("mypkg")
	ast := &blackfriday.Node{
		Type:    blackfriday.Text,
		Literal: []byte("hello"),
	}

	bs := g.Render(ast)
	assert.Equal(t,
		"/*\nPackage mypkg is the hello*/\npackage mypkg\n",
		string(bs),
	)
}

func TestDocumentHeader(t *testing.T) {
	out := &bytes.Buffer{}
	g := &GodocRenderer{
		pkg: "fun",
	}
	g.DocumentHeader(out)
	assert.Equal(t, out.String(), "/*\nPackage fun is the ")
}

func TestDocumentFooter(t *testing.T) {
	out := &bytes.Buffer{}
	g := &GodocRenderer{
		pkg: "fun",
	}
	g.DocumentFooter(out)
	assert.Equal(t, out.String(), "*/\npackage fun\n")
}

func TestBlockCode(t *testing.T) {
	buff := &bytes.Buffer{}
	code := []byte(`fmt.Println("Hello, world")`)
	g := &GodocRenderer{}
	g.blockCode(buff, code, "go")

	assert.Equal(t, buff.String(), `  fmt.Println("Hello, world")`)
}
