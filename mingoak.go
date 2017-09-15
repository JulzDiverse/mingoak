package mingoak

import (
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type FileMode uint32

type FileInfo interface {
	Name() string
	Size() int64
	Mode() FileMode
	ModTime() time.Time
	IsDir() bool
	Sys() interface{}
}

func MkRoot() *Dir {
	return &Dir{
		components: map[string]FileInfo{},
		name:       "root",
		time:       time.Now(),
	}
}

func (d Dir) WriteFile(path string, file []byte) error {
	if path == "" {
		return errors.New("No file name or path provided!")
	}

	current := d
	sl := slicePath(path)
	for i, name := range sl {
		if i == len(sl)-1 {
			current.components[name] = File{
				content: file,
				name:    name,
				time:    time.Now(),
			}
			break
		}
		current = current.components[name].(Dir)
	}
	return nil
}

func (d Dir) ReadFile(path string) ([]byte, error) {
	current := d
	for _, name := range slicePath(path) {
		if result, ok := current.components[name]; ok && result.IsDir() {
			current = result.(Dir)
		} else if result, ok := current.components[name]; ok {
			file := result.(File)
			return file.content, nil
		}
	}
	return nil, errors.New(fmt.Sprintf("File %s not found!", path))
}

func (d Dir) MkDirAll(path string) {
	current := d
	for _, name := range slicePath(path) {
		current.components[name] = Dir{
			components: map[string]FileInfo{},
			name:       name,
			time:       time.Now(),
		}
		current = current.components[name].(Dir)
	}
}

func (d Dir) ReadDir(dirname string) ([]FileInfo, error) {
	leafs := []FileInfo{}
	dir, err := d.getDir(dirname)
	if err != nil {
		return nil, err
	}

	for _, v := range dir.components {
		leafs = append(leafs, v)
	}

	return leafs, nil
}

func (d Dir) Walk(path string) ([]string, error) {
	dir, err := d.getDir(path)
	if err != nil {
		return nil, err
	}

	files := walkRecursion(dir, path)
	return files, nil
}

func walkRecursion(dir Dir, basepath string) []string {
	files := []string{}
	for k, v := range dir.components {
		if v.IsDir() {
			subFiles := walkRecursion(v.(Dir), filepath.Join(basepath, k))
			files = append(files, subFiles...)
		} else {
			files = append(files, filepath.Join(basepath, k))
		}
	}
	sort.Strings(files)
	return files
}

func (d Dir) getDir(path string) (Dir, error) {
	current := d
	for _, name := range slicePath(path) {
		if result, ok := current.components[name]; ok && result.IsDir() {
			current = result.(Dir)
		} else {
			return Dir{}, errors.New(fmt.Sprintf("Directory %s not found!", path))
		}
	}
	return current, nil
}

func slicePath(path string) []string {
	sl := strings.Split(path, "/")
	if sl[len(sl)-1] == "" {
		sl = sl[:len(sl)-1]
	}
	return sl
}
