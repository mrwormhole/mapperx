# MapperX

## What we aim to do
* 1-to-1 struct-struct mapping
* Maps underlying nested structs 
* Maps same types with same variable names by default
* Maps same types with different names by tag specification
* Does use reflection to generate mapperx package

## What we don't aim to do
* Doesn't aim to do aggregation
* Doesn't map not equal types
* Doesn't use reflection at runtime

## Getting Started
```
    Mapperx heavily relies on code generation. This means that you need to specify 2 arguments source(file path and struct type) and target(file path and struct type)
    Then mapperx will generate a directory and a package called mapperx, when you run it. Easiest way to use is to use go generate compiler directive in your definitions at the start of a file.
```

```go
package main

//go:generate go run github.com/mrwormhole/mapperx/main.go github.com/yourusername/yourproject/domain.User github.com/yourusername/yourproject/domain.Admin
type User struct {
    Name string
    ID string
    Country string `mapperx:"Score"`
    Highscore string
}

type Admin struct {
    Name string
    ID string
    Country string
    Score string
    Permissions []string
}
```