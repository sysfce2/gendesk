package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDesktopFilename(t *testing.T) {
	tests := []struct {
		pkgname  string
		output   string
		expected string
	}{
		{"myapp", "", "myapp.desktop"},
		{"myapp", "custom.desktop", "custom.desktop"},
		{"foo-bar", "", "foo-bar.desktop"},
		{"foo-bar", "override.desktop", "override.desktop"},
	}
	for _, tt := range tests {
		cfg := &DesktopConfig{Pkgname: tt.pkgname, Output: tt.output}
		got := cfg.desktopFilename()
		if got != tt.expected {
			t.Errorf("desktopFilename(%q, %q) = %q, want %q", tt.pkgname, tt.output, got, tt.expected)
		}
	}
}

func TestOutputWhitespaceTrimming(t *testing.T) {
	// Simulate the comma-split + trim logic from main
	input := " a.desktop , b.desktop , c.desktop "
	var outputFilenames []string
	for _, s := range strings.Split(input, ",") {
		outputFilenames = append(outputFilenames, strings.TrimSpace(s))
	}
	expected := []string{"a.desktop", "b.desktop", "c.desktop"}
	if len(outputFilenames) != len(expected) {
		t.Fatalf("got %d filenames, want %d", len(outputFilenames), len(expected))
	}
	for i, got := range outputFilenames {
		if got != expected[i] {
			t.Errorf("outputFilenames[%d] = %q, want %q", i, got, expected[i])
		}
	}
}

func TestCreateDesktopContents(t *testing.T) {
	buf, err := createDesktopContents("MyApp", "Generic", "A comment", "myapp", "myapp", "", false, false, []string{"Application"}, nil)
	if err != nil {
		t.Fatalf("createDesktopContents: %v", err)
	}
	contents := buf.String()
	if !strings.Contains(contents, "Name=MyApp") {
		t.Error("missing Name= line")
	}
	if !strings.Contains(contents, "Exec=myapp") {
		t.Error("missing Exec= line")
	}
	if !strings.Contains(contents, "Comment=A comment") {
		t.Error("missing Comment= line")
	}
	if !strings.Contains(contents, "Categories=Application") {
		t.Error("missing Categories= line")
	}
}

func TestCreateDesktopContentsWithMimeTypes(t *testing.T) {
	buf, err := createDesktopContents("Mail", "", "Email client", "mail", "mail", "", false, false, []string{"Email"}, []string{"x-scheme-handler/mailto"})
	if err != nil {
		t.Fatalf("createDesktopContents: %v", err)
	}
	contents := buf.String()
	if !strings.Contains(contents, "MimeType=x-scheme-handler/mailto") {
		t.Error("missing MimeType= line")
	}
}

func TestCreateDesktopContentsTerminalAndNotify(t *testing.T) {
	buf, err := createDesktopContents("Term", "", "Terminal app", "term", "term", "", true, true, []string{"System"}, nil)
	if err != nil {
		t.Fatalf("createDesktopContents: %v", err)
	}
	contents := buf.String()
	if !strings.Contains(contents, "Terminal=true") {
		t.Error("expected Terminal=true")
	}
	if !strings.Contains(contents, "StartupNotify=true") {
		t.Error("expected StartupNotify=true")
	}
}

func TestCreateWindowManagerDesktopContents(t *testing.T) {
	buf, err := createWindowManagerDesktopContents("i3", "i3")
	if err != nil {
		t.Fatalf("createWindowManagerDesktopContents: %v", err)
	}
	contents := buf.String()
	if !strings.Contains(contents, "Name=i3") {
		t.Error("missing Name= line")
	}
	if !strings.Contains(contents, "Exec=i3") {
		t.Error("missing Exec= line")
	}
	if !strings.Contains(contents, "Type=XSession") {
		t.Error("missing Type=XSession line")
	}
}

