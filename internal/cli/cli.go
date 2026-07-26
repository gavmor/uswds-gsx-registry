// Package cli implements the uswds-gsx command: add, list. It composes
// with gsxui rather than replacing it — the consumer project is expected to
// be gsxui-initialized, and this CLI reads the same gsxui.json to learn
// where vendored Go packages and CSS live. The vendoring contract mirrors
// gsxui's: copy source into the project (you own the code), never touch a
// locally modified file without --overwrite, then run gsx generate.
//
// Components land in <ui>/uswds/ as package uswds — their own subpackage,
// so they never collide with vendored gsxui components in <ui>/. CSS lands
// in <cssdir>/uswds/; the two @import lines are the consumer's one manual
// step (printed after add; we don't edit your stylesheet).
package cli

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"

	uswdsgsx "github.com/gavmor/uswds-gsx-registry"
)

// Config is the consumer's gsxui.json: where vendored Go packages, JS, and
// CSS live, relative to the module root. Owned by gsxui; read, never
// written, here.
type Config struct {
	UI  string `json:"ui"`
	JS  string `json:"js"`
	CSS string `json:"css"`
}

func loadConfig(dir string) (Config, error) {
	data, err := os.ReadFile(filepath.Join(dir, "gsxui.json"))
	if err != nil {
		if os.IsNotExist(err) {
			return Config{}, fmt.Errorf("gsxui.json not found — uswds-gsx vendors into a gsxui project; run 'gsxui init' first")
		}
		return Config{}, fmt.Errorf("reading gsxui.json: %w", err)
	}
	var c Config
	if err := json.Unmarshal(data, &c); err != nil {
		return Config{}, fmt.Errorf("parsing gsxui.json: %w", err)
	}
	return c, nil
}

// runCommand is the seam for external processes (go, gsx). Unit tests stub
// it; the real implementation streams output through.
var runCommand = func(dir, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Run dispatches the uswds-gsx subcommands.
func Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: uswds-gsx <add|list> [args]")
	}
	switch args[0] {
	case "add":
		return runAdd(args[1:])
	case "list":
		return runList(args[1:])
	default:
		return fmt.Errorf("unknown command %q (want add or list)", args[0])
	}
}

// components derives the component list from the embedded ui/ tree: a
// component is a .gsx file basename, exactly gsxui's registry model.
func components() ([]string, error) {
	entries, err := fs.ReadDir(uswdsgsx.Files, "ui")
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if name, ok := strings.CutSuffix(e.Name(), ".gsx"); ok && !e.IsDir() {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	return names, nil
}

func runList(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("usage: uswds-gsx list")
	}
	names, err := components()
	if err != nil {
		return err
	}
	for _, n := range names {
		fmt.Println(n)
	}
	return nil
}

// runAdd vendors the requested components into <cfg.UI>/uswds/ and the
// component CSS into <dir(cfg.CSS)>/uswds/. The CSS ships on every add —
// it is the components' single stylesheet, not per-component — and the
// content-compare in writeVendored makes repeats free.
func runAdd(args []string) error {
	fs2 := flag.NewFlagSet("add", flag.ContinueOnError)
	overwrite := fs2.Bool("overwrite", false, "replace locally modified files")
	if err := fs2.Parse(args); err != nil {
		return err
	}
	names := fs2.Args()
	if len(names) == 0 {
		return fmt.Errorf("usage: uswds-gsx add [--overwrite] <component>...")
	}
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	cfg, err := loadConfig(dir)
	if err != nil {
		return err
	}
	known, err := components()
	if err != nil {
		return err
	}
	for _, name := range names {
		found := false
		for _, k := range known {
			if k == name {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("unknown component %q (run 'uswds-gsx list')", name)
		}
	}
	fmt.Printf("adding: %s\n", strings.Join(names, " "))

	for _, name := range names {
		src, err := fs.ReadFile(uswdsgsx.Files, "ui/"+name+".gsx")
		if err != nil {
			return err
		}
		if err := writeVendored(filepath.Join(dir, cfg.UI, "uswds", name+".gsx"), src, *overwrite); err != nil {
			return err
		}
	}

	cssDir := filepath.Join(dir, filepath.Dir(cfg.CSS), "uswds")
	for _, css := range []string{"uswds-tokens.css", "uswds-components.css"} {
		src, err := fs.ReadFile(uswdsgsx.Files, "css/"+css)
		if err != nil {
			return err
		}
		if err := writeVendored(filepath.Join(cssDir, css), src, *overwrite); err != nil {
			return err
		}
	}

	notice, err := fs.ReadFile(uswdsgsx.Files, "NOTICE.md")
	if err != nil {
		return err
	}
	if err := writeVendored(filepath.Join(dir, cfg.UI, "uswds", "NOTICE.md"), notice, true); err != nil {
		return err
	}

	if err := runCommand(dir, "go", "tool", "gsx", "generate"); err != nil {
		return fmt.Errorf("gsx generate: %w — if the gsx tool is missing, run 'gsxui init' (or 'go get -tool github.com/gsxhq/gsx/cmd/gsx@latest')", err)
	}

	fmt.Printf(`done — if not already present, add to the top of %s (with your other @import lines):

  @import "./uswds/uswds-tokens.css";
  @import "./uswds/uswds-components.css";

(projects with their own --usa-* token layer should import exactly one
source of the primitives — see uswds-tokens.css), then rebuild your CSS
and: go build ./...
`, cfg.CSS)
	return nil
}

// writeVendored writes content to path unless an identical file is already
// there; a differing file is an error without --overwrite (the file is
// yours once vendored — same contract as gsxui).
func writeVendored(path string, content []byte, overwrite bool) error {
	if existing, err := os.ReadFile(path); err == nil {
		if string(existing) == string(content) {
			return nil
		}
		if !overwrite {
			return fmt.Errorf("%s differs from the uswds-gsx version — pass --overwrite to replace it", path)
		}
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, content, 0o644)
}
