package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"time"

	"github.com/yuin/goldmark"
	highlighting "github.com/yuin/goldmark-highlighting"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/renderer/html"
)

const (
	defaultTemplate = `<!DOCTYPE html>
<html>
  <head>
    <meta http-equiv="content-type" content="text/html; charset=utf-8">
      <title>{{ .Title }}</title>
  </head>
  <body>
  {{ .Body }}
  </body>
</html>
`
)

var browsers = map[string]string{
	"chrome":  "google-chrome",
	"firefox": "firefox",
}

var Version = "0.3.0"

// content type represents the HTML content to add into the template
type content struct {
	Title string
	Body  template.HTML
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "mdp - preview markdown in your browser.\n\n")
		fmt.Fprintln(flag.CommandLine.Output(), "Usage information:")
		flag.PrintDefaults()
	}
	filename := flag.String("file", "", "Markdown file to preview.")
	skipPreview := flag.Bool("skip", false, "Skip auto-preview and prevent auto-delete of html file.")
	browser := flag.String("browser", "chrome", "Which browser to use. Default is Chrome. Options are chrome and firefox.")
	version := flag.Bool("version", false, "Print version and exit.")
	flag.Parse()

	if *version {
		fmt.Printf("mdp: v%s\n", Version)
		os.Exit(0)
	}
	if *filename == "" {
		flag.Usage()
		os.Exit(0)
	}
	if _, ok := browsers[*browser]; !ok {
		fmt.Println(*browser)
		fmt.Println("browser must be 'chrome' or 'firefox'.")
		os.Exit(1)
	}
	if err := run(*filename, os.Stdout, *skipPreview, *browser); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(filename string, out io.Writer, skipPreview bool, browser string) error {
	// Read all the data from the input file and check for errors
	input, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	htmlData, err := parseContent(input, filename)
	if err != nil {
		return err
	}

	// Create temporary file
	temp, err := ioutil.TempFile("", "mdp.*.html")
	if err != nil {
		return err
	}
	if err := temp.Close(); err != nil {
		return err
	}
	outName := temp.Name()
	fmt.Fprintln(out, outName)

	if err := saveHTML(outName, htmlData); err != nil {
		return err
	}

	if skipPreview {
		return nil
	}

	defer os.Remove(outName)
	return preview(outName, browser)
}

func parseContent(source []byte, filename string) ([]byte, error) {
	// Convert markdown to HTML
	var con bytes.Buffer
	md := goldmark.New(
		goldmark.WithExtensions(
			highlighting.Highlighting,
			extension.Table,
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
			html.WithXHTML(),
		))
	if err := md.Convert(source, &con); err != nil {
		return nil, err
	}

	// Parse the contents of the defaultTemplate const into a new Template
	t, err := template.New("mdp").Parse(defaultTemplate)
	if err != nil {
		return nil, err
	}

	// Add markdown to template
	c := content{
		Title: filename,
		Body:  template.HTML(con.String()),
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, c); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func preview(fname string, browser string) error {
	// Locate Chrome in the PATH
	if browser == "chrome" {
		browser = "google-chrome"
	}
	browserPath, err := exec.LookPath(browser)
	if err != nil {
		return err
	}

	// Open the file in the browser
	if err := exec.Command(browserPath, fname).Start(); err != nil {
		return err
	}

	// Give the browser some time to open the file before deleting it
	time.Sleep(2 * time.Second)
	return nil
}

func saveHTML(outFname string, data []byte) error {
	return ioutil.WriteFile(outFname, data, 0644)
}