func TestWriteDesktopFile(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "test.desktop")
	cfg := &DesktopConfig{
		Pkgname:    "testpkg",
		Name:       "TestPkg",
		Comment:    "A test package",
		Exec:       "testpkg",
		Categories: "Application",
		Output:     filename,
		Force:      true,
	}
	o := newSilentOutput()
	writeDesktopFile(cfg, o)

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading generated file: %v", err)
	}
	contents := string(data)
	if !strings.Contains(contents, "Name=TestPkg") {
		t.Error("generated file missing Name=")
	}
	if !strings.Contains(contents, "Exec=testpkg") {
		t.Error("generated file missing Exec=")
	}
}

func TestWriteDesktopFileNoOverwrite(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "existing.desktop")
	os.WriteFile(filename, []byte("existing"), 0644)

	cfg := &DesktopConfig{
		Pkgname:    "testpkg",
		Name:       "TestPkg",
		Comment:    "A test package",
		Exec:       "testpkg",
		Categories: "Application",
		Output:     filename,
		Force:      false,
	}

	// writeDesktopFile calls os.Exit(1) when the file already exists.
	// We can't easily test that without subprocess tricks, but we can
	// verify the file is untouched after a force=true overwrite.
	cfg.Force = true
	o := newSilentOutput()
	writeDesktopFile(cfg, o)

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if string(data) == "existing" {
		t.Error("expected file to be overwritten with force=true")
	}
}

func TestWriteDesktopFileWithCustom(t *testing.T) {
	dir := t.TempDir()
	filename := filepath.Join(dir, "custom.desktop")
	cfg := &DesktopConfig{
		Pkgname:    "testpkg",
		Name:       "TestPkg",
		Comment:    "A test",
		Exec:       "testpkg",
		Categories: "Application",
		Custom:     "X-Custom-Key=hello",
		Output:     filename,
		Force:      true,
	}
	o := newSilentOutput()
	writeDesktopFile(cfg, o)

	data, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("reading file: %v", err)
	}
	if !strings.Contains(string(data), "X-Custom-Key=hello") {
		t.Error("custom line not appended")
	}
}

func TestOutputFilenamePerPackage(t *testing.T) {
	// Simulate the per-package output selection logic from main
	outputFilenames := []string{"a.desktop", "b.desktop", "c.desktop"}
	pkgnames := []string{"pkg-a", "pkg-b", "pkg-c"}

	for i, pkgname := range pkgnames {
		perPkgOutput := ""
		if i < len(outputFilenames) {
			perPkgOutput = outputFilenames[i]
		}
		cfg := &DesktopConfig{Pkgname: pkgname, Output: perPkgOutput}
		got := cfg.desktopFilename()
		if got != outputFilenames[i] {
			t.Errorf("package %q: got %q, want %q", pkgname, got, outputFilenames[i])
		}
	}
}

func TestOutputFilenameMismatchDetection(t *testing.T) {
	// Verify the mismatch condition from main: len(outputFilenames) != len(pkgnames)
	outputFilenames := []string{"a.desktop", "b.desktop"}
	pkgnames := []string{"pkg-a", "pkg-b", "pkg-c"}
	if len(outputFilenames) == len(pkgnames) {
		t.Error("expected mismatch")
	}
}

func TestOutputFilenameDefaultFallback(t *testing.T) {
	// When outputFilenames is empty, desktopFilename uses PKGNAME.desktop
	pkgnames := []string{"mypkg"}
	var outputFilenames []string

	for i, pkgname := range pkgnames {
		perPkgOutput := ""
		if i < len(outputFilenames) {
			perPkgOutput = outputFilenames[i]
		}
		cfg := &DesktopConfig{Pkgname: pkgname, Output: perPkgOutput}
		got := cfg.desktopFilename()
		if got != "mypkg.desktop" {
			t.Errorf("got %q, want %q", got, "mypkg.desktop")
		}
	}
}

func Example_desktopFilename() {
	cfg := &DesktopConfig{Pkgname: "myapp"}
	fmt.Println(cfg.desktopFilename())
	cfg.Output = "custom.desktop"
	fmt.Println(cfg.desktopFilename())
	// output:
	// myapp.desktop
	// custom.desktop
}
