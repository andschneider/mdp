# mdp

A simple CLI that will preview a markdown file in a browser.

## Installation

There are binaries for Linux, Mac, and Windows in the [releases page](https://github.com/andschneider/mdp/releases).

Or use `go get`:

```bash
go get github.com/andschneider/mdp
```

## Usage

- Pass in the name of the file you'd like to preview with the `-file` flag. 

```bash
mdp -file README.md
```

- To skip the autopreview use the `-skip` flag. This also doesn't remove the generated .html file if you'd like to save or view it later.

- Use `-help` flag to display the full usage information.

```bash
Usage of mdp:
  -file string
    	Markdown file to preview
  -skip
    	Skip auto-preview and prevent auto-delete of html file.
```

### inspiration

This repo is based on chapter 3 from the book [Powerful Command-Line Applications in Go](https://pragprog.com/titles/rggo/)
