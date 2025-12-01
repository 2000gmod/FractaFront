package ast

import (
	"encoding/gob"
	"fmt"
)

func init() {
	fmt.Println("Registering AST types")
	gob.Register(ExprBase{})
	gob.Register(Literal{})
	gob.Register(Identifier{})
	gob.Register(Unary{})
	gob.Register(Binary{})
	gob.Register(Paren{})
	gob.Register(Call{})
	gob.Register(Indexed{})

	gob.Register(StmtBase{})
	gob.Register(ArgPair{})
	gob.Register(FunctionDeclaration{})
	gob.Register(ReturnStatement{})
	gob.Register(ExpressionStatement{})
	gob.Register(BlockStatement{})

	gob.Register(NamedType{})
}
