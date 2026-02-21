package sema

import (
	"errors"
	"fmt"
	"fracta/internal/ast"
	"fracta/internal/diag"
	"fracta/internal/token"
)

func (a *SemanticAnalyzer) addErrorStmt(stmt *ast.StmtBase, f string, v ...any) {
	msg := fmt.Sprintf(f, v...)
	o := diag.CreateError(msg, a.currentFile, stmt.Line)
	a.errors = append(a.errors, o)
}

func (a *SemanticAnalyzer) addErrorExpr(expr *ast.ExprBase, f string, v ...any) {
	msg := fmt.Sprintf(f, v...)
	o := diag.CreateError(msg, a.currentFile, expr.Line)
	a.errors = append(a.errors, o)
}

func (a *SemanticAnalyzer) createScope() {
	a.currentScope = a.currentScope.newChildScope()
}

func (a *SemanticAnalyzer) dropScope() {
	a.currentScope = a.currentScope.parent
}

func (a *SemanticAnalyzer) Analyze() ([]*ast.FileSourceNode, error) {
	for _, fileAst := range a.packageAsts {
		a.currentFile = fileAst.Filename
		a.populatePackageSymbolTable(fileAst)
	}

	if len(a.errors) == 0 {
		for _, fn := range a.packageAsts {
			a.currentFile = fn.Filename
			a.analyzeFileNode(fn)
		}
	}

	errList := make([]error, 0)
	diag.AppendError(a.errors...)

	for _, e := range a.errors {
		errList = append(errList, e)
	}

	if len(a.errors) > 0 {
		return nil, errors.Join(errList...)
	}

	return a.packageAsts, nil
}

func (a *SemanticAnalyzer) populatePackageSymbolTable(fileTree *ast.FileSourceNode) {
	for _, stmt := range fileTree.Statements {
		switch s := stmt.(type) {
		case *ast.FunctionDeclaration:
			a.populateFunctionDecl(s)
		default:
			a.addErrorStmt(stmt.StmtNode(), "invalid statement, only declarations are allowed in top-level scope")
		}
	}
}

func (a *SemanticAnalyzer) populateFunctionDecl(fd *ast.FunctionDeclaration) {
	a.pkgScope.addSymbol(fd.Name.Identifier, &functionSymbol{
		symbolBase: symbolBase{pkg: a.packageName},
		fType:      ast.FuncDeclToFuncType(fd),
	})
}

func (a *SemanticAnalyzer) analyzeFileNode(fn *ast.FileSourceNode) {
	for _, st := range fn.Statements {
		a.analyzeTopLevelStatement(st)
	}
}

func (a *SemanticAnalyzer) analyzeTopLevelStatement(st ast.Statement) {
	switch s := st.(type) {
	case *ast.FunctionDeclaration:
		a.analyzeFunctionDecl(s)
	default:
		a.addErrorStmt(st.StmtNode(), "invalid top level statement")
	}
}

func (a *SemanticAnalyzer) analyzeFunctionDecl(fd *ast.FunctionDeclaration) {
	a.currentFunction = fd
	defer func() { a.currentFunction = nil }()

	if fd.Body != nil {
		body, ok := fd.Body.(*ast.BlockStatement)
		if !ok {
			a.addErrorStmt(&fd.StmtBase, "only block statements are allowed in a function body")
			return
		}
		a.analyzeBlockStatement(body)
	}
}

func (a *SemanticAnalyzer) analyzeStatement(st ast.Statement) {
	switch s := st.(type) {
	case *ast.ReturnStatement:
		a.analyzeReturnStatement(s)
	case *ast.BlockStatement:
		a.analyzeBlockStatement(s)
	case *ast.ExpressionStatement:
		a.analyzeExpressionStatement(s)
	default:
		a.addErrorStmt(st.StmtNode(), "invalid statement in this position")
	}
}

func (a *SemanticAnalyzer) analyzeReturnStatement(ret *ast.ReturnStatement) {
	if a.currentFunction.ReturnType == nil {
		if ret.Value != nil {
			a.addErrorStmt(&ret.StmtBase, "return has value in a void function")
		}
		return
	}

	if ret.Value != nil {
		a.analyzeExpression(ret.Value)
	}

	retType := ret.Value.ExprNode().Type

	if !ast.CompareTypes(retType, a.currentFunction.ReturnType) {
		a.addErrorStmt(&ret.StmtBase, "return type mismatch, expression of type %q, expected %q", retType.String(), a.currentFunction.ReturnType.String())
		return
	}

}

func (a *SemanticAnalyzer) analyzeBlockStatement(bl *ast.BlockStatement) {
	a.createScope()
	defer a.dropScope()

	for _, st := range bl.Body {
		a.analyzeStatement(st)
	}
}

func (a *SemanticAnalyzer) analyzeExpressionStatement(est *ast.ExpressionStatement) {
	a.analyzeExpression(est.Expression)
}

func (a *SemanticAnalyzer) analyzeExpression(expr ast.Expression) {
	switch e := expr.(type) {
	case *ast.Literal:
		a.analyzeLiteralExpr(e)
	case *ast.Identifier:
		a.analyzeIdentifierExpr(e)
	case *ast.Unary:
		a.analyzeUnaryExpr(e)
	case *ast.Binary:
		a.analyzeBinaryExpr(e)
	case *ast.Call:
		a.analyzeCallExpr(e)
	case *ast.Indexed:
		a.analyzeIndexedExpr(e)
	default:
		a.addErrorExpr(expr.ExprNode(), "unknown expression kind")
	}
}

func (a *SemanticAnalyzer) analyzeLiteralExpr(e *ast.Literal) {
	etype, ok := ast.TokenLiteralMap[e.Value.Kind]
	if !ok {
		a.addErrorExpr(&e.ExprBase, "literal not yet supported: %s", e.Value.String())
		return
	}

	e.Type = etype
}

func (a *SemanticAnalyzer) analyzeIdentifierExpr(e *ast.Identifier) {
	sym, ok := a.currentScope.getSymbol(e.Ident.Identifier)
	if !ok {
		a.addErrorExpr(&e.ExprBase, "used but not defined: %s", e.Ident.Identifier)
		return
	}
	e.Type = sym.getExprType()
}

func (a *SemanticAnalyzer) analyzeUnaryExpr(e *ast.Unary) {
	a.analyzeExpression(e.SubExpr)

	switch e.Op.Kind {
	case token.TokOpPlus, token.TokOpMinus:
		if !ast.IsNumeric(e.Type) {
			a.addErrorExpr(&e.ExprBase, "non-numeric expression type for unary expression")
			return
		}
	default:
		a.addErrorExpr(&e.ExprBase, "invalid operator for unary expression")
		return
	}
}

func (a *SemanticAnalyzer) analyzeBinaryExpr(e *ast.Binary) {
	a.analyzeExpression(e.Left)
	a.analyzeExpression(e.Right)

	if !ast.CompareTypes(e.Left.ExprNode().Type, e.Right.ExprNode().Type) {
		a.addErrorExpr(&e.ExprBase, "mismatched types for expression")
		return
	}
}

func (a *SemanticAnalyzer) analyzeCallExpr(e *ast.Call) {
	a.addErrorExpr(&e.ExprBase, "call expression not supported yet")
}

func (a *SemanticAnalyzer) analyzeIndexedExpr(e *ast.Indexed) {
	a.addErrorExpr(&e.ExprBase, "index expression not supported yet")
}
