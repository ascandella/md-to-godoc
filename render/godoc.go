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

	"github.com/russross/blackfriday"
)

// GodocExtensions are the default markdown extensions for blackfriday
const GodocExtensions = 0 |
	blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
	blackfriday.EXTENSION_TABLES |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_SPACE_HEADERS |
	blackfriday.EXTENSION_HEADER_IDS |
	blackfriday.EXTENSION_BACKSLASH_LINE_BREAK

var (
	nl       = []byte("\n")
	nlnl     = []byte("\n\n")
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
	pkg string
}

// block-level callbacks
func (g *GodocRenderer) BlockCode(out *bytes.Buffer, text []byte, lang string) {
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

func (g *GodocRenderer) BlockQuote(out *bytes.Buffer, text []byte) {
	// TODO
}

func (g *GodocRenderer) BlockHtml(out *bytes.Buffer, text []byte) {
	// TODO
}

func (g *GodocRenderer) Header(out *bytes.Buffer, text func() bool, level int, id string) {
	marker := out.Len()

	if !text() {
		out.Truncate(marker)
		return
	}

	out.Write(nlnl)
}

func (g *GodocRenderer) HRule(out *bytes.Buffer) {
	// TODO
}

func (g *GodocRenderer) List(out *bytes.Buffer, text func() bool, flags int) {
	marker := out.Len()

	if !text() {
		out.Truncate(marker)
		return
	}

	out.Write(nlnl)
}

func (g *GodocRenderer) ListItem(out *bytes.Buffer, text []byte, flags int) {
	out.WriteString("• ")
	out.Write(text)
	out.Write(nlnl)
}

func (g *GodocRenderer) Paragraph(out *bytes.Buffer, text func() bool) {
	marker := out.Len()
	if !text() {
		out.Truncate(marker)
		return
	}
	out.Write(nlnl)
}

func (g *GodocRenderer) Table(out *bytes.Buffer, header []byte, body []byte, columnData []int) {
	// TODO
}

func (g *GodocRenderer) TableRow(out *bytes.Buffer, text []byte) {
	// TODO
}

func (g *GodocRenderer) TableHeaderCell(out *bytes.Buffer, text []byte, flags int) {
	// TODO
}

func (g *GodocRenderer) TableCell(out *bytes.Buffer, text []byte, flags int) {
	// TODO
}

func (g *GodocRenderer) Footnotes(out *bytes.Buffer, text func() bool) {
	// TODO
}

func (g *GodocRenderer) FootnoteItem(out *bytes.Buffer, name, text []byte, flags int) {
	// TODO
}

func (g *GodocRenderer) TitleBlock(out *bytes.Buffer, text []byte) {
	// TODO
}

// Span-level callbacks
func (g *GodocRenderer) AutoLink(out *bytes.Buffer, link []byte, kind int) {
	out.Write(link)
}

func (g *GodocRenderer) CodeSpan(out *bytes.Buffer, text []byte) {
	out.Write(text)
}

func (g *GodocRenderer) DoubleEmphasis(out *bytes.Buffer, text []byte) {
	out.Write(starstar)
	out.Write(text)
	out.Write(starstar)
}

func (g *GodocRenderer) Emphasis(out *bytes.Buffer, text []byte) {
	out.Write(star)
	out.Write(text)
	out.Write(star)
}

func (g *GodocRenderer) Image(out *bytes.Buffer, link []byte, title []byte, alt []byte) {
	// TODO
}

// LineBreak outputs a newline
func (g *GodocRenderer) LineBreak(out *bytes.Buffer) {
	out.Write(nl)
}

func (g *GodocRenderer) Link(out *bytes.Buffer, link []byte, title []byte, content []byte) {
	out.Write(content)
	out.WriteString(" (")
	out.Write(link)
	out.WriteString(")")
	if len(title) > 0 {
		out.WriteString(": ")
		out.Write(content)
	}
}

func (g *GodocRenderer) RawHtmlTag(out *bytes.Buffer, tag []byte) {
	// TODO
}

func (g *GodocRenderer) TripleEmphasis(out *bytes.Buffer, text []byte) {
	// TODO
}

func (g *GodocRenderer) StrikeThrough(out *bytes.Buffer, text []byte) {
	// TODO
}

func (g *GodocRenderer) FootnoteRef(out *bytes.Buffer, ref []byte, id int) {
	// TODO
}

// Low-level callbacks
func (g *GodocRenderer) Entity(out *bytes.Buffer, entity []byte) {
	out.Write(entity)
}

func (g *GodocRenderer) NormalText(out *bytes.Buffer, text []byte) {
	out.Write(text)
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

// GetFlags seems unused in blackfriday, but is implemented to satisfy the
// interface.
func (g *GodocRenderer) GetFlags() int {
	return 0
}
