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
	"io"
	"log"
	"os"

	"github.com/russross/blackfriday"
)

// GodocExtensions are the default markdown extensions for blackfriday
const GodocExtensions = blackfriday.CommonExtensions

var (
	nl         = []byte("\n")
	indent     = []byte("  ")
	space      = []byte(" ")
	star       = []byte("*")
	starstar   = []byte("**")
	slashslash = []byte("//")
)

// Godoc returns a blackfriday renderer for doc.go style package documentation.
func Godoc(pkg string, badges bool) blackfriday.Renderer {
	return &GodocRenderer{
		pkg:     pkg,
		noBadge: !badges,
	}
}

// GodocRenderer implements the blackfriday.Render interface for doc.go style
// package documentation
type GodocRenderer struct {
	pkg     string
	noBadge bool

	pkgHeaderWritten bool
	lastOutputLen    int
	imageInLink      bool
	inLink           bool
	newline          bool
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
	debug := os.Getenv("DEBUG_IT") != ""
	if debug {
		log.Printf("Type: %+v Val: |%+v|, /%+v/ %v\n", node.Type, string(node.Literal), node, entering)
	}

	switch node.Type {
	case blackfriday.Text:
		if g.inLink && g.imageInLink && g.noBadge {
			if debug {
				log.Println("Skipping text node because we're in a badge")
			}
			return blackfriday.GoToNext
		}
		lines := bytes.Split(node.Literal, nl)
		for _, line := range lines {
			// Trim off trailing space for OCD
			if len(line) > 0 && string(line[len(line)-1]) == " " {
				line = line[0 : len(line)-1]
			}
			g.out(w, line)
			if len(lines) > 1 {
				g.cr(w)
			}
		}

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
		if !entering {
			g.cr(w)
		}

	case blackfriday.Link:
		g.inLink = entering
		if g.imageInLink && g.noBadge {
			if debug {
				log.Println("Skipping link due to badge skip flag")
			}
			g.imageInLink = false
			return blackfriday.GoToNext
		}
		if entering {
			if len(node.LinkData.Title) > 0 {
				g.out(w, node.LinkData.Title)
			}
		} else {
			dest := node.LinkData.Destination
			// Reset this for badge detection
			g.imageInLink = false
			g.out(w, []byte(" ("))
			g.out(w, dest)
			g.out(w, []byte(")"))
		}

	case blackfriday.Emph:
		if entering {
			g.out(w, space)
		}
		g.out(w, star)
	case blackfriday.Strong:
		if entering {
			g.out(w, space)
		}
		g.out(w, starstar)

	case blackfriday.Image:
		if entering && g.inLink {
			// There's an image inside a link, this is most likely a badge
			// (e.g. Travis/Coveralls. Ignore it altogether)
			g.imageInLink = true
		}
		// no support for images

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
	if g.newline && len(text) > 0 && string(text) != "//" && string(text) != "\n" {
		w.Write(space)
		g.newline = false
	}
	w.Write(text)
	g.lastOutputLen = len(text)
}

func (g *GodocRenderer) cr(w io.Writer) {
	if g.lastOutputLen > 0 {
		g.out(w, nl)
		g.out(w, slashslash)
		g.newline = true
	}
}

func (g *GodocRenderer) blockCode(out io.Writer, text []byte, lang string) {
	s := bufio.NewScanner(bytes.NewBuffer(text))
	for s.Scan() {
		b := s.Bytes()
		if len(b) > 0 {
			g.out(out, indent)
			g.out(out, s.Bytes())
		}
		g.cr(out)
	}
	g.cr(out)
}

// DocumentHeader writes the beginning of the package documentation.
func (g *GodocRenderer) DocumentHeader(out *bytes.Buffer) {
	out.WriteString("// Package " + g.pkg + " is the ")
}

// DocumentFooter writes the end of the package documentation
func (g *GodocRenderer) DocumentFooter(out *bytes.Buffer) {
	out.WriteString("\npackage " + g.pkg + "\n")
}
