package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/jamesainslie/gomd2svg"
	"github.com/jamesainslie/gomd2svg/theme"
)

var version = "dev"

func main() {
	if err := run(os.Args[1:], os.Stdin, os.Stdout, os.Stderr); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	if len(args) == 0 {
		return usage(stderr)
	}

	switch args[0] {
	case "render":
		return runRender(args[1:], stdin, stdout, stderr)
	case "themes":
		return runThemes(stdout)
	case "version":
		fmt.Fprintf(stdout, "gomd2svg %s\n", version)
		return nil
	case "help", "-h", "--help":
		return usage(stderr)
	default:
		return fmt.Errorf("unknown command: %s", args[0])
	}
}

func runRender(args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	fs := flag.NewFlagSet("render", flag.ContinueOnError)
	fs.SetOutput(stderr)
	output := fs.String("o", "", "output file (default: stdout)")
	themeName := fs.String("theme", "", "theme name (modern|default|dark|forest|neutral)")
	timing := fs.Bool("timing", false, "print timing info to stderr")
	if err := fs.Parse(args); err != nil {
		return err
	}

	var input []byte
	var err error
	if fs.NArg() > 0 {
		input, err = os.ReadFile(fs.Arg(0))
		if err != nil {
			return err
		}
	} else {
		input, err = io.ReadAll(stdin)
		if err != nil {
			return err
		}
	}

	if len(strings.TrimSpace(string(input))) == 0 {
		return errors.New("empty input")
	}

	opts := gomd2svg.Options{}
	if *themeName != "" {
		opts.ThemeName = *themeName
	}

	if *timing {
		result, err := gomd2svg.RenderWithTiming(string(input), opts)
		if err != nil {
			return err
		}
		fmt.Fprintf(stderr, "parse: %dus  layout: %dus  render: %dus  total: %.1fms\n",
			result.ParseUs, result.LayoutUs, result.RenderUs, result.TotalMs())
		return writeOutput(*output, result.SVG, stdout)
	}

	svg, err := gomd2svg.RenderWithOptions(string(input), opts)
	if err != nil {
		return err
	}
	return writeOutput(*output, svg, stdout)
}

func writeOutput(path, svg string, stdout io.Writer) error {
	if path == "" {
		_, err := io.WriteString(stdout, svg)
		return err
	}
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
	}
	return os.WriteFile(path, []byte(svg), 0o600)
}

func runThemes(w io.Writer) error {
	for _, name := range theme.Names() {
		fmt.Fprintln(w, name)
	}
	return nil
}

func usage(w io.Writer) error {
	fmt.Fprintln(w, `Usage: gomd2svg <command> [options]

Commands:
  render [file]   Render a .mmd file to SVG
  themes          List available themes
  version         Print version

Render options:
  -o <file>       Output file (default: stdout)
  -theme <name>   Theme: modern, default, dark, forest, neutral
  -timing         Print timing info to stderr

Examples:
  gomd2svg render diagram.mmd -o diagram.svg
  gomd2svg render -theme dark diagram.mmd > out.svg
  cat diagram.mmd | gomd2svg render > out.svg
  gomd2svg render -theme forest -timing diagram.mmd -o out.svg`)
	return nil
}
