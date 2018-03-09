### jmgen

install
```
go get -v -u github.com/dtynn/jmgen
```

usage
```
add tags for go structs

Usage:
  jmgen [flags]

Flags:
  -f, --format string     field name format type, default empty. "camel" or "snake" allowed
  -h, --help              help for jmgen
  -i, --input string      input file path
  -l, --lines ints        specified lines
  -r, --rewrite           rewrite src file, default false
  -s, --structs strings   specified struct names
  -t, --tags strings      tags to add
```

example
```
jmgen -i ./example/example.go -t json,db,validate:required
```
