package parser

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/fatih/structtag"
)

type structType struct {
	name string
	node ast.Node
}

// Options Parser Options
type Options struct {
	Input       string
	Structs     []string
	Lines       []int
	Tags        []string
	FieldFormat FieldFormat
	Rewrite     bool
	Out         io.Writer

	fset    *token.FileSet
	astFile ast.Node
}

func (o *Options) validate() error {
	err := o.FieldFormat.validate()
	if err != nil {
		return err
	}

	if o.Input == "" {
		return fmt.Errorf("input path required")
	}

	if o.Input, err = filepath.Abs(o.Input); err != nil {
		return err
	}

	if len(o.Tags) == 0 {
		o.Tags = []string{"json"}
	}

	return nil
}

// Parse parse input file in
func (o *Options) Parse() error {
	if err := o.validate(); err != nil {
		return err
	}

	if err := o.parse(); err != nil {
		return err
	}

	structs := o.selectStructs()
	if len(structs) == 0 {
		return nil
	}

	if err := o.processStructs(structs); err != nil {
		return err
	}

	if !o.Rewrite {
		if o.Out == nil {
			o.Out = os.Stdout
		}

		return o.output(o.Out)
	}

	var buf bytes.Buffer
	if err := o.output(&buf); err != nil {
		return err
	}

	return ioutil.WriteFile(o.Input, buf.Bytes(), 0)
}

func (o *Options) parse() error {
	o.fset = token.NewFileSet()

	var err error
	o.astFile, err = parser.ParseFile(o.fset, o.Input, nil, parser.ParseComments)
	return err
}

func (o *Options) selectStructs() []structType {
	snameMap := map[string]struct{}{}
	for _, sname := range o.Structs {
		snameMap[sname] = struct{}{}
	}

	structs := make([]structType, 0, 100)
	ast.Inspect(o.astFile, func(n ast.Node) bool {
		var t ast.Expr
		var name string

		switch n := n.(type) {
		case *ast.TypeSpec:
			if n.Type == nil {
				return true
			}

			t = n.Type
			name = n.Name.Name

		case *ast.CompositeLit:
			t = n.Type

		}

		_, ok := t.(*ast.StructType)
		if !ok {
			return true
		}

		if len(snameMap) > 0 {
			if _, ok := snameMap[name]; !ok {
				return true
			}
		}

		structs = append(structs, structType{
			name: name,
			node: n,
		})

		return true
	})

	return structs
}

func (o *Options) processStructs(structs []structType) error {
	lines := map[int]struct{}{}
	for _, line := range o.Lines {
		lines[line] = struct{}{}
	}

	var errs multierr
	ast.Inspect(o.astFile, func(n ast.Node) bool {
		s, ok := n.(*ast.StructType)
		if !ok {
			return true
		}

		for len(structs) > 0 {
			if s.Pos() >= structs[0].node.Pos() && s.End() <= structs[0].node.End() {
				break
			}

			if s.End() < structs[0].node.Pos() {
				return true
			}

			structs = structs[1:]
		}

		if len(structs) == 0 {
			return true
		}

		for _, field := range s.Fields.List {
			if len(field.Names) == 0 {
				continue
			}

			pos := o.fset.Position(field.Pos())

			if len(lines) > 0 {
				if _, ok := lines[pos.Line]; !ok {
					continue
				}
			}

			fieldname := field.Names[0].Name

			if field.Tag == nil {
				field.Tag = &ast.BasicLit{}
			}

			tag, err := o.processField(fieldname, field.Tag.Value)
			if err != nil {
				errs = append(errs, fmt.Errorf("%s:%d:%d:%s", pos.Filename, pos.Line, pos.Column, err))
				continue
			}

			field.Tag.Value = tag
		}

		return true
	})

	return errs.err()
}

func (o *Options) processField(name, tag string) (string, error) {
	var str string
	var err error
	if tag != "" {
		str, err = strconv.Unquote(tag)
		if err != nil {
			return "", err
		}
	}

	oldTags, err := structtag.Parse(str)
	if err != nil {
		return "", err
	}

	newTags := structtag.Tags{}

	for _, t := range o.Tags {
		pieces := strings.SplitN(t, ":", 2)
		tkey := pieces[0]

		tag, err := oldTags.Get(tkey)
		if err != nil {
			// tag not found
			var tname string

			if len(pieces) == 2 {
				tname = pieces[1]
			} else {
				tname = o.FieldFormat.transform(tkey, name)
			}

			tag = &structtag.Tag{
				Key:  tkey,
				Name: tname,
			}
		}

		newTags.Set(tag)
	}

	oldTags.Delete(o.Tags...)

	if remainKeys := oldTags.Keys(); len(remainKeys) > 0 {
		for _, rKey := range remainKeys {
			rTag, _ := oldTags.Get(rKey)
			newTags.Set(rTag)
		}
	}

	res := newTags.String()
	if res != "" {
		res = "`" + res + "`"
	}

	return res, nil
}

func (o *Options) output(w io.Writer) error {
	return format.Node(w, o.fset, o.astFile)
}
