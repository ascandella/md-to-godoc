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

// Package render contains the godoc render logic, which may be useful to others.
package render

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/russross/blackfriday"
)

// GodocExtensions are the default markdown extensions for blackfriday
const GodocExtensions = blackfriday.CommonExtensions

var (
	nl       = []byte("\n")
	indent   = []byte("  ")
	star     = []byte("*")
	starstar = []byte("**")
)

// Godoc returns a blackfriday renderer for doc.go style package documentation.
func Godoc(pkg string) blackfriday.Renderer {
	return &GodocRenderer{
		pkg: pkg,
	}
}

// GodocRenderer implements the blackfriday.Render interface for doc.go style
// package documentation
type GodocRenderer struct {
	pkg              string
	pkgHeaderWritten bool
	lastOutputLen    int
	inImage          bool
}

// Render walks the specified (sub)tree and returns a godoc document.
func (g *GodocRenderer) Render(ast *blackfriday.Node) []byte {
	var buff bytes.Buffer
	g.DocumentHeader(&buff)

	ast.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		return g.RenderNode(&buff, node, entering)
	})

	g.DocumentFooter(&buff)
	return buff.Bytes()
}

// RenderNode is a default renderer of a single node of a syntax tree. For
// block nodes it will be called twice: first time with entering=true, second
// time with entering=false, so that it could know when it's working on an open
// tag and when on close. It writes the result to w.
//
// The return value is a way to tell the calling walker to adjust its walk
// pattern: e.g. it can terminate the traversal by returning Terminate. Or it
// can ask the walker to skip a subtree of this node by returning SkipChildren.
// The typical behavior is to return GoToNext, which asks for the usual
// traversal to the next node.
func (g *GodocRenderer) RenderNode(w io.Writer, node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
	// you don't know, til you know
	debug := os.Getenv("DEBUG_IT") == "yass"
	if debug {
		log.Printf("Type: %+v Val: |%+v|, /%+v/\n", node.Type, string(node.Literal), node)
	}

	switch node.Type {
	case blackfriday.Text:
		g.out(w, node.Literal)
	case blackfriday.Softbreak:
		g.cr(w)
	case blackfriday.Hardbreak:
		g.cr(w)
		g.cr(w)

	case blackfriday.Header:
		if entering {
			g.out(w, node.Literal)
		} else {
			if !g.pkgHeaderWritten {
				g.out(w, []byte("."))
				g.pkgHeaderWritten = true
			}
			g.cr(w)
			g.cr(w)
		}

	case blackfriday.Paragraph:
		if entering {
			g.out(w, node.Literal)
		} else {
			g.cr(w)
			g.cr(w)
		}

	case blackfriday.Document:
		break
	case blackfriday.List:
		g.out(w, nl)

	case blackfriday.Link:
		if g.inImage {
			if debug {
				log.Println("Skip children for image")
			}
			return blackfriday.SkipChildren
		}
		if entering {
			if len(node.LinkData.Title) > 0 {
				g.out(w, node.LinkData.Title)
			}
		} else {
			dest := node.LinkData.Destination
			g.out(w, []byte(" ("))
			g.out(w, dest)
			g.out(w, []byte(")"))
		}

	case blackfriday.Emph:
		g.out(w, star)
	case blackfriday.Strong:
		g.out(w, starstar)

	case blackfriday.Image:
		if debug {
			fmt.Printf("Setting inImage to %v\n", entering)
		}
		g.inImage = entering
		// nope
		return blackfriday.SkipChildren

	case blackfriday.Item:
		if entering {
			if node.Prev != nil {
				g.cr(w)
			}
			g.out(w, []byte("• "))
		} else {
			g.out(w, node.Literal)
			g.cr(w)
		}

	case blackfriday.CodeBlock:
		g.blockCode(w, node.Literal, string(node.Info))

	case blackfriday.Code:
		// Sadly, no inline code support or emphasis
		g.out(w, node.Literal)

	case blackfriday.Table:
		// unsupported, do nothing
	case blackfriday.TableCell, blackfriday.TableRow, blackfriday.TableBody, blackfriday.TableHead:
		// unsupported, do nothing

	default:
		panic("Unknown node type " + node.Type.String())
	}

	return blackfriday.GoToNext
}

func (g *GodocRenderer) out(w io.Writer, text []byte) {
	w.Write(text)
	g.lastOutputLen = len(text)
}

func (g *GodocRenderer) cr(w io.Writer) {
	if g.lastOutputLen > 0 {
		g.out(w, nl)
	}
}

func (g *GodocRenderer) blockCode(out io.Writer, text []byte, lang string) {
	s := bufio.NewScanner(bytes.NewBuffer(text))
	for s.Scan() {
		b := s.Bytes()
		if len(b) > 0 {
			out.Write(indent)
			out.Write(s.Bytes())
		}
		out.Write(nl)
	}
}

// DocumentHeader writes the beginning of the package documentation.
func (g *GodocRenderer) DocumentHeader(out *bytes.Buffer) {
	out.WriteString("/*\n")
	out.WriteString("Package " + g.pkg + " is the ")
}

// DocumentFooter writes the end of the package documentation
func (g *GodocRenderer) DocumentFooter(out *bytes.Buffer) {
	out.WriteString("*/\npackage " + g.pkg + "\n")
}
