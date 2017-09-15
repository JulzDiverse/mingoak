package mingoak

import (
	"errors"
	"fmt"
	"path/filepath"
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
		components: []FileInfo{},
		name:       "root",
		time:       time.Now(),
	}
}

func (d *Dir) WriteFile(path string, file []byte) error {
	if path == "" {
		return errors.New("No file name or path provided!")
	}

	current := d
	sl := slicePath(path)
	for i, name := range sl {
		if i == len(sl)-1 {
			file := File{
				content: file,
				name:    name,
				time:    time.Now(),
			}
			current.components = append(current.components, file)
			break
		}
		dir, err := getDirForName(current, name)
		if err != nil {
			return errors.New(fmt.Sprintf("error %s: %s", path, err))
		}
		current = dir
	}
	return nil
}

func (d *Dir) ReadFile(path string) ([]byte, error) {
	current := d
	var file []byte
	sl := slicePath(path)
	for i, name := range sl {
		if i == len(sl)-1 {
			fileInfo, err := getFileForName(current, name)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error %s: %s", path, err.Error()))
			}

			file = fileInfo.content
		}
		dir, err := getDirForName(current, name)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("error %s: %s", path, err))
		}
		current = dir
	}
	return file, nil
}

func (d *Dir) MkDirAll(path string) {
	current := d
	for _, name := range slicePath(path) {
		dir := Dir{
			components: []FileInfo{},
			name:       name,
			time:       time.Now(),
		}
		current.components = append(current.components, &dir)
		current, _ = getDirForName(current, name)
	}
}

func (d Dir) Walk(path string) ([]string, error) {
	fi, err := d.ReadDir(path)
	if err != nil {
		return nil, err
	}

	files := walkRecursion(fi, path)
	return files, nil
}

func walkRecursion(fileInfos []FileInfo, basepath string) []string {
	files := []string{}
	for _, v := range fileInfos {
		if v.IsDir() {
			dir := v.(*Dir)
			subFiles := walkRecursion(dir.components, filepath.Join(basepath, v.Name()))
			files = append(files, subFiles...)
		} else {
			files = append(files, filepath.Join(basepath, v.Name()))
		}
	}
	return files
}

func (d *Dir) ReadDir(path string) ([]FileInfo, error) {
	current := d
	sl := slicePath(path)
	var result []FileInfo
	for i, name := range sl {
		if i == len(sl)-1 {
			dir, err := getDirForName(current, name)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("error %s: %s", path, err))
			}
			result = dir.components
		}
		current, _ = getDirForName(current, name)
	}
	return result, nil
}

func slicePath(path string) []string {
	sl := strings.Split(path, "/")
	if sl[len(sl)-1] == "" {
		sl = sl[:len(sl)-1]
	}
	return sl
}

func getDirForName(dir *Dir, name string) (*Dir, error) {
	for _, v := range dir.components {
		if v.Name() == name && v.IsDir() {
			d := v.(*Dir)
			return d, nil
		}
	}
	return &Dir{}, errors.New("dir not found!")
}

func getFileForName(dir *Dir, name string) (File, error) {
	for _, v := range dir.components {
		if v.Name() == name && !v.IsDir() {
			return v.(File), nil
		}
	}
	return File{}, errors.New("file not found!")
}

func (d Dir) PrintFilePaths() {
	for _, v := range d.components {
		fmt.Println("NAME: ", v.Name())
	}
}
