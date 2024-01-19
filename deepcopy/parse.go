package deepcopy

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

const (
	Tag             = "deepcopy"
	Cloneable       = "@copyable"     // 标记要不要实现Copy
	PtrRecv         = "@ptrrecv"      // 标记是否是指针接受
	ExportedOnly    = "@exportedonly" // 标记是否只导出
	Name            = "@name"         // 标记函数名称
	DefaultFuncName = "DeepCopy"      // 默认函数名称
)

// Option 包含了生成器选项的结构体
type Option struct {
	Name         string
	FuncName     string
	Generate     bool
	PtrRecv      bool
	Object       object
	ExportedOnly bool
}

// GeneratorOption 包含了生成器选项的切片
type GeneratorOption struct {
	Option
}

// ParseGeneratorOptions 从包中解析生成器选项
func ParseGeneratorOptions(p *packages.Package) [][]GeneratorOption {
	generatorOptions := make([][]GeneratorOption, len(p.CompiledGoFiles))
	for i, file := range p.Syntax {
		var options []GeneratorOption
		for _, decl := range file.Decls {
			if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.TYPE && gen.Doc != nil {
				var opt Option
				opt.PtrRecv = true
				opt.FuncName = DefaultFuncName
				opt.ExportedOnly = true
				for _, comment := range gen.Doc.List {
					opt.parseComment(comment.Text)
				}

				if opt.Generate {
					tp, _ := gen.Specs[0].(*ast.TypeSpec)
					obj, err := locateType(tp.Name.Name, p)
					if err != nil {
						panic(err)
					}
					opt.Name = tp.Name.Name
					opt.Object = obj
					options = append(options, GeneratorOption{Option: opt})
				}
			}
		}
		generatorOptions[i] = options
	}
	return generatorOptions
}

func (o *Option) parseComment(comment string) error {
	commentLine := strings.TrimSpace(strings.TrimLeft(comment, "/"))
	if len(commentLine) == 0 {
		return nil
	}

	fields := strings.Fields(commentLine)
	attr := strings.ToLower(fields[0])
	var lineRemainder string
	if len(fields) > 1 {
		lineRemainder = fields[1]
	}

	switch attr {
	case ExportedOnly:
		val := strings.ToLower(lineRemainder)
		o.ExportedOnly = val != "false"
	case PtrRecv:
		val := strings.ToLower(lineRemainder)
		o.PtrRecv = val != "false"
	case Cloneable:
		o.Generate = true
	case Name:
		o.FuncName = lineRemainder
	}
	return nil
}
