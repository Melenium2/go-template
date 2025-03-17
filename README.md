# Go template

A "business" go template that can be used to quickly initializing
a new service. The example contains guidelines on how and where to write the
code.

## Initializing

Here we have used `gonew` tool to initialize the template.

### Installing gonew

More about `gonew` [here](https://go.dev/blog/gonew).

```bash
go install golang.org/x/tools/cmd/gonew@latest
```

### Initializing the project

```bash
gonew github.com/Melenium2/go-template <new go.mod name>
```

For example

```bash
gonew github.com/Melenium2/go-template \
  github.com/Melenium2/myservice

```

As a result the `myservice` folder will be created in the current directory. 
The project will be initialized with specified module name.

```bash
# bash
> ls

- ...
- myservice
- ...
```

```go
// go.mod
module github.com/Melenium2/myservice

go 1.24
```
