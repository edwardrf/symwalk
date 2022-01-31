package symwalk

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestWalk(t *testing.T) {
	dir, err := ioutil.TempDir("", "walktest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("failed to eval temp dir: %v", err)
	}

	mkdir(t, filepath.Join(dir, "a"))
	mkdir(t, filepath.Join(dir, "b", "c"))
	symlink(t, filepath.Join(dir, "a"), filepath.Join(dir, "b", "c", "Y"))
	symlink(t, filepath.Join(dir, "b", "c"), filepath.Join(dir, "a", "X"))
	touch(t, filepath.Join(dir, "a", "d"))

	var res strings.Builder
	Walk(dir, func(path string, info os.FileInfo, err error) error {
		fmt.Fprintln(&res, path[len(dir):])
		return nil
	})

	const expected = `
/a
/a/X
/b/c
/b/c/Y
/a/d
/b
`

	if res.String() != expected {
		t.Errorf("walk did not match expected dirs:\nExpected:\n%s---\nBut found:\n%s", expected, res.String())
	}
}

func TestLoop(t *testing.T) {
	dir, err := ioutil.TempDir("", "walklooptest")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir)
	dir, err = filepath.EvalSymlinks(dir)
	if err != nil {
		t.Fatalf("failed to eval temp dir: %v", err)
	}
	symlink(t, filepath.Join(dir, "a"), filepath.Join(dir, "b"))
	symlink(t, filepath.Join(dir, "b"), filepath.Join(dir, "a"))

	var res strings.Builder
	Walk(dir, func(path string, info os.FileInfo, err error) error {
		fmt.Fprintln(&res, path[len(dir):])
		return nil
	})
	const expected = `
/a
`
	if res.String() != expected {
		t.Errorf("walk did not match expected dirs:\nExpected:\n%s---\nBut found:\n%s", expected, res.String())
	}
}

func mkdir(t *testing.T, name string) {
	err := os.MkdirAll(name, 0755)
	if err != nil {
		t.Fatalf("failed to create dir %v: %v", name, err)
	}
}

func symlink(t *testing.T, from, to string) {
	err := os.Symlink(from, to)
	if err != nil {
		t.Fatalf("failed to create symbolic link %v->%v: %v", from, to, err)
	}
}

func touch(t *testing.T, name string) {
	f, err := os.Create(name)
	if err != nil {
		t.Fatalf("failed to create file %v: %v", name, err)
	}
	f.Close()
}
