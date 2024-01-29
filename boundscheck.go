package linters

import (
	"go/ast"
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
		//ast.Inspect(file, func(n ast.Node) bool {
		//	if comment, ok := n.(*ast.Comment); ok {
		//		if strings.HasPrefix(comment.Text, "// TODO:") || strings.HasPrefix(comment.Text, "// TODO():") {
		//			pass.Report(analysis.Diagnostic{
		//				Pos:            comment.Pos(),
		//				End:            0,
		//				Category:       "todo",
		//				Message:        "TODO comment has no author",
		//				SuggestedFixes: nil,
		//			})
		//		}
		//	}
		//
		//	return true
		//})

		ast.Inspect(file, func(n ast.Node) bool {

			switch x := n.(type) {
			case *ast.IfStmt:
				//a := n.(*ast.IfStmt)
				//b := a.Cond.(*ast.BinaryExpr)
				//left := b.X.(*ast.CallExpr)
				//
				//var rightValue string
				//
				//right, ok := b.Y.(*ast.BasicLit)
				//if ok {
				//	rightValue = right.Value
				//} else {
				//	//rightValue = b.Y.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.BasicLit).Value
				//	//if ident, ok := b.Y.(*ast.Ident); ok {
				//	//	if assign, ok := ident.Obj.Decl.(*ast.AssignStmt); ok {
				//	//		if lit, ok := assign.Rhs[0].(*ast.BasicLit); ok {
				//	//			rightValue = lit.Value
				//	//		}
				//	//	}
				//	//}
				//	ident, ok := b.Y.(*ast.Ident)
				//	if !ok {
				//		//return false
				//
				//	}
				//	assign, ok := ident.Obj.Decl.(*ast.AssignStmt)
				//	if !ok {
				//		//return false
				//	}
				//	lit, ok := assign.Rhs[0].(*ast.BasicLit)
				//	if !ok {
				//		//return false
				//	}
				//	rightValue = lit.Value
				//	lenCheckMap[left.Args[0].(*ast.Ident).Name] = rightValue
				//}
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
					rightValue = b.Y.(*ast.Ident).Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.BasicLit).Value
				}

				ident, ok := left.Args[0].(*ast.Ident)
				if !ok {
					return false
				}
				lenCheckMap[ident.Name] = rightValue

			case *ast.IndexExpr:
				ident, ok := x.X.(*ast.Ident)
				if !ok {
					return false
				}
				name := ident.Name
				var value string

				switch i := x.Index.(type) {
				case *ast.BasicLit:
					value = i.Value

				case *ast.Ident:
					a, ok := i.Obj.Decl.(*ast.AssignStmt).Rhs[0].(*ast.BasicLit)
					if !ok {
						return false
					}
					value = a.Value
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
	}

	return nil, nil
}
