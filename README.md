
# StratConClone-Go

StratConClone-Go is a clone of a turn-based strategy game originally released for the Apple II computer in 1982.

https://en.wikipedia.org/wiki/Strategic_Conquest


## status

a work-in-progress


## notes

### build
```
go build
```

### test
```
go test ./...
```

### test coverage
```
go get golang.org/x/tools/cmd/cover
go test -coverprofile cover.out
go tool cover -html=cover.out
```

### format code
```
gofmt -w .
```

### fix imports and format code
```
go install golang.org/x/tools/cmd/goimports@latest
goimports -w .
```

### calculate cyclomatic complexities of functions in Go source code
```
go get github.com/fzipp/gocyclo
gocyclo .
```
