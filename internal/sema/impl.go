package sema

import (
	"errors"
	"fmt"
	"fracta/internal/ast"
	"fracta/internal/diag"
)

func (a *SemanticAnalyzer) addError(stmt *ast.StmtBase, f string, v ...any) {
	msg := fmt.Sprintf(f, v...)
	o := diag.CreateError(msg, a.currentFile, stmt.Line)
	a.errors = append(a.errors, o)
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
			a.addError(stmt.StmtNode(), "invalid statement, only declarations are allowed in top-level scope")
		}
	}
}

func (a *SemanticAnalyzer) populateFunctionDecl(fd *ast.FunctionDeclaration) {
	a.pkgScope.addSymbol(fd.Name.Identifier, &functionSymbol{
		symbolBase: symbolBase{pkg: a.packageName},
		decl:       fd,
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
		a.addError(st.StmtNode(), "invalid top level statement")
	}
}

func (a *SemanticAnalyzer) analyzeFunctionDecl(fd *ast.FunctionDeclaration) {
	a.currentFunction = fd
	defer func() { a.currentFunction = nil }()

	if fd.Body != nil {
		body, ok := fd.Body.(*ast.BlockStatement)
		if !ok {
			a.addError(&fd.StmtBase, "only block statements are allowed in a function body")
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
		a.addError(st.StmtNode(), "invalid statement in this position")
	}
}

func (a *SemanticAnalyzer) analyzeReturnStatement(ret *ast.ReturnStatement) {
	if a.currentFunction.ReturnType == nil {
		if ret.Value != nil {
			a.addError(&ret.StmtBase, "return has value in a void function")
			return
		}
	}

	if ret.Value != nil {
		a.analyzeExpression(ret.Value)
	}

	retType := ret.Value.ExprNode().Type

	if !ast.CompareTypes(retType, a.currentFunction.ReturnType) {
		a.addError(&ret.StmtBase, "type mismatch, expression of type %q, expected %q", retType.String(), a.currentFunction.ReturnType.String())
		return
	}

}

func (a *SemanticAnalyzer) analyzeBlockStatement(bl *ast.BlockStatement) {
	for _, st := range bl.Body {
		a.analyzeStatement(st)
	}
}

func (a *SemanticAnalyzer) analyzeExpressionStatement(est *ast.ExpressionStatement) {

}

func (a *SemanticAnalyzer) analyzeExpression(expr ast.Expression) {

}
