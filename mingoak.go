package mingoak

import (
	"errors"
	"path/filepath"
	"strings"
)

type Component interface {
	IsDir() bool
}

type Dir struct {
	components map[string]Component
}

type File struct {
	content []byte
}

type FileInfo struct {
	Name  string
	IsDir bool
}

func MkRoot() *Dir {
	return &Dir{
		components: map[string]Component{},
	}
}

func (f File) IsDir() bool {
	return false
}

func (d Dir) IsDir() bool {
	return true
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
	return nil, errors.New("File not found!")
}

func (d Dir) MkDirAll(path string) {
	current := d
	for _, name := range slicePath(path) {
		current.components[name] = Dir{
			map[string]Component{},
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

	for k, v := range dir.components {
		leafs = append(leafs, FileInfo{
			Name:  k,
			IsDir: v.IsDir(),
		})
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
	return files
}

func (d Dir) getDir(path string) (Dir, error) {
	current := d
	for _, name := range slicePath(path) {
		if result, ok := current.components[name]; ok && result.IsDir() {
			current = result.(Dir)
		} else {
			return Dir{}, errors.New("Directory not found!")
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
