# MingOak

A lightweight in-memory filesystem.

```
$ go get github.com/JulzDiverse/mingoak
```

## Usage

```go

  root := mingoak.MkRoot()

  root.MkDirAll("path/to/dir/")
  root.WriteFile("path/to/dir/file", []byte("test"))

  fileInfo := root.ReadDir("path/to/dir")
  for _, v := fileInfo {
     fmt.Println(v.IsDir) //true or false
     fmt.Println(v.Name)  //name of file/dir
  }

  file = root.ReadFile("path/to/dir/file")
```


