// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package testhelpers

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func RunGoldenTests(t *testing.T, basedir string, fn func(t *testing.T, dir string)) {
	files, err := ioutil.ReadDir(basedir)
	if err != nil {
		t.Fatalf("ReadDir(%q) failed: %v", basedir, err)
	}
	for _, file := range files {
		name := file.Name()
		absPath := filepath.Join(basedir, name)
		t.Run(name, func(t *testing.T) {
			fn(t, absPath)
		})
	}
}

func MustReadFile(t *testing.T, p string) []byte {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		t.Fatalf("failed to read file %q: %v", p, err)
	}
	return b
}

func CompareGoldenFile(t *testing.T, p string, got []byte) {
	if os.Getenv("WRITE_GOLDEN_OUTPUT") != "" {
		// Short-circuit when the output is correct
		b, err := ioutil.ReadFile(p)
		if err == nil && bytes.Equal(b, got) {
			return
		}

		if err := ioutil.WriteFile(p, got, 0644); err != nil {
			t.Fatalf("failed to write golden output %s: %v", p, err)
		}
		t.Errorf("wrote output to %s", p)
	} else {
		want := MustReadFile(t, p)
		if diff := cmp.Diff(string(want), string(got)); diff != "" {
			t.Errorf("unexpected diff in %s: %s", p, diff)
		}
	}
}

func CopyDir(src, dest string) error {
	srcFiles, err := ioutil.ReadDir(src)
	if err != nil {
		return fmt.Errorf("ReadDir(%q) failed: %w", src, err)
	}
	for _, srcFile := range srcFiles {
		srcPath := filepath.Join(src, srcFile.Name())
		destPath := filepath.Join(dest, srcFile.Name())
		if srcFile.IsDir() {
			if err := CopyDir(srcPath, destPath); err != nil {
				return err
			}
		} else {
			if err := CopyFile(srcPath, destPath); err != nil {
				return err
			}
		}
	}
	return nil
}

func CopyFile(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("OpenFile(%q) failed: %w", src, err)
	}
	defer in.Close()

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("Create(%q) failed: %w", dest, err)
	}

	if _, err := io.Copy(out, in); err != nil {
		out.Close()
		return fmt.Errorf("byte copy from %s to %s failed: %w", src, dest, err)
	}

	if err := out.Close(); err != nil {
		return fmt.Errorf("Close(%q) failed: %w", dest, err)
	}

	return nil
}

func CompareDir(t *testing.T, expectDir, actualDir string) {
	expectFiles, err := ioutil.ReadDir(expectDir)
	if err != nil {
		t.Fatalf("failed to read expectation directory %q: %v", expectDir, err)
	}
	expectFileMap := make(map[string]fs.FileInfo)
	for _, expectFile := range expectFiles {
		expectFileMap[expectFile.Name()] = expectFile
	}

	actualFiles, err := ioutil.ReadDir(actualDir)
	if err != nil {
		t.Fatalf("failed to read actual directory %q: %v", actualDir, err)
	}
	actualFileMap := make(map[string]fs.FileInfo)
	for _, actualFile := range actualFiles {
		actualFileMap[actualFile.Name()] = actualFile
	}

	for _, expectFile := range expectFiles {
		name := expectFile.Name()
		actualFile := actualFileMap[name]
		if actualFile == nil {
			t.Errorf("expected file %s not found", name)
			continue
		}

		if expectFile.IsDir() {
			if !actualFile.IsDir() {
				t.Errorf("expected file %s was not a directory", name)
				continue
			}
			CompareDir(t, filepath.Join(expectDir, name), filepath.Join(actualDir, name))
		} else {
			if actualFile.IsDir() {
				t.Errorf("expected file %s was not a file", name)
				continue
			}
			CompareFile(t, expectDir, actualDir, name)
		}
	}

	for _, actualFile := range actualFiles {
		name := actualFile.Name()
		expectFile := expectFileMap[name]
		if expectFile == nil {
			t.Errorf("additional file %s found in output", name)
			continue
		}
	}
}

func CompareFile(t *testing.T, expectDir, actualDir string, relPath string) {
	expectAbs := filepath.Join(expectDir, relPath)

	actualAbs := filepath.Join(actualDir, relPath)
	actualBytes, err := ioutil.ReadFile(actualAbs)
	if err != nil {
		if os.IsNotExist(err) {
			t.Errorf("expected file %s not found", relPath)
		} else {
			t.Fatalf("error reading file %s: %v", actualAbs, err)
		}
	}

	CompareGoldenFile(t, expectAbs, actualBytes)
}
