# MingOak

A lightweight, easy-to-use, in-memory file tree.

```
$ go get github.com/JulzDiverse/mingoak
```

## Usage

```go

  root := mingoak.MkRoot()

  root.MkDirAll("path/to/dir/")
  root.WriteFile("path/to/dir/file", []byte("test"))

  fileInfo, _ := root.ReadDir("path/to/dir")
  for _, v := fileInfo {
     fmt.Println(v.IsDir) //true or false
     fmt.Println(v.Name)  //name of file/dir
  }

  file, _ = root.ReadFile("path/to/dir/file")
```


