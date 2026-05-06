// Package toolversions_test tests the `toolversions` package.
package toolversions_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Snehil-Shah/zxcv/internal/testutil"
	"github.com/Snehil-Shah/zxcv/internal/toolversions"
)

func TestWriteEntry_NewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tool-versions")
	if err := toolversions.WriteEntry(path, "nodejs", "20.11.0"); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "nodejs 20.11.0\n"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteEntry_UpdatesExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tool-versions")
	testutil.WriteFile(t, path, "nodejs 18.19.0\npython 3.11.0\n")
	if err := toolversions.WriteEntry(path, "nodejs", "20.11.0"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	want := "nodejs 20.11.0\npython 3.11.0\n"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteEntry_AppendsWhenAbsent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tool-versions")
	testutil.WriteFile(t, path, "nodejs 20.11.0\n")
	if err := toolversions.WriteEntry(path, "python", "3.11.0"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	want := "nodejs 20.11.0\npython 3.11.0\n"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteEntry_PreservesCommentsAndBlankLines(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".tool-versions")
	original := "# header\n\nnodejs 18.19.0 # company standard\n\npython 3.11.0\n"
	testutil.WriteFile(t, path, original)
	if err := toolversions.WriteEntry(path, "nodejs", "20.11.0"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	want := "# header\n\nnodejs 20.11.0 # company standard\n\npython 3.11.0\n"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestWriteEntry_DoesNotMatchPrefix(t *testing.T) {
	// "nodejs-old" should not be confused with "nodejs".
	dir := t.TempDir()
	path := filepath.Join(dir, ".tool-versions")
	testutil.WriteFile(t, path, "nodejs-old 14.0.0\n")
	if err := toolversions.WriteEntry(path, "nodejs", "20.11.0"); err != nil {
		t.Fatal(err)
	}
	got, _ := os.ReadFile(path)
	want := "nodejs-old 14.0.0\nnodejs 20.11.0\n"
	if string(got) != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
