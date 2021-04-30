// Package symwalk provide a walk function similar to filepath.Walk but follow symbolic links
// but avoid visiting directories more than once even if there is a symbolic link loop
package symwalk

import (
	"os"
	"path/filepath"
)

// Walk is similar to filepath.Walk (https://golang.org/pkg/path/filepath/#Walk) except it follows
// the symbolic links it finds. The walk function keeps a list of all visited directories to avoid
// endless loop resulted from symbolic loops
func Walk(root string, fn filepath.WalkFunc) error {
	rr, err := filepath.EvalSymlinks(root) // Find real base if there is any symlinks in the path
	if err != nil {
		return err
	}

	visitedDirs := make(map[string]struct{})
	return filepath.Walk(rr, getWalkFn(visitedDirs, fn))
}

func getWalkFn(visitedDirs map[string]struct{}, fn filepath.WalkFunc) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fn(path, info, err)
		}

		if info.IsDir() {
			if _, ok := visitedDirs[path]; ok {
				return filepath.SkipDir
			}
			visitedDirs[path] = struct{}{}
		}

		if err := fn(path, info, err); err != nil {
			return err
		}

		if info.Mode()&os.ModeSymlink == 0 {
			return nil
		}

		// path is a symlink
		rp, err := filepath.EvalSymlinks(path)
		if err != nil {
			return err
		}

		ri, err := os.Stat(rp)
		if err != nil {
			return err
		}

		if ri.IsDir() {
			return filepath.Walk(rp, getWalkFn(visitedDirs, fn))
		}

		return nil
	}
}
