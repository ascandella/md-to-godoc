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

// Package main is the Markdown to Godoc converter.
//
// Sort of like godocdown (https://github.com/robertkrimen/godocdown), but in
// reverse.
//
//
// md-to-godoc takes markdown as input, and generates godoc-formatted package
// documentation.
//
//
// Status
//
// Way, **way** alpha. Barebones. The minimalest.
//
// Code example
//
// Mostly here so we can see some code in godoc:
//
// Sample list
//
// • This is a test
//
// • And another test
//
//   func main() {
//     fmt.Println("Hello, world")
//   }
//
// Usage
//
// First, install the binary:
//
//   go get -u github.com/sectioneight/md-to-godoc
//
// Then, run it on one or more packages. If you'd like to generate a doc.go file
// in the current package (that already has a
// README.md), simply run
// md-to-godoc with no flags:
//
//   md-to-godoc
//
// Advanced usage
//
// To generate doc.go for all subpackages, you can do something like the
// following:
//
//
//   find . -name README.md \
//          -not -path "./vendor/*" | \
//          xargs -I% md-to-godoc -input=%
//
// Projects using md-to-godoc
//
// • UberFx, on GitHub (https://github.com/uber-go/fx) and
// godoc.org (https://godoc.org/go.uber.org/fx)
//
// • Jaeger, on Github (https://github.com/uber/jaeger) and
// godoc.org (https://godoc.org/github.com/uber/jaeger/services/agent)
//
// Licence
//
// Apache 2.0 (https://www.apache.org/licenses/LICENSE-2.0)
//
//
package main
