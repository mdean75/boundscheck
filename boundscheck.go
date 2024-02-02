package linters

import (
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"strconv"
)

var BoundsAnalyzer = &analysis.Analyzer{
	Name: "boundscheck",
	Doc:  "ensures slice index bounds check before use",
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	lenCheckMap := make(map[string]string)
	//valid := true
	//errMsg := ""

	for _, file := range pass.Files {

		ast.Inspect(file, func(n ast.Node) bool {

			switch x := n.(type) {
			case *ast.IfStmt:

				a := n.(*ast.IfStmt)
				b, ok := a.Cond.(*ast.BinaryExpr)
				if !ok {
					return false
				}
				call, ok := b.X.(*ast.CallExpr)
				if !ok {
					return false
				}
				left := call

				var rightValue string

				right, ok := b.Y.(*ast.BasicLit)
				if ok {
					rightValue = right.Value
				} else {
					rightIdent, ok := b.Y.(*ast.Ident)
					if !ok {
						return false
					}
					if rightIdent == nil {
						return false
					}
					if rightIdent.Obj == nil {
						return false
					}
					rightAssign, ok := rightIdent.Obj.Decl.(*ast.AssignStmt)
					if !ok {
						return false
					}
					rightLit, ok := rightAssign.Rhs[0].(*ast.BasicLit) //.Value
					if !ok {
						return false
					}
					rightValue = rightLit.Value
				}

				ident, ok := left.Args[0].(*ast.Ident)
				if !ok {
					return false
				}
				lenCheckMap[ident.Name] = rightValue

			case *ast.IndexExpr:

				if a, ok := x.Index.(*ast.BinaryExpr); ok {
					if aa, ok := a.X.(*ast.CallExpr); ok {
						if funname, ok := aa.Fun.(*ast.Ident); ok {
							if funname.Name == "len" {
								return false
							}

						}
					}
				}

				ident, ok := x.X.(*ast.Ident)
				if !ok {
					return false
				}

				name := ident.Name
				var value string

				if ident.Obj == nil {
					return false
				}
				if ident.Obj.Decl == nil {
					return false
				}
				if _, ok := ident.Obj.Decl.(*ast.Field); ok {
					return false
				}
				if _, ok := ident.Obj.Decl.(*ast.ValueSpec); ok {
					return false
				}
				if a, ok := ident.Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.CallExpr); ok {
					if _, ok := a.Args[0].(*ast.MapType); ok {
						return false
					}

				}
				switch i := x.Index.(type) {
				case *ast.BasicLit:
					// TODO: NOT QUITE RIGHT, NEED TO CHECK THE VARIABLE TYPE IS NOT MAP
					m, ok := x.X.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.CompositeLit) //.Args[0].(*ast.MapType); !ok { //.Decl.(*ast.AssignStmt).Rhs[0].(*ast.CallExpr).Args[0].(*ast.MapType); !ok {
					if !ok {
						return false
					}
					_, ok = m.Type.(*ast.ArrayType)
					if !ok {
						return false
					}

					//_, ok = m.Args[0].(*ast.MapType)
					//if !ok {
					//	return false
					//}

					//x.X.Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.CallExpr).Args[0] == maptype
					if i.Kind != token.INT {
						return false
					}
					value = i.Value

				case *ast.Ident:
					if i.Obj == nil {
						return false
					}
					if i.Obj.Decl == nil {
						return false
					}

					if _, ok := i.Obj.Decl.(*ast.Field); ok {
						return false
					}
					a, ok := i.Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.BasicLit)
					if !ok {
						return false
					}
					value = a.Value

				case *ast.CallExpr:
					_, ok := i.Args[0].(*ast.MapType)
					if ok {
						return false
					}

				}

				// lookup bounds check
				chk, ok := lenCheckMap[name]
				if !ok {
					//valid = false
					//errMsg = "missing bounds check"

					pass.Reportf(x.Pos(), "missing bounds check")
					return true
				}
				chkindex, _ := strconv.Atoi(chk)
				ivalue, _ := strconv.Atoi(value)

				if chkindex < ivalue {
					//valid = false
					//errMsg = "ineffective bounds check, index check less than index use"
					pass.Reportf(x.Pos(), "ineffective bounds check, index check less than index use")
					return true
				}

			}

			return true
		})
		//fmt.Println(lenCheckMap)
	}

	return nil, nil
}
