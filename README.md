# MingOak

A lightweight in-memory file tree.

![mingoak](https://i.pinimg.com/736x/7f/4a/e3/7f4ae3efff8e660f80dfdfc8eb368d11--red-oak-tree-white-oak-tree.jpg =200x)

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


