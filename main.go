package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
)

func main() {
	// 解析源代码文件
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "source.go", nil, parser.ParseComments)
	if err != nil {
		fmt.Println("解析文件失败:", err)
		return
	}

	// 创建新文件
	newFile := &ast.File{
		Name:  &ast.Ident{Name: file.Name.Name + "aop"},
		Decls: []ast.Decl{},
	}

	// 遍历文件的AST节点
	ast.Inspect(file, func(n ast.Node) bool {
		// 找到函数声明节点
		if fn, ok := n.(*ast.FuncDecl); ok {
			// 创建新的函数声明节点
			// 生成参数

			newFn := &ast.FuncDecl{
				Name: fn.Name,
				Type: &ast.FuncType{
					Params:  fn.Type.Params,
					Results: fn.Type.Results,
				},
				Body: &ast.BlockStmt{
					List: []ast.Stmt{
						// 在新函数体中调用旧函数
						&ast.ExprStmt{
							X: &ast.CallExpr{
								Fun: &ast.SelectorExpr{
									X:   &ast.Ident{Name: file.Name.Name},
									Sel: fn.Name,
								},
								Args: createArgs(fn.Type.Params), // 在这里添加参数
							},
						},
					},
				},
			}

			// 添加新函数到新文件的Decls字段中
			newFile.Decls = append(newFile.Decls, newFn)
		}
		return true
	})

	// 生成新文件
	outputFile, err := os.Create("modified.go")
	if err != nil {
		fmt.Println("创建文件失败:", err)
		return
	}
	defer outputFile.Close()

	err = printer.Fprint(outputFile, fset, newFile)
	if err != nil {
		fmt.Println("格式化文件失败:", err)
		return
	}

	fmt.Println("AOP代码生成成功")
}

// 创建参数列表
func createArgs(params *ast.FieldList) []ast.Expr {
	args := make([]ast.Expr, len(params.List))
	for i, param := range params.List {
		// 使用参数的标识符作为新函数的参数
		args[i] = &ast.Ident{Name: param.Names[0].Name}
	}
	return args
}
