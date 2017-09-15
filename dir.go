package mingoak

import "time"

type Dir struct {
	components  map[string]FileInfo
	componentsl []FileInfo
	name        string
	time        time.Time
}

func (d Dir) IsDir() bool {
	return true
}

func (d Dir) Name() string {
	return d.name
}

func (d Dir) Size() int64 {
	return 1
}

func (d Dir) Mode() FileMode {
	return 777
}

func (d Dir) ModTime() time.Time {
	return d.time
}

func (d Dir) Sys() interface{} {
	return nil
}
