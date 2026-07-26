package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// scratchProject fabricates a minimal gsxui-initialized module and chdirs
// into it; gsx generate is stubbed out through the runCommand seam. The
// returned slice records every stubbed invocation.
func scratchProject(t *testing.T) (string, *[][]string) {
	t.Helper()
	dir := t.TempDir()
	must := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}
	must(os.WriteFile(filepath.Join(dir, "go.mod"), []byte("module example.com/app\n\ngo 1.26\n"), 0o644))
	must(os.WriteFile(filepath.Join(dir, "gsxui.json"), []byte(`{"ui":"ui","js":"web/gsxui","css":"web/gsxui.css"}`), 0o644))
	wd, err := os.Getwd()
	must(err)
	must(os.Chdir(dir))
	t.Cleanup(func() { os.Chdir(wd) })
	ran := &[][]string{}
	prev := runCommand
	runCommand = func(dir, name string, args ...string) error {
		*ran = append(*ran, append([]string{name}, args...))
		return nil
	}
	t.Cleanup(func() { runCommand = prev })
	return dir, ran
}

func TestAddVendorsComponentAndCSS(t *testing.T) {
	dir, ran := scratchProject(t)
	if err := Run([]string{"add", "button", "input"}); err != nil {
		t.Fatal(err)
	}
	if len(*ran) == 0 {
		t.Error("add never ran gsx generate")
	}
	for _, path := range []string{
		"ui/uswds/button.gsx",
		"ui/uswds/input.gsx",
		"ui/uswds/NOTICE.md",
		"web/uswds/uswds-tokens.css",
		"web/uswds/uswds-components.css",
	} {
		if _, err := os.Stat(filepath.Join(dir, path)); err != nil {
			t.Errorf("missing %s: %v", path, err)
		}
	}
	src, _ := os.ReadFile(filepath.Join(dir, "ui/uswds/button.gsx"))
	if !strings.Contains(string(src), "package uswds") {
		t.Error("vendored button.gsx should be package uswds")
	}
}

func TestAddRefusesToClobberLocalEdits(t *testing.T) {
	dir, _ := scratchProject(t)
	if err := Run([]string{"add", "button"}); err != nil {
		t.Fatal(err)
	}
	edited := filepath.Join(dir, "ui/uswds/button.gsx")
	if err := os.WriteFile(edited, []byte("package uswds // locally edited\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := Run([]string{"add", "button"}); err == nil {
		t.Fatal("re-add over a local edit should error without --overwrite")
	}
	if err := Run([]string{"add", "--overwrite", "button"}); err != nil {
		t.Fatalf("add --overwrite should succeed: %v", err)
	}
	src, _ := os.ReadFile(edited)
	if strings.Contains(string(src), "locally edited") {
		t.Error("--overwrite should have replaced the local edit")
	}
}

func TestAddUnknownComponent(t *testing.T) {
	_, _ = scratchProject(t)
	err := Run([]string{"add", "accordion"})
	if err == nil || !strings.Contains(err.Error(), "unknown component") {
		t.Fatalf("want unknown-component error, got %v", err)
	}
}

func TestAddOutsideGsxuiProject(t *testing.T) {
	dir := t.TempDir()
	wd, _ := os.Getwd()
	os.Chdir(dir)
	t.Cleanup(func() { os.Chdir(wd) })
	err := Run([]string{"add", "button"})
	if err == nil || !strings.Contains(err.Error(), "gsxui init") {
		t.Fatalf("want gsxui-project guidance, got %v", err)
	}
}
